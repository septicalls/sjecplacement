package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Unprotected routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/drive/:id", dynamic.ThenFunc(app.driveView))
	router.Handler(http.MethodPost, "/drive/:id", dynamic.ThenFunc(app.roleAddPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected routes
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/create", protected.ThenFunc(app.driveCreate))
	router.Handler(http.MethodPost, "/create", protected.ThenFunc(app.driveCreatePost))
	router.Handler(http.MethodPost, "/publish/:id", protected.ThenFunc(app.publishDrivePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
