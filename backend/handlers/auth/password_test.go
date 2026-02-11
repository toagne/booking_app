package auth

import (
	"testing"
)

func TestPasswordHandling(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	t.Run("hashing password", func(t *testing.T) {
		if err != nil {
			t.Errorf("error hashing password: %v", err)
		}
		if hash == "" {
			t.Error("expected hash to be not empty")
		}
		if hash == password {
			t.Error("expected hash to be different from password")
		}
	})

	t.Run("comparing hash and password", func(t *testing.T) {
		if !VerifyPassword(password, hash) {
			t.Errorf("expected hash to match password")
		}
		if VerifyPassword("wrongpassword", hash) {
			t.Errorf("expected password to not match hash")
		}
	})
}