// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// IsBase64String is a SchemaValidateFunc which tests if the provided value is a valid base64
// encoded string. Unlike IsBase64StringOrEmpty it also rejects empty and whitespace-only values,
// which base64 decodes happily but which callers never want.
func IsBase64String(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("%q must not be empty", k)}
	}

	if _, err := base64.StdEncoding.DecodeString(v); err != nil {
		return nil, []error{fmt.Errorf("%q must be a valid base64 encoded string", k)}
	}

	return nil, nil
}
