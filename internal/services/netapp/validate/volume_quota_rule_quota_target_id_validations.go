// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// lintignore:V001 // numeric ID parsing and range checks, not a plain regex match
func ValidateUnixUserIDOrGroupID(v interface{}, k string) (warnings []string, errors []error) {
	var value int64
	var err error

	switch v := v.(type) {
	case int:
		value = int64(v)
	case string:
		if _, err := strconv.ParseInt(v, 10, 64); err != nil {
			errors = append(errors, fmt.Errorf("%q must be an integer or a string that can be converted to an integer", k))
			return warnings, errors
		}
		if _, err := strconv.ParseInt(v, 10, 64); err == nil && !regexp.MustCompile(`^\d+$`).MatchString(v) {
			errors = append(errors, fmt.Errorf("%q must be an integer or a string that contains only digits", k))
			return warnings, errors
		}
		value, err = strconv.ParseInt(v, 10, 64)
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

func ValidateWindowsSID(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^S-1-5-(0|18|\d{1,9})(-\d{1,10}){0,14}$`), "must be a valid Windows security identifier (SID)")(v, k)
}
