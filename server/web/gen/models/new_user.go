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

// NewUser new user
//
// swagger:model NewUser
type NewUser struct {

	// email
	Email string `json:"email,omitempty"`

	// from platform
	// Required: true
	FromPlatform *string `json:"from_platform"`

	// tags
	// Required: true
	Tags []string `json:"tags"`

	// trusted
	// Required: true
	Trusted *bool `json:"trusted"`

	// uesr unique id
	// Required: true
	UesrUniqueID *string `json:"uesr_unique_id"`

	// user verbose
	// Required: true
	UserVerbose *string `json:"user_verbose"`
}

// Validate validates this new user
func (m *NewUser) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFromPlatform(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTrusted(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUesrUniqueID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserVerbose(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NewUser) validateFromPlatform(formats strfmt.Registry) error {

	if err := validate.Required("from_platform", "body", m.FromPlatform); err != nil {
		return err
	}

	return nil
}

func (m *NewUser) validateTags(formats strfmt.Registry) error {

	if err := validate.Required("tags", "body", m.Tags); err != nil {
		return err
	}

	return nil
}

func (m *NewUser) validateTrusted(formats strfmt.Registry) error {

	if err := validate.Required("trusted", "body", m.Trusted); err != nil {
		return err
	}

	return nil
}

func (m *NewUser) validateUesrUniqueID(formats strfmt.Registry) error {

	if err := validate.Required("uesr_unique_id", "body", m.UesrUniqueID); err != nil {
		return err
	}

	return nil
}

func (m *NewUser) validateUserVerbose(formats strfmt.Registry) error {

	if err := validate.Required("user_verbose", "body", m.UserVerbose); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this new user based on context it is used
func (m *NewUser) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *NewUser) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NewUser) UnmarshalBinary(b []byte) error {
	var res NewUser
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
