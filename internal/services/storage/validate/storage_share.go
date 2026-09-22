// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
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

// StorageShareName follows the naming convention as laid out in the docs https://msdn.microsoft.com/library/azure/dn167011.aspx
func StorageShareName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[0-9a-z-]+$`), "only lowercase alphanumeric characters and hyphens allowed"),
		validation.StringLenBetween(3, 63),
		validation.StringDoesNotMatch(regexp.MustCompile(`^-`), "cannot begin with a hyphen"),
		validation.StringDoesNotMatch(regexp.MustCompile(`[-]{2,}`), "does not allow consecutive hyphens"),
	)(v, k)
}
