package main

import (
	"net/http"

	"snippetbox.sarasproject.net/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.SessionManager.LoadAndSave,noSurf, app.authenticate)

	router.Handle(http.MethodGet, "/", toHandle(dynamic.ThenFunc(app.home)))
	router.Handle(http.MethodGet, "/snippet/view/:id", toHandle(dynamic.ThenFunc(app.snippetView)))
	router.Handle(http.MethodGet, "/snippet/create", toHandle(dynamic.ThenFunc(app.snippetCreate)))
	router.Handle(http.MethodPost, "/snippet/create", toHandle(dynamic.ThenFunc(app.snippetCreatePost)))

	router.Handle(http.MethodGet, "/user/signup", toHandle(dynamic.ThenFunc(app.userSignup)))
	router.Handle(http.MethodPost, "/user/signup", toHandle(dynamic.ThenFunc(app.userSignupPost)))

	router.Handle(http.MethodGet, "/user/login", toHandle(dynamic.ThenFunc(app.userLogin)))
	router.Handle(http.MethodPost, "/user/login", toHandle(dynamic.ThenFunc(app.userLoginPost)))
	router.Handle(http.MethodPost, "/user/logout", toHandle(dynamic.ThenFunc(app.userLogoutPost)))

	// protected := dynamic.Append(app.requireAuthentication)
	// router.Handler(http.MethodGet, "/snippet/create",
	// protected.ThenFunc(app.snippetCreate))
	// router.Handler(http.MethodPost, "/snippet/create",
	// protected.ThenFunc(app.snippetCreatePost))
	// router.Handler(http.MethodPost, "/user/logout",
	// protected.ThenFunc(app.userLogoutPost))


	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}

func toHandle(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
