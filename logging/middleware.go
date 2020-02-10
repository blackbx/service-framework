package logging

import (
	"fmt"
	"net/http"

	"github.com/BlackBX/service-framework/dependency"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// NewResponseLogger returns you an instance of a *ResponseLogger
func NewResponseLogger(w http.ResponseWriter) *ResponseLogger {
	return &ResponseLogger{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}
}

// ResponseLogger is a ResponseWriter that is able to log the
// status code of the response
type ResponseLogger struct {
	http.ResponseWriter
	Status int
}

// WriteHeader intercepts the call to the base ResponseWriter and logs the
// status code sent
func (l *ResponseLogger) WriteHeader(statusCode int) {
	l.Status = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

// NewMiddleware returns you a new instance of the Logger middleware
func NewMidlleware(logger *zap.Logger, config dependency.ConfigGetter) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("host", r.Host),
				zap.String("path", path(r)),
				zap.String("protocol", r.Proto),
				zap.Int64("request.content-length", r.ContentLength),
			}
			fields = append(fields, requestHeaders(r, config.GetStringSlice("excluded-headers"))...)
			fields = append(fields, queryParams(r)...)
			responseLogger := NewResponseLogger(rw)
			handler.ServeHTTP(responseLogger, r)
			fields = append(fields, zap.Int("status-code", responseLogger.Status))
			fields = append(fields, responseHeaders(responseLogger.Header())...)
			logger.Info("request log", fields...)
		})
	}
}

func path(r *http.Request) string {
	path := r.URL.Path
	muxPath, err := mux.
		CurrentRoute(r).
		GetPathTemplate()
	if err != nil {
		return path
	}
	return muxPath
}

func responseHeaders(headers http.Header) []zap.Field {
	fields := make([]zap.Field, 0, len(headers))
	for header := range headers {
		headerName := fmt.Sprintf("response.header.%s", header)
		field := zap.String(headerName, headers.Get(header))
		fields = append(fields, field)
	}
	return fields
}

func requestHeaders(r *http.Request, excludedHeaders []string) []zap.Field {
	headers := map[string]struct{}{}
	for _, header := range excludedHeaders {
		headers[header] = struct{}{}
	}
	fields := make([]zap.Field, 0, len(r.Header))
	for header := range r.Header {
		_, ok := headers[header]
		if ok {
			continue
		}
		headerName := fmt.Sprintf("request.header.%s", header)
		field := zap.String(headerName, r.Header.Get(header))
		fields = append(fields, field)
	}
	return fields
}

func queryParams(r *http.Request) []zap.Field {
	fields := make([]zap.Field, 0, len(r.Header))
	for param, values := range r.URL.Query() {
		queryName := fmt.Sprintf("query.%s", param)
		field := zap.Strings(queryName, values)
		fields = append(fields, field)
	}
	return fields
}
