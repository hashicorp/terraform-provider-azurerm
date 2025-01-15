// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func StorageTableDataPlaneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if client.StorageDomainSuffix == nil {
		return validation.IsURLWithPath(input, key)
	}

	if _, err := tables.ParseTableID(v, *client.StorageDomainSuffix); err != nil {
		errors = append(errors, err)
	}

	return
}

func StorageTableName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if value == "table" {
		errors = append(errors, fmt.Errorf(
			"Table Storage %q cannot use the word `table`: %q",
			k, value))
	}
	if !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]{2,62}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"Table Storage %q cannot begin with a numeric character, only alphanumeric characters are allowed and must be between 3 and 63 characters long: %q",
			k, value))
	}

	return warnings, errors
}
