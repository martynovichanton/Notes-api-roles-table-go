package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	start := time.Now()
	// 	uri := r.RequestURI
	// 	method := r.Method

	// 	next.ServeHTTP(w, r)

	// 	duration := time.Since(start)
	// 	log.Println(uri, method, duration)

	// 	// next(w, r)
	// }

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)
		log.Printf("%s %s %s %d %s", r.Method, r.RequestURI, r.RemoteAddr, wrappedWriter.statusCode, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
