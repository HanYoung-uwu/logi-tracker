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
	r.POST("/register", loginmiddleware.CreateUserFromInvitationLinkHandler)
	r.POST("/admin/create_admin", loginmiddleware.CreateAdminOnFirstRequestHandler)

	authorized := r.Group("/user", loginmiddleware.DefaultAuthHandler)
	authorized.GET("/all_items", handlers.GetAllItemsHandler)
	authorized.POST("/create_stockpile", handlers.CreateStockpileHandler)
	authorized.POST("/update_item", handlers.InsertOrUpdateItemHandler)
	authorized.GET("/all_stockpiles", handlers.GetAllLocationsHandler)

	clanAdmins := r.Group("/clan", loginmiddleware.ClanAdminAuthHandler)
	clanAdmins.GET("/invitation", loginmiddleware.GenerateInvitationLinkHandler)

	admins := r.Group("/admin", loginmiddleware.AdminAuthHandler)
	admins.POST("/invite_clan", loginmiddleware.GenerateClanAdminInvitationLinkHandler)

	r.Run()
}
