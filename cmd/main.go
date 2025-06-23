package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config := LoadConfig()
	log.Printf("Starting with config: port=%s, base_domain=%s",
		config.Port, config.BaseDomain)

	devMode := os.Getenv("DEV_MODE") == "true"

	var resolver ServiceResolver
	if devMode {
		log.Println("Running in development mode - using mock resolver")
		resolver = NewMockResolver()
	} else {
		k8sConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed to create in-cluster config: %v", err)
		}

		clientset, err := kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			log.Fatalf("Failed to create clientset: %v", err)
		}
		resolver = NewKubernetesResolver(clientset)
	}

	parser := NewHostParser(config.BaseDomain, config.DefaultEnv, config.AllowedNamespaces)

	server := NewServer(config, parser, resolver)

	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      server.Handler(),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		log.Printf("Starting reverse proxy server on :%s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
