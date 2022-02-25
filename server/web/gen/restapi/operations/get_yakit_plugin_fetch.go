// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"online/server/web/gen/models"
)

// GetYakitPluginFetchHandlerFunc turns a function with the right signature into a get yakit plugin fetch handler
type GetYakitPluginFetchHandlerFunc func(GetYakitPluginFetchParams, *models.Principle) middleware.Responder

// Handle executing the request and returning a response
func (fn GetYakitPluginFetchHandlerFunc) Handle(params GetYakitPluginFetchParams, principal *models.Principle) middleware.Responder {
	return fn(params, principal)
}

// GetYakitPluginFetchHandler interface for that can handle valid get yakit plugin fetch params
type GetYakitPluginFetchHandler interface {
	Handle(GetYakitPluginFetchParams, *models.Principle) middleware.Responder
}

// NewGetYakitPluginFetch creates a new http.Handler for the get yakit plugin fetch operation
func NewGetYakitPluginFetch(ctx *middleware.Context, handler GetYakitPluginFetchHandler) *GetYakitPluginFetch {
	return &GetYakitPluginFetch{Context: ctx, Handler: handler}
}

/* GetYakitPluginFetch swagger:route GET /yakit/plugin/fetch getYakitPluginFetch

GetYakitPluginFetch get yakit plugin fetch API

*/
type GetYakitPluginFetch struct {
	Context *middleware.Context
	Handler GetYakitPluginFetchHandler
}

func (o *GetYakitPluginFetch) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetYakitPluginFetchParams()
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
