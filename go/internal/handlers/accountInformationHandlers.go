package handlers

import (
	"hanyoung/logi-tracker/internal/database"
	"hanyoung/logi-tracker/pkg/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NameQuery struct {
	Name string `form:"name" json:"name" xml:"name"  binding:"required"`
}

func GetBasicAccountInfoHandler(c *gin.Context) {
	account := utility.GetAccount(c)
	if account != nil {
		c.JSON(http.StatusOK, gin.H{"Name": account.Name, "Clan": account.Clan, "Permission": account.Permission})
	}
}

func CheckAccountNameExistHandler(c *gin.Context) {
	var json NameQuery
	err := c.ShouldBind(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, "must supply name")
		return
	}
	switch database.GetInstance().IsNameExist(json.Name) {
	case true:
		c.JSON(http.StatusOK, gin.H{"exist": true})
	case false:
		c.JSON(http.StatusOK, gin.H{"exist": false})
	}
}

func CheckClanExistHandler(c *gin.Context) {
	var json NameQuery
	err := c.ShouldBind(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, "must supply name")
		return
	}
	switch database.GetInstance().IsClanExist(json.Name) {
	case true:
		c.JSON(http.StatusOK, gin.H{"exist": true})
	case false:
		c.JSON(http.StatusOK, gin.H{"exist": false})
	}
}
