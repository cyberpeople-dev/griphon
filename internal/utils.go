package internal

import (
	"math/rand"
	"github.com/fatih/color"
//	"sync"
//	"time"
)

func RandomUserAgent(userAgents []string) string {
	return userAgents[rand.Intn(len(userAgents))]
}

func RandomPayload(payloads []Payload) Payload {
	return payloads[rand.Intn(len(payloads))]
}

func RandomEndpoint(endpoints []Endpoint) Endpoint {
	return endpoints[rand.Intn(len(endpoints))]
}

func TimeoutLog(method, url string) {
	color.New(color.FgRed).Printf("[TIMEOUT] %s %s\n", method, url)
}

// Simple semaphore for throttling
type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{ch: make(chan struct{}, max)}
}

func (s *Semaphore) Acquire() { s.ch <- struct{}{} }
func (s *Semaphore) Release() { <-s.ch }

func PrintStatus(method, baseURL string, status int) {
	var statusColor *color.Color
	switch {
	case status >= 200 && status < 300:
		statusColor = color.New(color.FgGreen)
	case status >= 300 && status < 400:
		statusColor = color.New(color.FgCyan)
	case status >= 400 && status < 500:
		statusColor = color.New(color.FgYellow)
	default:
		statusColor = color.New(color.FgRed)
	}
	statusColor.Printf("[%d] ", status)
	color.New(color.FgHiBlue).Printf("%s %s\n", method, baseURL)
} 