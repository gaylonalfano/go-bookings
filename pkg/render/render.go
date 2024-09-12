package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gaylonalfano/go-bookings/pkg/config"
	"github.com/gaylonalfano/go-bookings/pkg/models"
)

// Let's pull in our AppConfig
var app *config.AppConfig

// NewTemplates sets the config for the render package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultTemplateData(td *models.TemplateData) *models.TemplateData {
	// TODO: Add custom default data to be available to all templates
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// Use our AppConfig template cache
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		// Couldn't get the Template from cache
		log.Fatal("could not retrieve passed template from template cache")
	}

	td = AddDefaultTemplateData(td)

	// NOTE: Optional check by first attempting to Execute
	// the template inside a temp Buffer, instead of in
	// the ResponseWriter. If error then it's from our map.
	buf := new(bytes.Buffer)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// Render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	// t := template.Must(template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.html"))
	// t.Execute(w, nil)
}

// NOTE:! When we have pages that use layouts, we have to parse them
// together into a single Template type.
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get all template files
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return cache, err
	}

	// Now have slice of string ("./templates/home.page.html")
	for i, page := range pages {
		log.Printf("iteration %d", i)
		// Get the actual file name
		name := filepath.Base(page)
		// Parse the file to create a Template
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		// Now get any Layouts we've created
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			// We have some layout paths. Time to parse and create template and add to our template set
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return cache, err
			}
		}

		// Finally, add new page and layout templates to cache
		cache[name] = ts
	}

	return cache, nil
}

// ----- Simple Cache Setup -----
// NOTE: We want to cache our templates so we don't have to
// read from disk on each request.
var tc = make(map[string]*template.Template)

func RenderTemplateSimple(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// Check to see if we already have Template in cache
	_, isCached := tc[t]
	if !isCached {
		// Need to create and render the template
		// Store template in our template cache
		log.Println("creating template and adding to cache")
		err = createTemplateCacheSimple(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// Template already in cache. Need to render
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCacheSimple(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.html",
	}

	// Parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// Add template to cache
	tc[t] = tmpl

	return nil
}
