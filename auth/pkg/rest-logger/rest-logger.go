package restLogger

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// LogRequest is a struct that contains the request information.
type LogRequest struct {
	Method       string
	URL          string
	RemoteAddr   string
	StatusCode   int
	Latency      time.Duration
	RequestBody  string
	ResponseBody string
}

// LogMiddleware is a middleware that logs the request and response information.
func LogMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			requestBody, err := readRequestBody(r)
			if err != nil {
				logger.Error("Failed to read request body", slog.String("error", err.Error()))
			}

			wrappedWriter := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(wrappedWriter, r)

			responseBody, err := readResponseBody(wrappedWriter)
			if err != nil {
				logger.Error("Failed to read response body", slog.String("error", err.Error()))
			}

			logRequest(logger, r, wrappedWriter.statusCode, time.Since(startTime), requestBody, responseBody)
		})
	}
}

func readRequestBody(r *http.Request) (string, error) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	r.Body = io.NopCloser(bytes.NewReader(requestBody))

	return string(requestBody), nil
}

// readResponseBody reads the response body.
func readResponseBody(rw *responseWriter) (string, error) {
	responseBody, err := io.ReadAll(&rw.buf)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

// logRequest logs the request information.
func logRequest(logger *slog.Logger, r *http.Request, statusCode int, latency time.Duration, requestBody, responseBody string) {
	logger.Info("HTTP request",
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()),
		slog.String("remote_addr", r.RemoteAddr),
		slog.Int("status_code", statusCode),
		slog.Duration("latency", latency),
		slog.String("request_body", requestBody),
		slog.String("response_body", responseBody),
	)
}

// responseWriter is a wrapper for the http.ResponseWriter that captures the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

// WriteHeader writes the status code and captures it.
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write writes the response body and captures it.
func (rw *responseWriter) Write(b []byte) (int, error) {
	_, err := rw.buf.Write(b)
	if err != nil {
		return 0, err
	}
	return rw.ResponseWriter.Write(b)
}
