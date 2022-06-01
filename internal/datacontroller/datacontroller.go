package datacontroller

import (
	dbmanager "hanyoung/logi-tracker/internal/database"

	gin "github.com/gin-gonic/gin"
)

type DataController struct {
	db     *dbmanager.DataBaseManager
	admins map[string]bool
}

func GetInstance() *DataController {
	db := dbmanager.GetInstance()
	admins := make(map[string]bool)
	return &DataController{db, admins}
}

func (d *DataController) ValidateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		if d.admins[user] {
			c.Next()
		} else {
			c.Abort()
		}
	}
}
