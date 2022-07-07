package handlers

import (
	"hanyoung/logi-tracker/internal/database"
	"hanyoung/logi-tracker/pkg/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Name = utility.Name

type Stockpile struct {
	Location string `form:"name" json:"name" xml:"name"  binding:"required"`
	Code     string `form:"code" json:"code" xml:"code"  binding:"-"`
}

type Faction struct {
	Faction int `form:"faction" json:"faction" xml:"faction"  binding:"required"`
}

func GetAllItemsHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		result := database.GetInstance().GetAllItems(account)
		c.JSON(http.StatusOK, result)
	}
}

func CreateStockpileHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		var json Stockpile
		if err := c.ShouldBindJSON(&json); err != nil || len(json.Code) == 0 {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			return
		}

		database.GetInstance().CreateStockpile(json.Location, json.Code, account.Clan, account.Name)
		c.JSON(http.StatusAccepted, "success")
	}
}

func InsertOrUpdateItemHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		var json database.StockpileItem
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := database.GetInstance().InsertOrUpdateItem(json.Location, json.ItemType, json.Size, account.Clan, account.Name)
		if err != nil {
			c.JSON(http.StatusNotModified, "unable to update item")
			return
		}
		c.JSON(http.StatusAccepted, "success")
	}
}

func GetAllLocationsHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		c.JSON(http.StatusAccepted, database.GetInstance().GetAllLocations(account.Clan))
	}
}

func DeleteStockpileHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		var json Stockpile
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.GetInstance().DeleteStockpile(json.Location, account.Clan, account.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, "stockpile doesn't exits")
			return
		}
		c.JSON(http.StatusOK, "success")
	}
}

func DeleteItemHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		var json database.StockpileItem
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		database.GetInstance().DeleteItem(json.Location, json.ItemType, account.Clan, account.Name)
		c.JSON(http.StatusOK, "success")
	}
}

func SetItemHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {

		var json database.StockpileItem
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		database.GetInstance().SetItem(json.Location, json.ItemType, json.Size, account.Clan, account.Name)
		c.JSON(http.StatusOK, "success")
	}
}

func RefreshStockpileHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		var json Stockpile
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		database.GetInstance().RefreshStockpile(json.Location, account.Clan)
		c.JSON(http.StatusOK, "success")
	}
}

func SetClanFactionHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		var json Faction
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if json.Faction != 1 && json.Faction != 2 {
			c.JSON(http.StatusBadRequest, "faction need to be either 1 or 2")
		}
		database.GetInstance().SetClanFaction(account.Clan, json.Faction)
		c.JSON(http.StatusOK, "success")
	}
}

func GetClanFactionHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		faction := database.GetInstance().GetClanFaction(account.Clan)
		c.JSON(http.StatusOK, gin.H{"faction": faction})
	}
}
