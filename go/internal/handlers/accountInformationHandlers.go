package handlers

import (
	"hanyoung/logi-tracker/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBasicAccountInfo(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"Name": _account.Name, "Clan": _account.Clan, "Permission": _account.Permission})
}