package config

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// NOTE: !! The key is to ONLY import from stdlib. Don't want to
// import other pkgs from our app!

// AppConfig holds the application config
type AppConfig struct {
	Router         http.Handler
	TemplateCache  map[string]*template.Template
	InfoLog        *log.Logger
	UseCache       bool
	InProduction   bool
	SessionManager *scs.SessionManager
}
