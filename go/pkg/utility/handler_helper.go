package utility

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAccount(c *gin.Context) *Account {
	account, exists := c.Get("account")
	if !exists {
		log.Println("can't get account")
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil
	}
	_account, ok := account.(*Account)
	if !ok {
		log.Panic("account is not a *Account")
	}
	return _account
}
