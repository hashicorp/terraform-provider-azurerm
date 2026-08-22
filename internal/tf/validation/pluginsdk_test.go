// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestIsURLWithHTTPorHTTPSOrEmpty(t *testing.T) {
	cases := map[string]struct {
		Value                  interface{}
		ExpectValidationErrors bool
	}{
		"accept empty string": {
			Value:                  "",
			ExpectValidationErrors: false,
		},
		"accept valid https URL": {
			Value:                  "https://index.docker.io",
			ExpectValidationErrors: false,
		},
		"accept valid http URL": {
			Value:                  "http://myregistry.example.com",
			ExpectValidationErrors: false,
		},
		"accept hostname without protocol": {
			Value:                  "index.docker.io",
			ExpectValidationErrors: false,
		},
		"accept hostname with path without protocol": {
			Value:                  "index.docker.io/v1",
			ExpectValidationErrors: false,
		},
		"accept private registry hostname": {
			Value:                  "myregistry.azurecr.io",
			ExpectValidationErrors: false,
		},
		"accept https URL with path": {
			Value:                  "https://index.docker.io/v1",
			ExpectValidationErrors: false,
		},
		"reject non-string value": {
			Value:                  123,
			ExpectValidationErrors: true,
		},
	}

	for tn, tc := range cases {
		_, errors := IsURLWithHTTPorHTTPSOrEmpty(tc.Value, tn)
		if len(errors) > 0 && !tc.ExpectValidationErrors {
			t.Errorf("%s: unexpected errors %s", tn, errors)
		} else if len(errors) == 0 && tc.ExpectValidationErrors {
			t.Errorf("%s: expected errors but got none", tn)
		}
	}
}

func TestValidateFloatInSlice(t *testing.T) {
	cases := map[string]struct {
		Value                  interface{}
		ValidateFunc           pluginsdk.SchemaValidateFunc
		ExpectValidationErrors bool
	}{
		"accept valid value": {
			Value:                  1.5,
			ValidateFunc:           FloatInSlice([]float64{1.0, 1.5, 2.0}),
			ExpectValidationErrors: false,
		},
		"accept valid negative value ": {
			Value:                  -1.0,
			ValidateFunc:           FloatInSlice([]float64{-1.0, 2.0}),
			ExpectValidationErrors: false,
		},
		"accept zero": {
			Value:                  0.0,
			ValidateFunc:           FloatInSlice([]float64{0.0, 2.0}),
			ExpectValidationErrors: false,
		},
		"reject out of range value": {
			Value:                  -1.0,
			ValidateFunc:           FloatInSlice([]float64{0.0, 2.0}),
			ExpectValidationErrors: true,
		},
		"reject incorrectly typed value": {
			Value:                  1,
			ValidateFunc:           FloatInSlice([]float64{0, 1, 2}),
			ExpectValidationErrors: true,
		},
	}

	for tn, tc := range cases {
		_, errors := tc.ValidateFunc(tc.Value, tn)
		if len(errors) > 0 && !tc.ExpectValidationErrors {
			t.Errorf("%s: unexpected errors %s", tn, errors)
		} else if len(errors) == 0 && tc.ExpectValidationErrors {
			t.Errorf("%s: expected errors but got none", tn)
		}
	}
}
