package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type Server struct {
	config   *Config
	parser   HostParser
	resolver ServiceResolver
}

func NewServer(config *Config, parser HostParser, resolver ServiceResolver) *Server {
	return &Server{
		config:   config,
		parser:   parser,
		resolver: resolver,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/healthz", s.healthHandler)
	mux.HandleFunc("/routes", s.routesHandler)
	mux.HandleFunc("/", s.proxyHandler)

	// Wrap the mux with a logging middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request started: %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

		mux.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("Request completed: %s %s %s (took %v)", r.RemoteAddr, r.Method, r.URL.Path, duration)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) routesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"config": map[string]interface{}{
			"base_domain":        s.config.BaseDomain,
			"default_env":        s.config.DefaultEnv,
			"allowed_namespaces": s.config.AllowedNamespaces,
		},
		"patterns": []string{
			"{service}.{namespace}.svc.{env}.{base_domain}",
			"{service}.{namespace}.{domain}",
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
	route := s.parser.ParseHostname(r.Host)
	if !route.IsValid {
		log.Printf("Invalid route from hostname: %s", r.Host)
		http.Error(w, "Invalid hostname format", http.StatusBadRequest)
		return
	}
	log.Printf("Received request for route: %s", route)
	if route.Namespace == "health" {
		s.healthHandler(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceInfo, err := s.resolver.ResolveService(ctx, route)
	if err != nil {
		log.Printf("Failed to resolve service %s.%s: %v", route.Service, route.Namespace, err)
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if !serviceInfo.Healthy {
		log.Printf("Service %s.%s is marked as unhealthy", route.Service, route.Namespace)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(serviceInfo.TargetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for %s -> %s: %v", r.Host, serviceInfo.TargetURL.String(), err)
		http.Error(w, "Service temporarily unavailable", http.StatusBadGateway)
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Proxy-Service", route.Service)
		resp.Header.Set("X-Proxy-Namespace", route.Namespace)
		resp.Header.Set("X-Proxy-Target", serviceInfo.TargetURL.String())
		return nil
	}
	log.Printf("Proxying %s %s -> %s", r.Method, r.URL.Path, serviceInfo.TargetURL.String())
	proxy.ServeHTTP(w, r)
}
