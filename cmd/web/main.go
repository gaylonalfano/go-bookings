package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gaylonalfano/go-bookings/pkg/config"
	"github.com/gaylonalfano/go-bookings/pkg/handler"
	"github.com/gaylonalfano/go-bookings/pkg/render"
)

// NOTE: General project structure template: https://github.com/golang-standards/project-layout

const portNumber = ":8080"

// The AppConfig will store template cache along with other stuff
// Check out the Redis example where App { db, router, config }
var app config.AppConfig

// Need sessionManager to be available to middleware pkg
// NOTE: To make sessionManager available to ALL pkgs, best
// to add to AppConfig.
var sessionManager *scs.SessionManager

func main() {
	app.InProduction = false

	// Session management using scs package
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction // https

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	repo := handler.NewRepo(&app)
	handler.NewHandlers(repo)
	render.NewTemplates(&app)

	r := loadRoutes(&app)
	// NOTE: !! Add custom MW by passing as Handler to Server!
	// Q: How to pass/use MULTIPLE middleware funcs?
	// REF: https://youtu.be/H7tbjKFSg58?t=462
	// REF: https://www.reddit.com/r/golang/comments/1f8dt5d/how_to_write_a_logging_middleware_for_nethttp/
	// A: Use Middleware Chaining approach! Otherwise you have to nest again and again!
	// REF: https://youtu.be/H7tbjKFSg58?t=539
	stack := CreateMiddlewareStack(
		Logging,
		// AllowCors,
		// IsAuthed,
		// CheckPermissions,
		NoSurfCSRF,
		SessionLoad,
		FirstMiddleware,
		WriteToConsole,
	)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        stack(r),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	app.UseCache = false
	app.TemplateCache = tc
	app.Router = r
	app.SessionManager = sessionManager

	fmt.Printf("App running on port %s \n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
