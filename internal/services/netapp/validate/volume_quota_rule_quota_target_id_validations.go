// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func ValidateUnixUserIDOrGroupID(v interface{}, k string) (warnings []string, errors []error) {
	var value int64
	var err error

	switch v.(type) {
	case int:
		value = v.(int64)
	case string:
		if _, err := strconv.ParseInt(v.(string), 10, 64); err != nil {
			errors = append(errors, fmt.Errorf("%q must be an integer or a string that can be converted to an integer", k))
			return warnings, errors
		}
		if _, err := strconv.ParseInt(v.(string), 10, 64); err == nil && !regexp.MustCompile(`^\d+$`).MatchString(v.(string)) {
			errors = append(errors, fmt.Errorf("%q must be an integer or a string that contains only digits", k))
			return warnings, errors
		}
		value, err = strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q must be an integer or a string that can be converted to an integer", k))
			return warnings, errors
		}
	default:
		errors = append(errors, fmt.Errorf("%q must be an integer or a string that can be converted to an integer", k))
		return warnings, errors
	}

	if value < 1 || value > models.MaxQuotaTargetIDSizeInKiB {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 4294967295", k))
		return warnings, errors
	}

	return warnings, errors
}

func ValidateWindowsSID(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(*string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q must be a string", k))
		return warnings, errors
	}

	if !regexp.MustCompile(`^S-1-5-(0|18|\d{1,9})(-\d{1,10}){0,14}$`).MatchString(pointer.From(value)) {
		errors = append(errors, fmt.Errorf("%q must be a valid Windows security identifier (SID)", k))
		return warnings, errors
	}

	return warnings, errors
}
