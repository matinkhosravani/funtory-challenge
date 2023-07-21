package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/matinkhosravani/funtory-challenge/api/router"
	"github.com/matinkhosravani/funtory-challenge/app"
	"log"
)

func main() {
	app.Boot()
	r := gin.New()
	api := r.Group("api/v1/")
	router.UserRoutes(api)
	// Start server
	err := r.Run(app.GetEnv().ServerAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}
