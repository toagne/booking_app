package auth

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1)
	if err != nil {
		t.Errorf("error creating token: %v", err)
	}
	if token == "" {
		t.Error("expected token to be not empty")
	}
}