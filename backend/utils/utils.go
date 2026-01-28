package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Email struct {
	To		string
	Subject	string
	Body	string
}

var EmailQueue = make(chan Email, 10)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userId int) (string, error) {
	claims := jwt.RegisteredClaims {
		ExpiresAt:	jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		IssuedAt:	jwt.NewNumericDate(time.Now()),
		Subject:	strconv.Itoa(userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func sendEmail(email Email) {
	time.Sleep(5 * time.Second)
	log.Printf("\n\n####SENDING NEW EMAIL####\n\nTo: %v\n\nSubject: %v\n\n%v\n\n#########################",
		email.To, email.Subject, email.Body,
	)
}

func StartEmailWorkers(n int) {
	for i := range n {
		go func(workerId int) {
			for email := range EmailQueue {
				sendEmail(email)
			}
		} (i)
	}
}