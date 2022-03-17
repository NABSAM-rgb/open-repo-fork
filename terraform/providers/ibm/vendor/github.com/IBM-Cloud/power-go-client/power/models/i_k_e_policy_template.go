// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// IKEPolicyTemplate i k e policy template
// swagger:model IKEPolicyTemplate
type IKEPolicyTemplate struct {

	// ikePolicy Authentication default value
	// Required: true
	Authentication *string `json:"authentication"`

	// ikePolicy DHGroup default value
	// Required: true
	DhGroup *int64 `json:"dhGroup"`

	// ikePolicy Encryption default value
	// Required: true
	Encryption *string `json:"encryption"`

	// key lifetime
	// Required: true
	KeyLifetime KeyLifetime `json:"keyLifetime"`

	// ikePolicy Version default value
	// Required: true
	Version *float64 `json:"version"`
}

// Validate validates this i k e policy template
func (m *IKEPolicyTemplate) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAuthentication(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDhGroup(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEncryption(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKeyLifetime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersion(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IKEPolicyTemplate) validateAuthentication(formats strfmt.Registry) error {

	if err := validate.Required("authentication", "body", m.Authentication); err != nil {
		return err
	}

	return nil
}

func (m *IKEPolicyTemplate) validateDhGroup(formats strfmt.Registry) error {

	if err := validate.Required("dhGroup", "body", m.DhGroup); err != nil {
		return err
	}

	return nil
}

func (m *IKEPolicyTemplate) validateEncryption(formats strfmt.Registry) error {

	if err := validate.Required("encryption", "body", m.Encryption); err != nil {
		return err
	}

	return nil
}

func (m *IKEPolicyTemplate) validateKeyLifetime(formats strfmt.Registry) error {

	if err := m.KeyLifetime.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("keyLifetime")
		}
		return err
	}

	return nil
}

func (m *IKEPolicyTemplate) validateVersion(formats strfmt.Registry) error {

	if err := validate.Required("version", "body", m.Version); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IKEPolicyTemplate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IKEPolicyTemplate) UnmarshalBinary(b []byte) error {
	var res IKEPolicyTemplate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}