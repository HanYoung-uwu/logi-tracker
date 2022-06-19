package loginmiddleware

import (
	"errors"
	"hanyoung/logi-tracker/internal/database"
	"hanyoung/logi-tracker/pkg/utility"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name     string `form:"Name" json:"Name" xml:"Name"  binding:"required"`
	Password string `form:"Password" json:"Password" xml:"Password" binding:"required"`
	Clan     string `form:"Clan" json:"Clan" xml:"Clan" binding:"-"`
}

type Token struct {
	value      string
	expireTime time.Time
	account    *database.Account
}

type TokenManager struct {
	tokens *sync.Map
	ticker *time.Ticker
}

var singleton *TokenManager
var lock = &sync.Mutex{}

func GetAccountManager() *TokenManager {
	if singleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleton == nil {
			singleton = &TokenManager{&sync.Map{}, time.NewTicker(10 * time.Hour)}
			go singleton.gc()
		}
	}
	return singleton
}

func (t *TokenManager) gc() {
	for {
		now := <-t.ticker.C
		cleaned := 0
		log.Println("Start gc for TokenManager")
		t.tokens.Range(func(key interface{}, val interface{}) bool {
			token, ok := val.(*Token)
			if !ok {
				log.Panic("invalid token")
			}
			if token.expireTime.Before(now) {
				t.tokens.Delete(key)
				cleaned += 1
			}
			return true
		})
		log.Print("Cleaned", cleaned, "tokens")
	}
}

func (t *TokenManager) ValidateAccountAndGenerateToken(name string, password string) (string, error) {
	dbMgr := database.GetInstance()
	account, err := dbMgr.GetAndValidateAccount(name, password)
	if err != nil {
		if err == database.ErrorIncorrectPassword || err == database.ErrorNoAccount {
			return "", err
		}
		log.Fatal(err)
	}

	accessCookie := utility.RandBytes(256)
	for {
		if _, exists := t.tokens.LoadOrStore(string(accessCookie), &Token{string(accessCookie), time.Now().Add(1000000000 * 3600 * 24 * 14), account}); exists {
			accessCookie = utility.RandBytes(256)
		} else {
			break
		}
	}

	return string(accessCookie), nil
}

var ErrorInvalideToken = errors.New("invalide token")

func (t *TokenManager) GetAccountByToken(token string) (*database.Account, error) {
	val, found := t.tokens.Load(token)
	if found {
		m_token, ok := val.(*Token)
		if !ok {
			log.Panic("m_token is not a *Token")
		}
		if m_token.expireTime.Before(time.Now()) {
			t.tokens.Delete(token)
			return nil, ErrorInvalideToken
		}
		return m_token.account, nil
	}
	return nil, ErrorInvalideToken
}

func (t *TokenManager) GenerateInvitationToken(clan string) string {
	account := &database.Account{Name: "tmp", Permission: database.InvitationLinkAccount, Clan: clan}

	token := utility.RandBytes(92)
	for {
		// invitation links are only valid for one day
		if _, exists := t.tokens.LoadOrStore(string(token), &Token{string(token), time.Now().Add(1000000000 * 3600 * 24), account}); exists {
			token = utility.RandBytes(92)
		} else {
			break
		}
	}

	return string(token)
}

func (t *TokenManager) GenerateClanAdminInvitationToken() string {
	account := &database.Account{Name: "tmp", Permission: database.ClanAdminInvitationLinkAccount, Clan: ""}

	token := utility.RandBytes(92)
	for {
		// invitation links are only valid for one day
		if _, exists := t.tokens.LoadOrStore(string(token), &Token{string(token), time.Now().Add(1000000000 * 3600 * 24), account}); exists {
			token = utility.RandBytes(92)
		} else {
			break
		}
	}

	return string(token)
}

func DefaultAuthHandler(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		account, err := GetAccountManager().GetAccountByToken(token)
		if err == nil {
			if account.Permission == database.AdminAccount || account.Permission == database.ClanAdminAccount || account.Permission == database.NormalAccount {
				c.Set("account", account)
				c.Next()
				return
			}
		}
	}
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{"reason": "unauthorized"})
}

