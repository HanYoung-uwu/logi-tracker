package main

import (
	"fmt"
	"hanyoung/logi-tracker/internal/handlers"
	"hanyoung/logi-tracker/internal/loginmiddleware.go"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello, World!")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/register", loginmiddleware.CreateUserHandler)
	r.POST("/login", loginmiddleware.LoginHandler)
	authorized := r.Group("/user", loginmiddleware.DefaultAuthHandler)

	authorized.GET("/all_items", handlers.GetAllItemsHandler)
	authorized.POST("/create_stockpile", handlers.CreateStockpileHandler)
	authorized.POST("/insert_item", handlers.InsertItemHandler)
	authorized.GET("/all_stockpiles", handlers.GetAllLocationsHandler)
	r.Run()
}
