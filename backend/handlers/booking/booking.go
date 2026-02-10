package booking

import (

	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/toagne/booking_app/utils"
	"github.com/toagne/booking_app/types"
	"fmt"
)

type Handler struct {
	bookingRepo types.BookingRepo
}

func NewHandler(repo types.BookingRepo) *Handler {
	return &Handler{bookingRepo: repo}
}

func (h *Handler) BookMatch(c *gin.Context) {
	userId := c.GetInt("userId")
	
	var req types.BookingPayload
	if err := c.BindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.bookingRepo.GetMatchByMatchId(req.GameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingId, err := h.bookingRepo.AddBooking(userId, req.GameId, req.NOfTickets)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := h.bookingRepo.GetBookingInfo(bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var email types.Email
	email.To = info.Username
	email.Subject = fmt.Sprintf("Booking Confirmation %v - %v" , info.Match.Team1, info.Match.Team2)
	email.Body = fmt.Sprintf("Here you can find the booking details:\nDate: %v\nTime: %v\nGame: %v - %v\nN of tickets: %v", info.Match.Date, info.Match.Time, info.Match.Team1, info.Match.Team2, info.Tickets)

	utils.EmailQueue <- email

	c.JSON(http.StatusOK, *info)
}