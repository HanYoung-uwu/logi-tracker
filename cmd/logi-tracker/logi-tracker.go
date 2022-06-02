package main

import (
	"fmt"
	"hanyoung/logi-tracker/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUser struct {
	Name     string `form:"Name" json:"Name" xml:"Name"  binding:"required"`
	Password string `form:"Password" json:"Password" xml:"Password" binding:"required"`
}

func main() {
	// manager := database.GetInstance()

	fmt.Println("Hello, World!")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/register", func(c *gin.Context) {
		var json CreateUser
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(json.Password) < 100 {
			padding := make([]byte, 100-len(json.Password))
			json.Password = string(append(padding, json.Password...))
		}
		database.GetInstance().AddAccount(json.Name, json.Password, "test", 0)
		c.JSON(200, gin.H{
			"message": "succeed",
		})
	})
	r.Run()
}
