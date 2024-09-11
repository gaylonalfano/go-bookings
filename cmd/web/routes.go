package main

import (
	"net/http"

	"github.com/gaylonalfano/go-study/basicapp/pkg/config"
	"github.com/gaylonalfano/go-study/basicapp/pkg/handler"
)

func loadRoutes(app *config.AppConfig) http.Handler {
	mux := http.NewServeMux()

	// NOTE: !! Use the middleware by wrapping it around original route handler
	// NOTE: !! REF: https://youtu.be/H7tbjKFSg58?t=430
	// Watch this! Dreams of Code explanation with 1.22
	// U: You pass it as the http.Server{Handler: middleware.Logging(router)}!!!!
	// mux.HandleFunc("GET /", FirstMiddleware(handler.Repo.Home))  // http.Handler ERROR
	// mux.HandleFunc("GET /", SecondMiddleware(handler.Repo.Home)) // func(w,r) WORKS!
	mux.HandleFunc("GET /", handler.Repo.Home)
	mux.HandleFunc("GET /about", handler.Repo.About)

	return mux
}
