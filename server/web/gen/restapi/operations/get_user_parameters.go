// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetUserParams creates a new GetUserParams object
// with the default values initialized.
func NewGetUserParams() GetUserParams {

	var (
		// initialize parameters with default values

		limitDefault = int64(20)

		orderDefault = string("desc")

		pageDefault = int64(1)
	)

	return GetUserParams{
		Limit: &limitDefault,

		Order: &orderDefault,

		Page: &pageDefault,
	}
}

// GetUserParams contains all the bound params for the get user operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetUser
type GetUserParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	  Default: 20
	*/
	Limit *int64
	/*
	  In: query
	*/
	Name *string
	/*
	  In: query
	  Default: "desc"
	*/
	Order *string
	/*
	  In: query
	*/
	OrderBy *string
	/*
	  In: query
	  Default: 1
	*/
	Page *int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetUserParams() beforehand.
func (o *GetUserParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qName, qhkName, _ := qs.GetOK("name")
	if err := o.bindName(qName, qhkName, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrder, qhkOrder, _ := qs.GetOK("order")
	if err := o.bindOrder(qOrder, qhkOrder, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrderBy, qhkOrderBy, _ := qs.GetOK("order_by")
	if err := o.bindOrderBy(qOrderBy, qhkOrderBy, route.Formats); err != nil {
		res = append(res, err)
	}

	qPage, qhkPage, _ := qs.GetOK("page")
	if err := o.bindPage(qPage, qhkPage, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *GetUserParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetUserParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int64", raw)
	}
	o.Limit = &value

	return nil
}

// bindName binds and validates parameter Name from query.
func (o *GetUserParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Name = &raw

	return nil
}

// bindOrder binds and validates parameter Order from query.
func (o *GetUserParams) bindOrder(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetUserParams()
		return nil
	}
	o.Order = &raw

	if err := o.validateOrder(formats); err != nil {
		return err
	}

	return nil
}

// validateOrder carries on validations for parameter Order
func (o *GetUserParams) validateOrder(formats strfmt.Registry) error {

	if err := validate.EnumCase("order", "query", *o.Order, []interface{}{"asc", "desc"}, true); err != nil {
		return err
	}

	return nil
}

// bindOrderBy binds and validates parameter OrderBy from query.
func (o *GetUserParams) bindOrderBy(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.OrderBy = &raw

	return nil
}

// bindPage binds and validates parameter Page from query.
func (o *GetUserParams) bindPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetUserParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("page", "query", "int64", raw)
	}
	o.Page = &value

	return nil
}
