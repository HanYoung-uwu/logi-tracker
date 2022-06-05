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
	basePath := r.Group("/api")
	basePath.POST("/login", loginmiddleware.LoginHandler)
	basePath.POST("/register", loginmiddleware.CreateUserFromInvitationLinkHandler)
	basePath.POST("/admin/create_admin", loginmiddleware.CreateAdminOnFirstRequestHandler)

	authorized := basePath.Group("/user", loginmiddleware.DefaultAuthHandler)
	authorized.GET("/all_items", handlers.GetAllItemsHandler)
	authorized.POST("/create_stockpile", handlers.CreateStockpileHandler)
	authorized.POST("/update_item", handlers.InsertOrUpdateItemHandler)
	authorized.GET("/all_stockpiles", handlers.GetAllLocationsHandler)
	authorized.GET("/delete_stockpile", handlers.DeleteStockpileHandler)

	clanAdmins := basePath.Group("/clan", loginmiddleware.ClanAdminAuthHandler)
	clanAdmins.GET("/invitation", loginmiddleware.GenerateInvitationLinkHandler)

	admins := basePath.Group("/admin", loginmiddleware.AdminAuthHandler)
	admins.POST("/invite_clan", loginmiddleware.GenerateClanAdminInvitationLinkHandler)

	r.Run()
}
