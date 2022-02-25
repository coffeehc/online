// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"online/server/web/gen/models"
)

// GetYakitPluginTagsHandlerFunc turns a function with the right signature into a get yakit plugin tags handler
type GetYakitPluginTagsHandlerFunc func(GetYakitPluginTagsParams, *models.Principle) middleware.Responder

// Handle executing the request and returning a response
func (fn GetYakitPluginTagsHandlerFunc) Handle(params GetYakitPluginTagsParams, principal *models.Principle) middleware.Responder {
	return fn(params, principal)
}

// GetYakitPluginTagsHandler interface for that can handle valid get yakit plugin tags params
type GetYakitPluginTagsHandler interface {
	Handle(GetYakitPluginTagsParams, *models.Principle) middleware.Responder
}

// NewGetYakitPluginTags creates a new http.Handler for the get yakit plugin tags operation
func NewGetYakitPluginTags(ctx *middleware.Context, handler GetYakitPluginTagsHandler) *GetYakitPluginTags {
	return &GetYakitPluginTags{Context: ctx, Handler: handler}
}

/* GetYakitPluginTags swagger:route GET /yakit/plugin/tags getYakitPluginTags

GetYakitPluginTags get yakit plugin tags API

*/
type GetYakitPluginTags struct {
	Context *middleware.Context
	Handler GetYakitPluginTagsHandler
}

func (o *GetYakitPluginTags) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetYakitPluginTagsParams()
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
