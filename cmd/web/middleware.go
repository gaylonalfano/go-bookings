package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

// WARN: REF: https://www.youtube.com/watch?v=irSUU1K1qtU
// Q: Not sure of the difference between http.Handler and
// func(w http.ResponseWriter, r *http.Request). For a simple func,
// I can use the MW via SecondMiddleware(handler.Repo.Home).
// However, I can't use FirstMiddleware(handler.Repo.Home)...
// NOTE: !! REF: https://youtu.be/H7tbjKFSg58?t=430
// Watch this! Dreams of Code explanation with 1.22
// U: You pass it as the http.Server{Handler: middleware.Logging(router)}!!!!
// NOTE: To use the MW you simply pass original route handler to MW
//

// NOTE: !! Middleware Chaining approach with native net/http
// REF: https://youtu.be/H7tbjKFSg58?t=562
// REF: https://www.reddit.com/r/golang/comments/1aoxlsr/middleware_in_go_1220/
type Middleware func(http.Handler) http.Handler

// Q: Variadic array??
func CreateMiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		// return the topmost middleware
		return next
	}
}

// FirstMiddleware uses http.Handler
func FirstMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("inside middleware 1 - before handler")
		next.ServeHTTP(w, r)
		log.Print("inside middleware 1 - after handler")
	})
}

// SecondMiddleware uses func(http.ResponseWriter, *http.Request)
func SecondMiddleware(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("inside middleware 2 - before handler")
		f(w, r) // original function call
		log.Print("inside middleware 2 - after handler")
	}
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// REF: https://youtu.be/H7tbjKFSg58?t=483
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Q: How to log the response status code?
		// A: Create a custom type that conforms to the ResponseWriter interface
		wrappedWriter := &wrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)

		log.Println(wrappedWriter.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit the page")
		next.ServeHTTP(w, r)
	})
}

// NoSurfCSRF adds CSRF protection to all POST requests
func NoSurfCSRF(next http.Handler) http.Handler {
	fmt.Println("MW: NoSurfCSRF")
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: app.InProduction,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	fmt.Println("MW: SessionLoad")
	return sessionManager.LoadAndSave(next)
}
