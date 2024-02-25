package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buffer     *bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	// Also write to the original ResponseWriter

	n, err := lrw.ResponseWriter.Write(b)
	// Write to the buffer for logging purposes
	lrw.buffer.Write(b)
	return n, err

}

func (s *service) CustomLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now().UTC()

		// Create a log entry with logrus
		logEntry := logrus.WithFields(logrus.Fields{
			"Method":    r.Method,
			"Path":      r.URL.Path,
			"Query":     r.URL.Query(),
			"StartTime": startTime,
		})

		// Capture the request body
		requestBody := readRequestBody(r)
		logEntry = logEntry.WithField("RequestBody", requestBody)

		// Create a custom ResponseWriter to capture the status code and response body
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			buffer:         &bytes.Buffer{},
		}

		// Call the next handler in the chain
		next.ServeHTTP(lrw, r)

		// Capture the response body
		responseBody := lrw.buffer.String()
		logEntry = logEntry.WithFields(logrus.Fields{
			"StatusCode":   lrw.statusCode,
			"ResponseTime": time.Since(startTime),
			"ResponseBody": responseBody,
		})

		// Log the complete request and response details
		logEntry.Info("HTTP Request Processed")

	})

}
func readRequestBody(r *http.Request) map[string]interface{} {
	// Capture the request body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Warn("Error reading request body: ", err)
	}
	// Refill the emptied request body, reset the read position
	r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	// Check if the body is non-empty before attempting to unmarshal
	bodyFields := map[string]interface{}{}
	if len(requestBody) > 1 {
		if err := json.Unmarshal(requestBody, &bodyFields); err != nil {
			logrus.Warn("Error reading request body into map: ", err)
		}
	}
	return bodyFields
}
