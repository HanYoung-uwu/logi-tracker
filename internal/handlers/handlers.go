package handlers

import (
	"hanyoung/logi-tracker/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Stockpile struct {
	Location string `form:"name" json:"name" xml:"name"  binding:"required"`
	Code     string `form:"code" json:"code" xml:"code"  binding:"required"`
}

func GetAllItemsHandler(c *gin.Context) {
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

	result := database.GetInstance().GetAllItems(_account)
	c.JSON(http.StatusOK, result)
}

func CreateStockpileHandler(c *gin.Context) {
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

	var json Stockpile
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.GetInstance().CreateStockpile(json.Location, _account.Clan, json.Code)
	c.JSON(http.StatusAccepted, "success")
}

func InsertOrUpdateItemHandler(c *gin.Context) {
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

	var json database.StockpileItem
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.GetInstance().InsertOrUpdateItem(json.Location, json.ItemType, json.Size, _account.Clan)
	c.JSON(http.StatusAccepted, "success")
}

func GetAllLocationsHandler(c *gin.Context) {
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
	c.JSON(http.StatusAccepted, database.GetInstance().GetAllLocations(_account.Clan))
}
