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
	h := handler.ConnectWhatsappHandler{
		UserRepo: repository.NewUserRepository(),
		Client:   whatsapp.NewWhatsmeow(),
	}

	r.GET("/connect/:user_id", h.Handle)
}
