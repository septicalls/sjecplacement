package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.html")
}

func (app *application) driveView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFount(w)
		return
	}

	fmt.Fprintf(w, "You're looking at the drive %d.", id)
}

func (app *application) driveCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You're looking at the drive creation page")
}
