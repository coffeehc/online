// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// YakitPluginsResponse yakit plugins response
//
// swagger:model YakitPluginsResponse
type YakitPluginsResponse struct {
	Paging

	// data
	// Required: true
	Data []*YakitPlugin `json:"data"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *YakitPluginsResponse) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 Paging
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.Paging = aO0

	// now for regular properties
	var propsYakitPluginsResponse struct {
		Data []*YakitPlugin `json:"data"`
	}
	if err := swag.ReadJSON(raw, &propsYakitPluginsResponse); err != nil {
		return err
	}
	m.Data = propsYakitPluginsResponse.Data

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m YakitPluginsResponse) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 1)

	aO0, err := swag.WriteJSON(m.Paging)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	// now for regular properties
	var propsYakitPluginsResponse struct {
		Data []*YakitPlugin `json:"data"`
	}
	propsYakitPluginsResponse.Data = m.Data

	jsonDataPropsYakitPluginsResponse, errYakitPluginsResponse := swag.WriteJSON(propsYakitPluginsResponse)
	if errYakitPluginsResponse != nil {
		return nil, errYakitPluginsResponse
	}
	_parts = append(_parts, jsonDataPropsYakitPluginsResponse)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this yakit plugins response
func (m *YakitPluginsResponse) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with Paging
	if err := m.Paging.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *YakitPluginsResponse) validateData(formats strfmt.Registry) error {

	if err := validate.Required("data", "body", m.Data); err != nil {
		return err
	}

	for i := 0; i < len(m.Data); i++ {
		if swag.IsZero(m.Data[i]) { // not required
			continue
		}

		if m.Data[i] != nil {
			if err := m.Data[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("data" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this yakit plugins response based on the context it is used
func (m *YakitPluginsResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with Paging
	if err := m.Paging.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *YakitPluginsResponse) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Data); i++ {

		if m.Data[i] != nil {
			if err := m.Data[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("data" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *YakitPluginsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *YakitPluginsResponse) UnmarshalBinary(b []byte) error {
	var res YakitPluginsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
