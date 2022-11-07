package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler, logger shared.ILogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			next.ServeHTTP(w, r)
		}

		responseWriter := NewLoggingResponseWriter(w)

		requestStart := time.Now()
		next.ServeHTTP(responseWriter, r)
		requestEnd := time.Since(requestStart)

		statusCode := responseWriter.statusCode

		logger.Info(fmt.Sprintf("%s [%s] %s %s %s %d %s %s \n",
			ip,
			requestStart.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL,
			r.Proto,
			statusCode,
			fmt.Sprint(requestEnd),
			r.UserAgent(),
		))
	})
}
