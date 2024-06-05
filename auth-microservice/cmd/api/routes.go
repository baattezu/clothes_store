package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthCheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/createModule", app.createClotheHandler)
	router.HandlerFunc(http.MethodGet, "/v1/getModule/:id", app.getClotheHandler)
	router.HandlerFunc(http.MethodGet, "/v1/getModule", app.listClotheHandler)
	router.HandlerFunc(http.MethodPut, "/v1/updateModule/:id", app.editClotheHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/deleteModule/:id", app.deleteClotheHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.requirePermission("read", app.getUserInfoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users", app.requirePermission("read", app.getAllUserInfoHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.requirePermission("write", app.editUserInfoHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.requirePermission("write", app.deleteUserInfoHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
