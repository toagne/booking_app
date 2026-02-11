package match

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/toagne/booking_app/types"
)

func TestMatchHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	matchRepo := &mockMatchRepo{}
	handler := NewHandler(matchRepo)

	t.Run("should fail match_id not convertible by Atoi", func(t *testing.T) {
		router := gin.Default()
		matchId := "a"
		router.POST("/matches/match/:id", handler.GetMatchByMatchId)
		req, err := http.NewRequest(
			http.MethodPost,
			"/matches/match/" + matchId,
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}

type mockMatchRepo struct {}

func (m *mockMatchRepo) GetMatchByMatchId(matchId int) (*types.Match, error) {
	return &types.Match{}, nil
}

func (m *mockMatchRepo) GetMatchesByTeam(teamId int) (*[]types.Match, error) {
	return &[]types.Match{}, nil
}

func (m *mockMatchRepo) GetMatchesByMatchday(matchday string) (*[]types.Match, error) {
	return &[]types.Match{}, nil
}

