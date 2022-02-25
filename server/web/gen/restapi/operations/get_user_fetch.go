// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"online/server/web/gen/models"
)

// GetUserFetchHandlerFunc turns a function with the right signature into a get user fetch handler
type GetUserFetchHandlerFunc func(GetUserFetchParams, *models.Principle) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUserFetchHandlerFunc) Handle(params GetUserFetchParams, principal *models.Principle) middleware.Responder {
	return fn(params, principal)
}

// GetUserFetchHandler interface for that can handle valid get user fetch params
type GetUserFetchHandler interface {
	Handle(GetUserFetchParams, *models.Principle) middleware.Responder
}

// NewGetUserFetch creates a new http.Handler for the get user fetch operation
func NewGetUserFetch(ctx *middleware.Context, handler GetUserFetchHandler) *GetUserFetch {
	return &GetUserFetch{Context: ctx, Handler: handler}
}

/* GetUserFetch swagger:route GET /user/fetch getUserFetch

GetUserFetch get user fetch API

*/
type GetUserFetch struct {
	Context *middleware.Context
	Handler GetUserFetchHandler
}

func (o *GetUserFetch) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetUserFetchParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principle
	if uprinc != nil {
		principal = uprinc.(*models.Principle) // this is really a models.Principle, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}