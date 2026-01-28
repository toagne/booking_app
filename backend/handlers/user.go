package handlers

import (
	"github.com/toagne/booking_app/db"
	"github.com/toagne/booking_app/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMatchByMatchId(c *gin.Context) {
	matchId, err := strconv.Atoi(c.Param("id"))
	if err != nil {

	}
	match, err := db.GetMatchByMatchId(matchId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, match)
}

func GetMatchesByTeam(c *gin.Context) {
	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {

	}
	matches, err := db.GetMatchesByTeam(teamId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, matches)
}

func GetMatchesByMatchday(c *gin.Context) {
	matchday := "Matchday " + c.Param("id")
	matches, err := db.GetMatchesByMatchday(matchday)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, matches)
}

func Signup(c *gin.Context) {
	var req struct {
		Email		string `json:"email" binding:"required,email"`
		Password	string `json:"password" binding:"required,min=8"`
	}
	if err := c.BindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
	}

	if err := db.AddUser(req.Email, hash); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created succesfully"})
}

func Login(c *gin.Context) {
	var req struct {
		Email		string `json:"email" binding:"required,email"`
		Password	string `json:"password" binding:"required,min=8"`
	}
	if err := c.BindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email"})
		return
	}

	log.Printf("username after GetUserByEmail in Login: %v\n", user.Username)
	log.Printf("hashed password after GetUserByEmail in Login: %v\n", user.HashedPassword)
	if !utils.VerifyPassword(req.Password, user.HashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}