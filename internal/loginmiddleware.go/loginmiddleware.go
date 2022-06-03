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
		return nil, ErrorInvalideToken
	}
	return _token.account, nil
}

func DefaultAuthHandler(r *gin.Context) {
	log.Println(r.Request.Header)
	token, err := r.Cookie("token")
	if err == nil {
		account, err := GetAccountManager().GetAccountByToken(token)
		if err == nil {
			r.Set("account", account)
			r.Next()
			return
		}
	}
	r.Abort()
	r.JSON(http.StatusUnauthorized, gin.H{"reason": "unauthorized"})
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
	database.GetInstance().AddAccount(json.Name, json.Password, "test", 0)
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
