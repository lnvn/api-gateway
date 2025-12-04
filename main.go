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
)

func main() {
	// 1. Load Configuration
	// For simplicity, we'll create a default config if the file doesn't exist or just hardcode for now for the first run,
	// but let's assume a config.json exists or we create one.
	// Let's create a sample config in memory for this step to ensure it runs without external dependencies immediately,
	// or better, let's write a config.json file in the next step.

	// For now, let's just print a message that we are starting.
	fmt.Println("Starting API Gateway...")

	// Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Set up Routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

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
