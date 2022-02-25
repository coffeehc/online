// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Operation operation
//
// swagger:model Operation
type Operation struct {
	GormBaseModel

	NewOperation
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *Operation) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 GormBaseModel
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.GormBaseModel = aO0

	// AO1
	var aO1 NewOperation
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.NewOperation = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m Operation) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.GormBaseModel)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.NewOperation)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this operation
func (m *Operation) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with GormBaseModel
	if err := m.GormBaseModel.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with NewOperation
	if err := m.NewOperation.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this operation based on the context it is used
func (m *Operation) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with GormBaseModel
	if err := m.GormBaseModel.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with NewOperation
	if err := m.NewOperation.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Operation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Operation) UnmarshalBinary(b []byte) error {
	var res Operation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
