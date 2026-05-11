package traefik_plugin_tenant_token_rewrite

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRewriteMatchesTenantHost(t *testing.T) {
	var gotPath string
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
	        t.Logf("got path: %s", gotPath)

	})

	mw, err := New(context.Background(), next, &Config{
		DomainSuffix:   "example.com",
		SourcePath:     "/openid-connect/token",
		TargetTemplate: "/auth/realms/{tenant}/protocol/openid-connect/token",
	}, "tenant-token-rewrite")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "http://test.example.com/openid-connect/token", nil)
	req.Host = "test.example.com"
	rr := httptest.NewRecorder()


	mw.ServeHTTP(rr, req)

	want := "/auth/realms/test/protocol/openid-connect/token"
	if gotPath != want {
		t.Fatalf("rewritten path = %q, want %q", gotPath, want)
	}
}

func TestNoRewriteWhenHostNotMatchingSuffix(t *testing.T) {
	var gotPath string
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
	        t.Logf("got path: %s", gotPath)
	})

	mw, _ := New(context.Background(), next, &Config{
		DomainSuffix:   "dev7.plainid.cloud",
		SourcePath:     "/openid-connect/token",
		TargetTemplate: "/auth/realms/{tenant}/protocol/openid-connect/token",
	}, "tenant-token-rewrite")

	req := httptest.NewRequest(http.MethodPost, "http://auth.other.cloud/openid-connect/token", nil)
	req.Host = "auth.other.cloud"
	rr := httptest.NewRecorder()
	t.Logf("got path: %s", gotPath)


	mw.ServeHTTP(rr, req)

	if gotPath != "/openid-connect/token" {
		t.Fatalf("path changed unexpectedly: %q", gotPath)
	}
}

func TestNoRewriteWhenPathDoesNotMatch(t *testing.T) {
	var gotPath string
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
	        t.Logf("got path: %s", gotPath)
	})

	mw, _ := New(context.Background(), next, &Config{
		DomainSuffix:   "test.example.com",
		SourcePath:     "/openid-connect/token",
		TargetTemplate: "/auth/realms/{tenant}/protocol/openid-connect/token",
	}, "tenant-token-rewrite")

	req := httptest.NewRequest(http.MethodGet, "http://test.example.com/runtime", nil)
	req.Host = "test.example.com"
	rr := httptest.NewRecorder()
	t.Logf("got path: %s", gotPath)


	mw.ServeHTTP(rr, req)

	if gotPath != "/runtime" {
		t.Fatalf("path changed unexpectedly: %q", gotPath)
	}
}