func UserInfoAuthHandler(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		account, err := GetAccountManager().GetAccountByToken(token)
		if err == nil {
			c.Set("account", account)
			c.Next()
			return
		}
	}
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{"reason": "unauthorized"})
}

func AdminAuthHandler(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		account, err := GetAccountManager().GetAccountByToken(token)
		if err == nil && account.Permission == 0 {
			c.Set("account", account)
			c.Next()
			return
		}
	}
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{"reason": "unauthorized"})
}

func ClanAdminAuthHandler(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		account, err := GetAccountManager().GetAccountByToken(token)
		if err == nil && account.Permission <= database.ClanAdminAccount {
			c.Set("account", account)
			c.Next()
			return
		}
	}
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{"reason": "unauthorized"})
}

func CreateUserHandler(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(json.Password) < 8 || len(json.Name) == 0 {
		c.Abort()
		c.JSON(http.StatusNotAcceptable, gin.H{"reason": "password or name too short"})
		return
	}
	if len(json.Password) > 72 {
		json.Password = json.Password[:71]
	}
	account, exists := c.Get("account")
	if !exists {
		log.Println("can't get account")
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_account, ok := account.(*database.Account)
	if !ok {
		log.Panic("account is not a *Account")
	}

	database.GetInstance().AddAccount(json.Name, json.Password, _account.Clan, 2)
	c.JSON(200, gin.H{
		"message": "succeed",
	})
}

func LoginHandler(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := GetAccountManager().ValidateAccountAndGenerateToken(json.Name, json.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "failed to validate"})
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 3600*24*14, "/", "", !utility.DebugEnvironment, !utility.DebugEnvironment)
	c.JSON(http.StatusAccepted, "success")
}

func CreateUserFromInvitationLinkHandler(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(json.Password) < 8 || len(json.Name) == 0 {
		c.Abort()
		c.JSON(http.StatusNotAcceptable, gin.H{"reason": "password or name too short"})
		return
	}
	if len(json.Password) > 72 {
		json.Password = json.Password[:71]
	}

	token, err := c.Cookie("token")
	if err == nil {
		account, err := GetAccountManager().GetAccountByToken(token)
		if err == nil && account.Permission >= database.InvitationLinkAccount {
			var permission int
			clan := account.Clan
			switch account.Permission {
			case database.ClanAdminInvitationLinkAccount:
				permission = database.ClanAdminAccount
				if json.Clan == "" {
					c.JSON(http.StatusBadRequest, "must specify clan")
					return
				} else {
					clan = json.Clan
				}
			case database.InvitationLinkAccount:
				permission = database.NormalAccount
			}

			err = database.GetInstance().AddAccount(json.Name, json.Password, clan, permission)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "succeed",
			})
			GetAccountManager().tokens.Delete(token)
			return
		}
	}
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{"reason": "unauthorized"})
}

func GenerateInvitationLinkHandler(c *gin.Context) {
	account, exists := c.Get("account")
	if !exists {
		log.Println("can't get account")
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_account, ok := account.(*database.Account)
	if !ok {
		log.Panic("account is not a *Account")
	}

	token := GetAccountManager().GenerateInvitationToken(_account.Clan)
	c.JSON(http.StatusAccepted, gin.H{"token": token})
}

func GenerateClanAdminInvitationLinkHandler(c *gin.Context) {
	token := GetAccountManager().GenerateClanAdminInvitationToken()
	c.JSON(http.StatusAccepted, gin.H{"token": token})
}

func LogoutHandler(c *gin.Context) {
	account, exists := c.Get("account")
	if !exists {
		log.Println("can't get account")
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_account, ok := account.(*database.Account)
	if !ok {
		log.Panic("account is not a *Account")
	}
	GetAccountManager().tokens.Delete(_account.Name)
	c.SetCookie("token", "", 0, "", "", !utility.DebugEnvironment, !utility.DebugEnvironment)
	c.JSON(http.StatusAccepted, "")
}
