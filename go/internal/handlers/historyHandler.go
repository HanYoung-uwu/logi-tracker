package handlers

import (
	"hanyoung/logi-tracker/internal/database"
	"hanyoung/logi-tracker/pkg/utility"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClanHistoryHandler(c *gin.Context) {
	account, exists := c.Get("account")
	if !exists {
		log.Println("can't get account")
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_account, ok := account.(*utility.Account)
	if !ok {
		log.Panic("account is not a *Account")
	}

	c.JSON(http.StatusOK, database.GetInstance().GetClanHistory(_account.Clan))
}
