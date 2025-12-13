package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/config"
	"api-gateway/middleware"
	"api-gateway/proxy"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("Starting API Gateway...")

	// Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Set up Routes
	mux := http.NewServeMux()

	// Liveness endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Readiness endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	for _, route := range cfg.Routes {
		// Create a proxy for each route
		p, err := proxy.NewProxy(route.BackendURL)
		if err != nil {
			log.Fatalf("Failed to create proxy for %s: %v", route.BackendURL, err)
		}

		// Handle the route
		// Note: http.StripPrefix might be needed depending on how we want to forward the path.
		// For this simple example, we'll just forward everything under the path.
		handler := proxy.ProxyHandler(p)

		// Wrap with middleware
		loggedHandler := middleware.Logger(handler)

		mux.Handle(route.Path, http.StripPrefix(route.Path, loggedHandler))

		fmt.Printf("Setup route: %s -> %s\n", route.Path, route.BackendURL)
	}

	// 3. Start Server with Graceful Shutdown
	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("API Gateway listening on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
