package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"
)

func TestServesIndexAndSPAFallback(t *testing.T) {
	t.Parallel()

	handler := mustHandler(t)

	root := httptest.NewRequest(http.MethodGet, "/", nil)
	rootRecorder := httptest.NewRecorder()
	handler.ServeHTTP(rootRecorder, root)

	if rootRecorder.Code != http.StatusOK {
		t.Fatalf("root status = %d, want %d", rootRecorder.Code, http.StatusOK)
	}
	if got := rootRecorder.Header().Get("Cache-Control"); got != htmlCacheControl {
		t.Fatalf("root cache-control = %q, want %q", got, htmlCacheControl)
	}

	spa := httptest.NewRequest(http.MethodGet, "/projects/checkboxes", nil)
	spaRecorder := httptest.NewRecorder()
	handler.ServeHTTP(spaRecorder, spa)

	if spaRecorder.Code != http.StatusOK {
		t.Fatalf("spa status = %d, want %d", spaRecorder.Code, http.StatusOK)
	}
	if body := spaRecorder.Body.String(); !strings.Contains(body, "naterpatater.com") {
		t.Fatalf("spa body = %q, want index.html", body)
	}
}

func TestServesAssetWithImmutableCachingAndGzip(t *testing.T) {
	t.Parallel()

	handler := mustHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("asset status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Cache-Control"); got != immutableAssetsCacheControl {
		t.Fatalf("asset cache-control = %q, want %q", got, immutableAssetsCacheControl)
	}
	if got := recorder.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("asset content-encoding = %q, want gzip", got)
	}
}

func TestServesRootStaticFileWithShortCache(t *testing.T) {
	t.Parallel()

	handler := mustHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("favicon status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Cache-Control"); got != rootStaticCacheControl {
		t.Fatalf("favicon cache-control = %q, want %q", got, rootStaticCacheControl)
	}
}

func TestRejectsUnsupportedMethod(t *testing.T) {
	t.Parallel()

	handler := mustHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("post status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
	if got := recorder.Header().Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("allow header = %q, want %q", got, "GET, HEAD")
	}
}

func TestMissingAssetReturnsNotFound(t *testing.T) {
	t.Parallel()

	handler := mustHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/assets/missing.js", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("missing asset status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestAddsSecurityHeadersAndHealthEndpoints(t *testing.T) {
	t.Parallel()

	handler := mustHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if body, err := io.ReadAll(recorder.Body); err != nil || string(body) != "ok" {
		t.Fatalf("health body = %q, err = %v, want ok", string(body), err)
	}
	if got := recorder.Header().Get("Content-Security-Policy"); got == "" {
		t.Fatal("expected content-security-policy header")
	}
	if got := recorder.Header().Get("Cache-Control"); got != healthCacheControl {
		t.Fatalf("health cache-control = %q, want %q", got, healthCacheControl)
	}
}

func TestIndexScriptNonceMatchesStrictCSP(t *testing.T) {
	t.Parallel()

	handler := mustHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	csp := recorder.Header().Get("Content-Security-Policy")
	match := regexp.MustCompile(`script-src 'nonce-([^']+)'`).FindStringSubmatch(csp)
	if len(match) != 2 {
		t.Fatalf("content-security-policy = %q, want script nonce", csp)
	}
	if !strings.Contains(csp, "'strict-dynamic'") {
		t.Fatalf("content-security-policy = %q, want strict-dynamic", csp)
	}
	if strings.Contains(csp, "pagead2.googlesyndication.com") {
		t.Fatalf("content-security-policy = %q, want no AdSense domain allowlist", csp)
	}

	wantScript := `<script nonce="` + match[1] + `" src="/assets/app.js"></script>`
	if body := recorder.Body.String(); !strings.Contains(body, wantScript) {
		t.Fatalf("index body = %q, want nonce-bearing script %q", body, wantScript)
	}
}

func mustHandler(t *testing.T) http.Handler {
	t.Helper()

	handler, err := New(fstest.MapFS{
		"dist/index.html":       {Data: []byte(`<!doctype html><title>naterpatater.com</title><script src="/assets/app.js"></script>`)},
		"dist/assets/app.js":    {Data: []byte(strings.Repeat("console.log('site');", 40))},
		"dist/favicon.ico":      {Data: []byte("icon")},
		"dist/site.webmanifest": {Data: []byte(`{"name":"naterpatater.com"}`)},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return handler
}
