package main

import (
	"path/filepath"
	"text/template"
	"time"

	"sjecplacement.in/internal/models"
)

type templateData struct {
	Drive  *models.Drive
	Drives []*models.Drive
	Roles  []*models.Role
	Form   any
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

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
