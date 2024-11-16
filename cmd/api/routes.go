package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {

	// setup a new router
	router := httprouter.New()
	// handle 404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	// handle 405
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
	// setup routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.updateCommentHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.deleteCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments", a.listCommentsHandler)

	// user routers:
	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)

	// We use PUT instead of POST because PUT is idempotent
	// and appropriate for this endpoint.  The activation
	// should only happens once, also we are not creating a resource
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activateUserHandler)

	// Request sent first to recoverPanic() then sent to rateLimit()
	// finally it is sent to the router.
	return a.recoverPanic(a.rateLimit(router))

}
