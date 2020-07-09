package examples

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/maykonlf/go-devkit/log"
)

type ResponseWriter struct {
	http.ResponseWriter
	status      int
	body        []byte
	wroteHeader bool
}

func (rw *ResponseWriter) Status() int {
	return rw.status
}

func (rw *ResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *ResponseWriter) Write(body []byte) (int, error) {
	rw.body = body
	return rw.ResponseWriter.Write(body)
}

func wrapResponseWriter(w http.ResponseWriter, defaultStatus int) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, status: defaultStatus}
}

func HTTPLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(r.Method, "PRI") {
			next.ServeHTTP(w, r)
			return
		}

		httpLoggingMiddleware(w, r, next)
	}

	return http.HandlerFunc(fn)
}

func httpLoggingMiddleware(w http.ResponseWriter, r *http.Request, next http.Handler) {
	start := time.Now()
	requestBody := readRequestBody(r)
	responseWrapped := wrapResponseWriter(w, 200)

	next.ServeHTTP(responseWrapped, r)
	duration := time.Since(start)

	status := responseWrapped.Status()
	log.Info(fmt.Sprintf("%s %s %s %d", r.Method, r.URL.Path, r.Proto, status),
		"http_response_status", status,
		"http_request_method", r.Method, "http_request_path", r.URL.Path,
		"http_request_args", r.URL.Query(), "http_request_headers", r.Header,
		"http_request_body", string(requestBody), "duration", duration.String(),
	)
}

func readRequestBody(r *http.Request) []byte {
	requestBody, _ := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	return requestBody
}
