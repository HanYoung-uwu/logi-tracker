package handlers

import (
	"hanyoung/logi-tracker/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Stockpile struct {
	Location string `form:"name" json:"name" xml:"name"  binding:"required"`
	Code     string `form:"code" json:"code" xml:"code"  binding:"-"`
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
	if err := c.ShouldBindJSON(&json); err != nil || len(json.Code) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.GetInstance().CreateStockpile(json.Location, json.Code, _account.Clan)
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

	err := database.GetInstance().InsertOrUpdateItem(json.Location, json.ItemType, json.Size, _account.Clan, _account.Name)
	if err != nil {
		c.JSON(http.StatusNotModified, "unable to update item")
		return
	}
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

func DeleteStockpileHandler(c *gin.Context) {
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
	err := database.GetInstance().DeleteStockpile(json.Location, _account.Clan)
	if err != nil {
		c.JSON(http.StatusBadRequest, "stockpile doesn't exits")
		return
	}
	c.JSON(http.StatusOK, "success")
}
