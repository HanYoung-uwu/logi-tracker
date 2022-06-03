package main

import (
	"hanyoung/logi-tracker/internal/handlers"
	"hanyoung/logi-tracker/internal/loginmiddleware.go"
	"hanyoung/logi-tracker/pkg/utility"

	"github.com/gin-gonic/gin"
)

func main() {
	utility.InitConfig()

	r := gin.Default()
	r.POST("/login", loginmiddleware.LoginHandler)

	authorized := r.Group("/user", loginmiddleware.DefaultAuthHandler)
	authorized.GET("/all_items", handlers.GetAllItemsHandler)
	authorized.POST("/create_stockpile", handlers.CreateStockpileHandler)
	authorized.POST("/update_item", handlers.InsertOrUpdateItemHandler)
	authorized.GET("/all_stockpiles", handlers.GetAllLocationsHandler)

	clanAdmins := r.Group("/clan", loginmiddleware.ClanAdminAuthHandler)
	clanAdmins.POST("/register", loginmiddleware.CreateUserHandler)
	clanAdmins.GET("/invitation", loginmiddleware.GenerateInvitationLinkHandler)
	r.Run()
}
