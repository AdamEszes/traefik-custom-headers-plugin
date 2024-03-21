package traefik_plugin_rename_headers

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
)

// Rename holds one rename configuration.
type renameData struct {
	ExistingHeaderName string `json:"existingHeaderName"`
	NewHeaderName      string `json:"newHeaderName"`
}

// Config holds the plugin configuration.
type Config struct {
	RenameData []renameData `json:"renameData"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// New creates and returns a new rewrite body plugin instance.
type renameHeaders struct {
	name    string
	next    http.Handler
	renames []renameData
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &renameHeaders{
		name:    name,
		next:    next,
		renames: config.RenameData,
	}, nil
}

func (r *renameHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedWriter := &responseWriter{
		writer:          rw,
		headersToRename: r.renames,
	}

	r.next.ServeHTTP(wrappedWriter, req)
}

type responseWriter struct {
	writer          http.ResponseWriter
	headersToRename []renameData
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	for _, headerToRename := range r.headersToRename {
		headerValues := r.writer.Header().Values(headerToRename.ExistingHeaderName)

		if len(headerValues) == 0 {
			continue
		}

		r.writer.Header().Del(headerToRename.ExistingHeaderName)
		r.writer.Header()[headerToRename.NewHeaderName] = headerValues
	}

	r.writer.WriteHeader(statusCode)
}

func (r *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.writer.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", r.writer)
	}

	return hijacker.Hijack()
}

func (r *responseWriter) Flush() {
	if flusher, ok := r.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
