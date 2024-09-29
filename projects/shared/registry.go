package shared

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Registry map[string]string

func (r Registry) Endpoint(name string) string {
	if endpoint, ok := r[name]; ok {
		return endpoint
	}
	return ""
}

func RegistryFromEnv(prefix string, services []string) Registry {
	registry := Registry{}
	for _, service := range services {
		env := fmt.Sprintf("%s_%s_ENDPOINT", strings.ToUpper(prefix), strings.ToUpper(service))
		endpoint := os.Getenv(env)
		if endpoint == "" {
			log.Fatalf("env %s is required", env)
		}
		registry[service] = endpoint
	}
	return registry
}
