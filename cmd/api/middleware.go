package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})

}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type limit struct {
		mu          sync.Mutex
		rateLimiter *rate.Limiter
		lastSeen    time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*limit)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		var found bool

		mu.Lock()
		if _, found = clients[ip]; !found {
			clients[ip] = &limit{mu: sync.Mutex{}, rateLimiter: rate.NewLimiter(2, 4)}
		}
		mu.Unlock()

		currentLimit := clients[ip]
		currentLimit.mu.Lock()
		if !currentLimit.rateLimiter.Allow() {
			currentLimit.mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}
		currentLimit.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
