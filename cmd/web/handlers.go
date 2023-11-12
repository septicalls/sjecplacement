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

type driveCreateForm struct {
	Title       string `form:"title"`
	Company     string `form:"company"`
	Description string `form:"description"`
	Date        string `form:"date"`
}

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
	app.render(w, http.StatusOK, "create.html", nil)
}

func (app *application) driveCreatePost(w http.ResponseWriter, r *http.Request) {
	var form driveCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", form.Date)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.drives.Insert(form.Title, form.Company, form.Description, date)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drive/%d", id), http.StatusSeeOther)
}
