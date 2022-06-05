package loginmiddleware

import (
	"errors"
	"hanyoung/logi-tracker/internal/database"
	"hanyoung/logi-tracker/pkg/utility"
	"io/fs"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var createAdminOnFirstRequestHandlerLock = &sync.Mutex{}

var used = false

func CreateAdminOnFirstRequestHandler(c *gin.Context) {
	if !used {
		createAdminOnFirstRequestHandlerLock.Lock()
		defer createAdminOnFirstRequestHandlerLock.Unlock()
		if !used {
			used = true

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
			if len(json.Password) > 72 {
				json.Password = json.Password[:71]
			}

			// if database file has already exist, is also not first time
			_, err := os.Stat(utility.DatabasePath)
			if err != nil && errors.Is(err, fs.ErrNotExist) {
				err := database.GetInstance().AddAccount(json.Name, json.Password, "admin", database.AdminAccount)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": err.Error(),
					})
					return
				} else {
					c.JSON(200, gin.H{
						"message": "succeed",
					})
					return
				}
			}
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "not first time request",
	})
}
