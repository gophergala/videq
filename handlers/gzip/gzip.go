package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	passHandle http.Handler
}

func NewHandler(passHandler http.Handler) *Handler {
	h := new(Handler)
	h.passHandle = passHandler
	return h
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.passHandle.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
	h.passHandle.ServeHTTP(gzw, r)
}
