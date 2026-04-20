package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
)

const (
	immutableAssetsCacheControl = "public, max-age=31536000, immutable"
	rootStaticCacheControl      = "public, max-age=3600"
	htmlCacheControl            = "no-cache"
	healthCacheControl          = "no-store"
)

type Handler struct {
	files     fs.FS
	indexHTML []byte
}

func New(files fs.FS) (http.Handler, error) {
	indexHTML, err := fs.ReadFile(files, "dist/index.html")
	if err != nil {
		return nil, fmt.Errorf("read embedded index.html: %w", err)
	}

	handler := &Handler{
		files:     files,
		indexHTML: indexHTML,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/livez", healthHandler)
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/", handler.serveHTTP)

	return withSecurityHeaders(mux), nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body := []byte("ok")
	w.Header().Set("Cache-Control", healthCacheControl)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = w.Write(body)
	}
}

func (h *Handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath, cacheControl, fallbackToIndex := resolvePath(r.URL.Path)
	if fallbackToIndex {
		h.writeStatic(w, r, http.StatusOK, "dist/index.html", h.indexHTML, htmlCacheControl)
		return
	}

	body, err := fs.ReadFile(h.files, filePath)
	if err != nil {
		if isSPAPath(r.URL.Path) {
			h.writeStatic(w, r, http.StatusOK, "dist/index.html", h.indexHTML, htmlCacheControl)
			return
		}
		http.NotFound(w, r)
		return
	}

	h.writeStatic(w, r, http.StatusOK, filePath, body, cacheControl)
}

func resolvePath(requestPath string) (filePath string, cacheControl string, fallbackToIndex bool) {
	cleaned := path.Clean("/" + requestPath)
	if cleaned == "/" {
		return "dist/index.html", htmlCacheControl, false
	}

	relative := strings.TrimPrefix(cleaned, "/")
	if strings.HasPrefix(relative, "assets/") {
		return path.Join("dist", relative), immutableAssetsCacheControl, false
	}
	if !isSPAPath(cleaned) {
		return path.Join("dist", relative), rootStaticCacheControl, false
	}
	return "dist/index.html", htmlCacheControl, true
}

func isSPAPath(requestPath string) bool {
	base := path.Base(requestPath)
	return !strings.Contains(base, ".")
}

func (h *Handler) writeStatic(w http.ResponseWriter, r *http.Request, status int, name string, body []byte, cacheControl string) {
	contentType := mime.TypeByExtension(path.Ext(name))
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	w.Header().Set("Cache-Control", cacheControl)
	w.Header().Set("Content-Type", contentType)
	w.Header().Add("Vary", "Accept-Encoding")

	responseBody := body
	if acceptsGzip(r.Header.Get("Accept-Encoding")) && isCompressible(contentType) && len(body) > 256 {
		var compressed bytes.Buffer
		writer, err := gzip.NewWriterLevel(&compressed, gzip.BestSpeed)
		if err == nil {
			_, _ = writer.Write(body)
			_ = writer.Close()
			responseBody = compressed.Bytes()
			w.Header().Set("Content-Encoding", "gzip")
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = w.Write(responseBody)
	}
}

func acceptsGzip(acceptEncoding string) bool {
	return strings.Contains(acceptEncoding, "gzip")
}

func isCompressible(contentType string) bool {
	return strings.HasPrefix(contentType, "text/") ||
		strings.HasPrefix(contentType, "application/javascript") ||
		strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(contentType, "application/xml") ||
		strings.HasPrefix(contentType, "image/svg+xml")
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Permissions-Policy", "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; "+
				"connect-src 'self' https://pagead2.googlesyndication.com https://googleads.g.doubleclick.net https://tpc.googlesyndication.com https://ep1.adtrafficquality.google https://cloudflareinsights.com; "+
				"img-src 'self' data: https://pagead2.googlesyndication.com https://googleads.g.doubleclick.net https://tpc.googlesyndication.com; "+
				"style-src 'self' 'unsafe-inline'; "+
				"script-src 'self' https://pagead2.googlesyndication.com https://static.cloudflareinsights.com; "+
				"frame-src https://googleads.g.doubleclick.net https://tpc.googlesyndication.com; "+
				"object-src 'none'; base-uri 'self'; frame-ancestors 'none';",
		)
		next.ServeHTTP(w, r)
	})
}
