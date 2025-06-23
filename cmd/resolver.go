package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type kubernetesResolver struct {
	clientset *kubernetes.Clientset
}

func NewKubernetesResolver(clientset *kubernetes.Clientset) ServiceResolver {
	return &kubernetesResolver{
		clientset: clientset,
	}
}

func (r *kubernetesResolver) ResolveService(ctx context.Context, route RouteInfo) (*ServiceInfo, error) {
	service, err := r.clientset.CoreV1().
		Services(route.Namespace).
		Get(ctx, route.Service, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	targetURL := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s.%s.svc.cluster.local:80", route.Service, route.Namespace),
	}

	if len(service.Spec.Ports) > 0 {
		port := service.Spec.Ports[0].Port
		targetURL.Host = fmt.Sprintf("%s.%s.svc.cluster.local:%d", route.Service, route.Namespace, port)

		if port == 443 || service.Spec.Ports[0].Name == "https" {
			targetURL.Scheme = "https"
		}
	}

	serviceInfo := &ServiceInfo{
		Service:   route.Service,
		Namespace: route.Namespace,
		TargetURL: targetURL,
		Healthy:   true, // We'll implement proper health checks later
		LastCheck: time.Now(),
		TTL:       time.Now().Add(5 * time.Minute), // Default TTL
	}

	return serviceInfo, nil
}

func (r *kubernetesResolver) HealthCheck(ctx context.Context, service *ServiceInfo) bool {

	_, err := r.clientset.CoreV1().
		Services(service.Namespace).
		Get(ctx, service.Service, v1.GetOptions{})

	service.LastCheck = time.Now()
	service.Healthy = err == nil

	return service.Healthy
}

// +++++++++++++++++++++++++++++Mock Resolver dev Environment++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type mockResolver struct{}

func NewMockResolver() ServiceResolver {
	return &mockResolver{}
}

func (r *mockResolver) ResolveService(ctx context.Context, route RouteInfo) (*ServiceInfo, error) {
	targetURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:3000",
	}

	serviceInfo := &ServiceInfo{
		Service:   route.Service,
		Namespace: route.Namespace,
		TargetURL: targetURL,
		Healthy:   true,
		LastCheck: time.Now(),
		TTL:       time.Now().Add(5 * time.Minute),
	}

	return serviceInfo, nil
}

func (r *mockResolver) HealthCheck(ctx context.Context, service *ServiceInfo) bool {
	service.LastCheck = time.Now()
	service.Healthy = true
	return service.Healthy
}
