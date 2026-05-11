package traefik_plugin_tenant_token_rewrite

import (
	"context"
	"net"
	"net/http"
	"strings"
)

type Config struct {
	DomainSuffix string `json:"domainSuffix,omitempty"`
	SourcePath string `json:"sourcePath,omitempty"`
	TargetTemplate string `json:"targetTemplate,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		SourcePath:     "/openid-connect/token",
		TargetTemplate: "/auth/realms/{tenant}/protocol/openid-connect/token",
	}
}

type Middleware struct {
	next http.Handler
	cfg  *Config
	name string
}

func New(_ context.Context, next http.Handler, cfg *Config, name string) (http.Handler, error) {
	if cfg == nil {
		cfg = CreateConfig()
	}
	if cfg.SourcePath == "" {
		cfg.SourcePath = "/openid-connect/token"
	}
	if cfg.TargetTemplate == "" {
		cfg.TargetTemplate = "/auth/realms/{tenant}/protocol/openid-connect/token"
	}

	return &Middleware{
		next: next,
		cfg:  cfg,
		name: name,
	}, nil
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == m.cfg.SourcePath {
		host := normalizeHost(req.Host)
		tenant, ok := extractTenant(host, m.cfg.DomainSuffix)
		if ok && tenant != "" {
			req.URL.Path = strings.ReplaceAll(m.cfg.TargetTemplate, "{tenant}", tenant)
			// Clear RequestURI so upstream uses updated URL.Path.
			req.RequestURI = ""
		}
	}

	m.next.ServeHTTP(rw, req)
}

func normalizeHost(host string) string {
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	return host
}

func extractTenant(host, suffix string) (string, bool) {
	suffix = strings.TrimPrefix(strings.TrimSpace(strings.ToLower(suffix)), ".")
	if suffix == "" {
		return "", false
	}

	dotSuffix := "." + suffix
	if !strings.HasSuffix(host, dotSuffix) {
		return "", false
	}

	left := strings.TrimSuffix(host, dotSuffix)
	if left == "" {
		return "", false
	}

	// Keep only first label as tenant.
	if i := strings.IndexByte(left, '.'); i >= 0 {
		left = left[:i]
	}
	return left, left != ""
}
