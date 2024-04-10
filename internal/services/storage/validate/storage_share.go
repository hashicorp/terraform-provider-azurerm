// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/file/shares"
)

func StorageShareDataPlaneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if client.StorageDomainSuffix == nil {
		return validation.IsURLWithPath(input, key)
	}

	if _, err := shares.ParseShareID(v, *client.StorageDomainSuffix); err != nil {
		errors = append(errors, err)
	}

	return
}

func StorageShareName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q: %q",
			k, value))
	}

	// Following the naming convention as laid out in the docs https://msdn.microsoft.com/library/azure/dn167011.aspx
	if len(value) < 3 || len(value) > 63 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 3 and 63 characters: %q", k, value))
	}
	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot begin with a hyphen: %q", k, value))
	}
	if regexp.MustCompile(`[-]{2,}`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q does not allow consecutive hyphens: %q", k, value))
	}
	return warnings, errors
}
