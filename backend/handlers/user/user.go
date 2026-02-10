package user

import (
	"github.com/toagne/booking_app/handlers/auth"
	"github.com/toagne/booking_app/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userRepo types.UserRepo
}

func NewHandler(repo types.UserRepo) *Handler {
	return &Handler{userRepo: repo}
}

func (h *Handler) Signup(c *gin.Context) {
	var req types.RegisterAndLoginPayload
	if err := c.BindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	if err := h.userRepo.AddUser(req.Email, hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created succesfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req types.RegisterAndLoginPayload
	if err := c.BindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email"})
		return
	}

	if !auth.VerifyPassword(req.Password, user.HashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := auth.GenerateToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}