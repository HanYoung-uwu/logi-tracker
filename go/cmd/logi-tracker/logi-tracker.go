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
	basePath.GET("/logout", loginmiddleware.LogoutHandler)
	basePath.POST("/login", loginmiddleware.LoginHandler)
	basePath.POST("/register", loginmiddleware.CreateUserFromInvitationLinkHandler)
	basePath.POST("/check_name", handlers.CheckAccountNameExistHandler)
	basePath.POST("/check_clan", handlers.CheckClanExistHandler)
	basePath.POST("/admin/create_admin", loginmiddleware.CreateAdminOnFirstRequestHandler)
	basePath.GET("/invite_info", loginmiddleware.InviteAccountInfoHandler)

	authorized := basePath.Group("/user", loginmiddleware.DefaultAuthHandler)
	authorized.GET("/info", handlers.GetBasicAccountInfoHandler)
	authorized.GET("/all_items", handlers.GetAllItemsHandler)
	authorized.POST("/create_stockpile", handlers.CreateStockpileHandler)
	authorized.POST("/update_item", handlers.InsertOrUpdateItemHandler)
	authorized.GET("/all_stockpiles", handlers.GetAllLocationsHandler)
	authorized.POST("/delete_stockpile", handlers.DeleteStockpileHandler)
	authorized.POST("/delete_item", handlers.DeleteItemHandler)
	authorized.POST("/set_item", handlers.SetItemHandler)
	authorized.POST("/refresh_stockpile", handlers.RefreshStockpileHandler)
	authorized.GET("/history", handlers.GetClanHistoryHandler)

	clanAdmins := basePath.Group("/clan", loginmiddleware.ClanAdminAuthHandler)
	clanAdmins.GET("/invitation", loginmiddleware.GenerateInvitationLinkHandler)
	clanAdmins.GET("/member_info", loginmiddleware.GetClanAccountInfoHandler)
	clanAdmins.POST("/kick_member", loginmiddleware.KickClanMemberHandler)
	clanAdmins.POST("/promote_member", loginmiddleware.PromoteClanMemberHandler)

	admins := basePath.Group("/admin", loginmiddleware.AdminAuthHandler)
	admins.GET("/invite_clan", loginmiddleware.GenerateClanAdminInvitationLinkHandler)
	r.Run()
}
