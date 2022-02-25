// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Paging paging
//
// swagger:model Paging
type Paging struct {

	// pagemeta
	// Required: true
	Pagemeta *PageMeta `json:"pagemeta"`
}

// Validate validates this paging
func (m *Paging) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePagemeta(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Paging) validatePagemeta(formats strfmt.Registry) error {

	if err := validate.Required("pagemeta", "body", m.Pagemeta); err != nil {
		return err
	}

	if m.Pagemeta != nil {
		if err := m.Pagemeta.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pagemeta")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this paging based on the context it is used
func (m *Paging) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePagemeta(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Paging) contextValidatePagemeta(ctx context.Context, formats strfmt.Registry) error {

	if m.Pagemeta != nil {
		if err := m.Pagemeta.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pagemeta")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Paging) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Paging) UnmarshalBinary(b []byte) error {
	var res Paging
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
