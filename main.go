package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	routes "daucu/routes"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()

	//Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://www.piesocket.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Set 10GB max upload size
	app.MaxMultipartMemory = 10000000000

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Daucu",
		})
	})

	// Routes
	auth := app.Group("/auth")
	{
		auth.POST("/register", routes.Register)
		auth.POST("/login", routes.Login)
		auth.POST("/profile", routes.Profile)
	}

	//Print
	fmt.Println("Server is running on port http://localhost:9000")
	//Run on port 9000
	app.Run(":9000")

}
