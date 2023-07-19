package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/matinkhosravani/funtory-challenge/app"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/types/events"
	"log"
	"net/http"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	case *events.Connected:
		fmt.Println("connected", v)
	case *events.LoggedOut:
		fmt.Println("Loggedout", v.Reason)
	case *events.PairSuccess:
		fmt.Println("Paired")
	case *events.QRScannedWithoutMultidevice:
		fmt.Println("Multidevice is not enabled")
	}
}

func main() {
	app.Boot()
	r := gin.New()
	api := r.Group("api/v1/")
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "HelloWorld")
	})
	// Start server
	err := r.Run(app.GetEnv().ServerAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}
