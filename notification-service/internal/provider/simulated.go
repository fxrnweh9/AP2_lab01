package provider

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type SimulatedEmailSender struct{}

func NewSimulatedEmailSender() *SimulatedEmailSender {
	return &SimulatedEmailSender{}
}

func (s *SimulatedEmailSender) Send(to, subject, body string) error {
	latency := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(latency)

	if rand.Float32() < 0.8 {
		return errors.New("simulated provider error: temporary failure")
	}

	log.Printf("[SimulatedEmail] TO=%s SUBJECT=%s", to, subject)
	return nil
}
