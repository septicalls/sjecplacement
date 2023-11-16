package main

import (
	"fmt"
	"path/filepath"
	"text/template"
	"time"

	"sjecplacement.in/internal/models"
)

type templateData struct {
	Drive   *models.Drive
	Role    *models.Role
	Drives  []*models.Drive
	Roles   []*models.Role
	DriveID int
	Flash   string
	Form    any
}

func humanDate(t time.Time) string {
	return t.Format("Monday, January 02, 2006")
}

func isDone(t time.Time) bool {
	return !t.After(time.Now())
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"isDone":    isDone,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Some of the HTMX files are also used as partials in the main templates.
		// This requires the htmx to be {{defined "like"}} this.
		ts, err = ts.ParseGlob("./ui/htmx/*.htmx")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	// Load HTMX templates
	htmxFiles, err := filepath.Glob("./ui/htmx/*.htmx")
	if err != nil {
		return nil, err
	}

	for _, htmxFile := range htmxFiles {
		name := filepath.Base(htmxFile)

		// This is to render the raw html from the .htmx files while supporting the
		// partial {{definition "for"}} templates.
		ts, err := template.New(name).Funcs(functions).Parse(fmt.Sprintf(`{{template "%s" .}}`, name))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(htmxFile)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
