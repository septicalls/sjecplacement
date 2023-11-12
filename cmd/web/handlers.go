package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"sjecplacement.in/internal/models"
	"sjecplacement.in/internal/validator"
)

type driveCreateForm struct {
	Title               string `form:"title"`
	Company             string `form:"company"`
	Description         string `form:"description"`
	Date                string `form:"date"`
	validator.Validator `form:"-"`
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
	data := templateData{}

	data.Form = driveCreateForm{
		Date: time.Now().Format("2006-01-02"),
	}

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) driveCreatePost(w http.ResponseWriter, r *http.Request) {
	var form driveCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Title is mandatory")
	form.CheckField(validator.NotBlank(form.Company), "company", "Company name is mandatory")
	form.CheckField(validator.NotBlank(form.Description), "description", "Description is mandatory")
	form.CheckField(validator.MaxChar(form.Title, 100), "title", "Title can only be 100 characters long")
	form.CheckField(validator.MaxChar(form.Company, 100), "company", "Company can only be 100 characters long")

	date, dateOk := validator.ValidDate(form.Date)
	form.CheckField(dateOk, "date", "Must be a valid date beyond today")

	if !form.Valid() {
		data := templateData{
			Form: form,
		}
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.drives.Insert(form.Title, form.Company, form.Description, date)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drive/%d", id), http.StatusSeeOther)
}
