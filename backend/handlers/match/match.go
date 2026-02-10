package match

import (
	"github.com/toagne/booking_app/types"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	matchRepo types.MatchRepo
}

func NewHandler(repo types.MatchRepo) *Handler {
	return &Handler{matchRepo: repo}
}

func (h *Handler) GetMatchByMatchId(c *gin.Context) {
	matchId, err := strconv.Atoi(c.Param("id"))
	if err != nil {

	}
	match, err := h.matchRepo.GetMatchByMatchId(matchId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, match)
}

func (h *Handler) GetMatchesByTeam(c *gin.Context) {
	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {

	}
	matches, err := h.matchRepo.GetMatchesByTeam(teamId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, matches)
}

func (h *Handler) GetMatchesByMatchday(c *gin.Context) {
	matchday := "Matchday " + c.Param("id")
	matches, err := h.matchRepo.GetMatchesByMatchday(matchday)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, matches)
}