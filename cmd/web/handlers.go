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
	drives, err := app.drives.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Drives: drives,
	}

	app.render(w, http.StatusOK, "home.html", data)
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

	data := &templateData{
		Drive: drive,
	}

	app.render(w, http.StatusOK, "drive.html", data)
}

func (app *application) driveCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You're looking at the drive creation page")
}

func (app *application) driveCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "GDSC Web Dev Project Screening"
	company := "Google Developer Student Club - SJEC"
	description := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vitae proin sagittis nisl rhoncus. In mollis nunc sed id semper. Enim diam vulputate ut pharetra. Diam vel quam elementum pulvinar etiam non quam lacus suspendisse. Neque gravida in fermentum et sollicitudin ac orci phasellus egestas. Netus et malesuada fames ac turpis egestas sed tempus urna. Curabitur gravida arcu ac tortor. In hendrerit gravida rutrum quisque non tellus orci. Elementum sagittis vitae et leo duis ut. Nisl purus in mollis nunc sed id semper risus in. Feugiat scelerisque varius morbi enim nunc faucibus. Mattis vulputate enim nulla aliquet. Tincidunt arcu non sodales neque. Proin sagittis nisl rhoncus mattis rhoncus.`
	date := time.Now().Truncate(24*time.Hour).AddDate(0, 0, 365)

	id, err := app.drives.Insert(title, company, description, date)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drive/%d", id), http.StatusSeeOther)
}
