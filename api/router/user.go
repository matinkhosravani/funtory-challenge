package router

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/funtory-challenge/api/handler"
	"github.com/matinkhosravani/funtory-challenge/client/whatsapp"
	"github.com/matinkhosravani/funtory-challenge/repository"
)

func UserRoutes(r *gin.RouterGroup) {
	setUpConnectWhatsappCheckRoute(r)
}

func setUpConnectWhatsappCheckRoute(r *gin.RouterGroup) {
	userRepo := repository.NewUserRepository()
	h := handler.ConnectWhatsappHandler{
		UserRepo: userRepo,
		Client:   whatsapp.NewWhatsmeow(userRepo),
	}

	r.GET("/connect/:user_id", h.Handle)
}
