package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"sjecplacement.in/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.html")
}

func (app *application) driveView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	drive, err := app.drives.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", drive)
}

func (app *application) driveCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You're looking at the drive creation page")
}

func (app *application) driveCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "The sceptix club recruitement"
	company := "The sceptix club"
	description := `LOREM IPSUM DOLOR AMET`
	date := time.Now().Truncate(24 * time.Hour)

	id, err := app.drives.Insert(title, company, description, date)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drive/%d", id), http.StatusSeeOther)
}
