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
}

type Token struct {
	value      string
	expireTime time.Time
	account    *database.Account
}

type TokenManager struct {
	tokens map[string]*Token
}

var singleton *TokenManager
var lock = &sync.Mutex{}

func GetAccountManager() *TokenManager {
	if singleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleton == nil {
			singleton = &TokenManager{make(map[string]*Token)}
		}
	}
	return singleton
}

func (t *TokenManager) ValidateAccountAndGenerateToken(name string, password string) (string, error) {
	dbMgr := database.GetInstance()
	account, err := dbMgr.GetAndValidateAccount(name, password)
	if err != nil {
		if err == database.ErrorIncorrectPassword {
			return "", err
		}
		log.Fatal(err)
	}

	accessCookie := utility.RandBytes(256)
	for {
		if _, exists := t.tokens[string(accessCookie)]; exists {
			accessCookie = utility.RandBytes(256)
		} else {
			break
		}
	}

	// default 14 day auto login
	t.tokens[string(accessCookie)] = &Token{string(accessCookie), time.Now().Add(1000000000 * 3600 * 24 * 14), account}
	return string(accessCookie), nil
}

var ErrorInvalideToken = errors.New("invalide token")

func (t *TokenManager) GetAccountByToken(token string) (*database.Account, error) {
	_token, found := t.tokens[token]
	if !found || _token.expireTime.Before(time.Now()) {
		delete(t.tokens, token)
		return nil, ErrorInvalideToken
	}
	return _token.account, nil
}

func (t *TokenManager) GenerateInvitationToken(clan string) string {
	account := &database.Account{Name: "tmp", Permission: 3, Clan: clan}

	token := utility.RandBytes(92)
	for {
		if _, exists := t.tokens[string(token)]; exists {
			token = utility.RandBytes(92)
		} else {
			break
		}
	}

	// invitation links are only valid for one day
	t.tokens[string(token)] = &Token{string(token), time.Now().Add(1000000000 * 3600 * 24), account}
	return string(token)
}

func DefaultAuthHandler(c *gin.Context) {
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
		if err == nil && account.Permission <= 1 {
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
	if len(json.Password) < 100 {
		padding := make([]byte, 100-len(json.Password))
		json.Password = string(append(padding, json.Password...))
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

	c.SetCookie("token", token, 3600*24*14, "/user", "", true, true)
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
	if len(json.Password) < 100 {
		padding := make([]byte, 100-len(json.Password))
		json.Password = string(append(padding, json.Password...))
	}

	token, err := c.Cookie("token")
	if err == nil {
		account, err := GetAccountManager().GetAccountByToken(token)
		if err == nil && account.Permission == database.InvitationLinkAccount {
			database.GetInstance().AddAccount(json.Name, json.Password, account.Clan, 2)
			c.JSON(200, gin.H{
				"message": "succeed",
			})
			delete(GetAccountManager().tokens, token)
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

	if _account.Permission > database.ClanAdminAccount {
		c.JSON(http.StatusUnauthorized, "not clan admin account")
		c.Abort()
		return
	}

	token := GetAccountManager().GenerateInvitationToken(_account.Clan)
	c.JSON(http.StatusAccepted, gin.H{"token": token})
}
