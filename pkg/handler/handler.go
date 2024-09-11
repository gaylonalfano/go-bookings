package handler

import (
	"net/http"

	"github.com/gaylonalfano/go-study/basicapp/pkg/config"
	"github.com/gaylonalfano/go-study/basicapp/pkg/models"
	"github.com/gaylonalfano/go-study/basicapp/pkg/render"
)

// NOTE: Repository pattern in Go: https://www.udemy.com/course/building-modern-web-applications-with-go/learn/lecture/22875035#questions/14813628

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Store user's IP address into our session manager
	remoteIP := r.RemoteAddr
	m.App.SessionManager.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some business logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from About page data."

	// Extract user's stored session data (e.g., IP address)
	remoteIP := m.App.SessionManager.GetString(r.Context(), "remote_ip")
	// Put inside template data
	stringMap["remote_ip"] = remoteIP

	// send the processed data to template
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{StringMap: stringMap})
}
