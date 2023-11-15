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

type roleCreateForm struct {
	Profile             string `form:"profile"`
	Description         string `form:"description"`
	Qualification       string `form:"qualification"`
	Cutoff              string `form:"cutoff"`
	Location            string `form:"location"`
	Stipend             string `form:"stipend"`
	CTC                 string `form:"ctc"`
	ServiceAgreement    string `form:"serviceagreement"`
	DriveID             int    `form:"-"`
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

	roles, err := app.roles.All(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := &templateData{
		Drive: drive,
		Roles: roles,
		Form:  roleCreateForm{DriveID: id},
		Flash: flash,
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

	form.CheckField(validator.NotBlank(form.Title), "title", "This field is mandatory")
	form.CheckField(validator.NotBlank(form.Company), "company", "This field is mandatory")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field is mandatory")
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

	app.sessionManager.Put(r.Context(), "flash", "Drive Created Successfully")

	http.Redirect(w, r, fmt.Sprintf("/drive/%d", id), http.StatusSeeOther)
}

func (app *application) roleAddPost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := roleCreateForm{DriveID: id}

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Profile), "profile", "This field is mandatory")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field is mandatory")
	form.CheckField(validator.NotBlank(form.Qualification), "qualification", "This field is mandatory")

	cutoff, cutoffOk := validator.ValidString(form.Cutoff, 100)
	form.CheckField(cutoffOk, "cutoff", "Cut-off can only be 100 characters long")

	location, locationOk := validator.ValidString(form.Location, 100)
	form.CheckField(locationOk, "location", "Location can only be 100 characters long")

	stipend, stipendOk := validator.ValidInt(form.Stipend)
	form.CheckField(stipendOk, "stipend", "Must be a valid integer")

	ctc, ctcOk := validator.ValidFloat(form.CTC)
	form.CheckField(ctcOk, "ctc", "Must be a valid floating point")

	srvAgr, srvAgrOk := validator.ValidFloat(form.ServiceAgreement)
	form.CheckField(srvAgrOk, "serviceagreement", "Must be a valid floating point")

	if !form.Valid() {
		data := templateData{
			DriveID: id,
			Form:    form,
		}

		app.templateCache["role.htmx"].Execute(w, data)
		return
	}

		data := &templateData{
			Form:  form,
			Drive: drive,
			Roles: roles,
		}

		app.render(w, http.StatusUnprocessableEntity, "drive.html", data)
		return
	}

	role := models.Role{
		Profile:          form.Profile,
		Description:      form.Description,
		Qualification:    form.Qualification,
		Cutoff:           cutoff,
		Location:         location,
		Stipend:          stipend,
		CTC:              ctc,
		ServiceAgreement: srvAgr,
		DriveID:          id,
	}

	_, err = app.roles.Insert(&role)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drive/%d", id), http.StatusSeeOther)
}
