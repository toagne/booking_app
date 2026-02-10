package main

import (
	"github.com/toagne/booking_app/db"
	"github.com/toagne/booking_app/handlers/auth"
	"github.com/toagne/booking_app/handlers/booking"
	"github.com/toagne/booking_app/handlers/match"
	"github.com/toagne/booking_app/handlers/user"
	"github.com/toagne/booking_app/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	newDb := db.InitDb()

	utils.StartEmailWorkers(3)

	router := gin.Default()

	repo := db.NewDbRepo(newDb)

	userHandler := user.NewHandler(repo)
	router.POST("/signup", userHandler.Signup)
	router.POST("/login", userHandler.Login)

	matchHandler := match.NewHandler(repo)
	router.GET("/matches/match/:id", matchHandler.GetMatchByMatchId)
	router.GET("/matches/team/:id", matchHandler.GetMatchesByTeam)
	router.GET("/matches/matchday/:id", matchHandler.GetMatchesByMatchday)

	authWrapper := router.Group("/auth")
	bookingHandler := booking.NewHandler(repo)
	authWrapper.Use(auth.AuthMiddleware())
	{
		authWrapper.POST("/book_match", bookingHandler.BookMatch)
	}

	router.Run(":8080")
}