package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthCheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/createClothe", app.createClotheHandler)
	router.HandlerFunc(http.MethodGet, "/getClothe/:id", app.getClotheHandler)
	router.HandlerFunc(http.MethodGet, "/getClothe", app.listClotheHandler)
	router.HandlerFunc(http.MethodPut, "/updateClothe/:id", app.editClotheHandler)
	router.HandlerFunc(http.MethodDelete, "/deleteClothe/:id", app.deleteClotheHandler)

	return router
}
