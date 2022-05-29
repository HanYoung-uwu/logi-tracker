package main

import (
	"fmt"
	database "hanyoung/logi-tracker/internal/database"

	gin "github.com/gin-gonic/gin"
)

func main() {
	database.Test()
	fmt.Println("Hello, World!")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
