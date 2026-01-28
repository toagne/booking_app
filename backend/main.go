package main

import (
	"github.com/toagne/booking_app/db"
	"github.com/toagne/booking_app/handlers"
	"github.com/toagne/booking_app/middleware"
	"github.com/toagne/booking_app/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDb()

	utils.StartEmailWorkers(3)

	router := gin.Default()

	router.GET("/matches/matchday/:id", handlers.GetMatchesByMatchday)
	router.GET("/matches/team/:id", handlers.GetMatchesByTeam)
	router.GET("/matches/match/:id", handlers.GetMatchByMatchId)
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)

	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/book_match", handlers.BookMatch)
	}

	router.Run(":8080")
}