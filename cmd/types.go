package main

import (
	"context"
	"net/url"
	"time"
)

type RouteInfo struct {
	Service    string
	Namespace  string
	Env        string
	BaseDomain string
	IsValid    bool
}

type ServiceInfo struct {
	Service   string
	Namespace string
	TargetURL *url.URL
	Healthy   bool
	LastCheck time.Time
	TTL       time.Time
}

type HostParser interface {
	ParseHostname(hostname string) RouteInfo
}

type ServiceResolver interface {
	ResolveService(ctx context.Context, route RouteInfo) (*ServiceInfo, error)
	HealthCheck(ctx context.Context, service *ServiceInfo) bool
}

/**
# Test the health endpoint
curl http://localhost:8080/health

# Test the routes endpoint
curl http://localhost:8080/routes

# Test proxy functionality with proper hostname
curl -H "Host: test-service.default.svc.dev.code-craft.in" http://localhost:8080/
*/
