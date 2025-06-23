package main

import (
	"log"
	"strings"
)

type hostParser struct {
	baseDomain        string
	defaultEnv        string
	allowedNamespaces []string
}

func NewHostParser(baseDomain, defaultEnv string, allowedNamespaces []string) HostParser {
	return &hostParser{
		baseDomain:        baseDomain,
		defaultEnv:        defaultEnv,
		allowedNamespaces: allowedNamespaces,
	}
}

func (p *hostParser) ParseHostname(hostname string) RouteInfo {
	parts := strings.Split(hostname, ".")

	if strings.Contains(parts[0], ":") {
		parts[0] = strings.Split(parts[0], ":")[0]
	}

	route := RouteInfo{IsValid: false}

	if len(parts) < 3 {
		log.Printf("Invalid hostname format: %s (too few parts)", hostname)
		return route
	}

	// Parse based on expected pattern: service.namespace.svc.env.base.domain
	if len(parts) >= 5 && parts[2] == "svc" {
		route.Service = parts[0]
		route.Namespace = parts[1]
		route.Env = parts[3]
		route.BaseDomain = strings.Join(parts[4:], ".")
		route.IsValid = true
	} else if len(parts) >= 3 {
		// Fallback: service.namespace.rest-as-domain
		route.Service = parts[0]
		route.Namespace = parts[1]
		route.Env = p.defaultEnv
		route.BaseDomain = strings.Join(parts[2:], ".")
		route.IsValid = true
	}

	if route.IsValid && len(p.allowedNamespaces) > 0 {
		allowed := false
		for _, ns := range p.allowedNamespaces {
			if ns == route.Namespace {
				allowed = true
				break
			}
		}
		if !allowed {
			log.Printf("Namespace %s not in allowed list", route.Namespace)
			route.IsValid = false
		}
	}

	return route
}
