package utils

import (
	"log"
	"time"
	"github.com/toagne/booking_app/types"
)

var EmailQueue = make(chan types.Email, 10)

func sendEmail(email types.Email) {
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