package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/toagne/booking_app/types"
)

func TestUserHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userRepo := &mockUserRepo{}
	handler := NewHandler(userRepo)
	t.Run("should fail if signing up with no email", func(t *testing.T) {
		router := gin.Default()
		router.POST("/signup", handler.Signup)
		body := `{
		"email":"",
		"password":"password"
		}`
		req, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBufferString(body),
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail if signing up with invalid email", func(t *testing.T) {
		router := gin.Default()
		router.POST("/signup", handler.Signup)
		body := `{
		"email":"aaa",
		"password":"password"
		}`
		req, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBufferString(body),
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail if signing up with password lenght < 8", func(t *testing.T) {
		router := gin.Default()
		router.POST("/signup", handler.Signup)
		body := `{
		"email":"aaa@aaa.aaa",
		"password":"123"
		}`
		req, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBufferString(body),
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should register the user correctly", func(t *testing.T) {
		router := gin.Default()
		router.POST("/signup", handler.Signup)
		body := `{
		"email":"aaa@aaa.aaa",
		"password":"password"
		}`
		req, err := http.NewRequest(
			http.MethodPost,
			"/signup",
			bytes.NewBufferString(body),
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserRepo struct {}

func (m *mockUserRepo) AddUser(email, hashedPassword string) error {
	return nil
}

func (m *mockUserRepo) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}