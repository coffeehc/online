// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"online/server/web/gen/models"
)

// DeleteYakitPluginHandlerFunc turns a function with the right signature into a delete yakit plugin handler
type DeleteYakitPluginHandlerFunc func(DeleteYakitPluginParams, *models.Principle) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteYakitPluginHandlerFunc) Handle(params DeleteYakitPluginParams, principal *models.Principle) middleware.Responder {
	return fn(params, principal)
}

// DeleteYakitPluginHandler interface for that can handle valid delete yakit plugin params
type DeleteYakitPluginHandler interface {
	Handle(DeleteYakitPluginParams, *models.Principle) middleware.Responder
}

// NewDeleteYakitPlugin creates a new http.Handler for the delete yakit plugin operation
func NewDeleteYakitPlugin(ctx *middleware.Context, handler DeleteYakitPluginHandler) *DeleteYakitPlugin {
	return &DeleteYakitPlugin{Context: ctx, Handler: handler}
}

/* DeleteYakitPlugin swagger:route DELETE /yakit/plugin deleteYakitPlugin

DeleteYakitPlugin delete yakit plugin API

*/
type DeleteYakitPlugin struct {
	Context *middleware.Context
	Handler DeleteYakitPluginHandler
}

func (o *DeleteYakitPlugin) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteYakitPluginParams()
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
