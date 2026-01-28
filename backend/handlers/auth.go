package handlers

import (

	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/toagne/booking_app/db"
	"github.com/toagne/booking_app/utils"
	"fmt"
)

func BookMatch(c *gin.Context) {
	userId := c.GetInt("userId")
	
	var req struct {
		GameId		int `json:"gameId"`
		NOfTickets	int `json:"tickets"`
	}
	if err := c.BindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.GetMatchByMatchId(req.GameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingId, err := db.AddBooking(userId, req.GameId, req.NOfTickets)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := db.GetBookingInfo(bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var email utils.Email
	email.To = info.Username
	email.Subject = fmt.Sprintf("Booking Confirmation %v - %v" , info.Match.Team1, info.Match.Team2)
	email.Body = fmt.Sprintf("Here you can find the booking details:\nDate: %v\nTime: %v\nGame: %v - %v\nN of tickets: %v", info.Match.Date, info.Match.Time, info.Match.Team1, info.Match.Team2, info.Tickets)

	utils.EmailQueue <- email

	c.JSON(http.StatusOK, *info)
}