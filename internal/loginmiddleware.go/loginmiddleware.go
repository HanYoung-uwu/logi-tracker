package loginmiddleware

import (
	"errors"
	"hanyoung/logi-tracker/internal/database"
	"hanyoung/logi-tracker/pkg/utility"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

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
	_token, err := t.tokens[token]
	if err || _token.expireTime.Before(time.Now()) {
		return nil, ErrorInvalideToken
	}
	return _token.account, nil
}

func DefaultAuth(r *gin.Context) {
	token, err := r.Cookie("token")
	if err != nil {
		r.Abort()
		return
	}
	account, err := GetAccountManager().GetAccountByToken(token)
	if err != nil {
		r.Abort()
		return
	}

	r.Set("account", account)
	r.Next()
}
