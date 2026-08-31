// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

func StorageContainerName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^\$root$|^\$web$|^[0-9a-z-]+$`), "only lowercase alphanumeric characters and hyphens allowed"),
		validation.StringLenBetween(3, 63),
		validation.StringDoesNotMatch(regexp.MustCompile(`^-`), "cannot begin with a hyphen"),
	)(v, k)
}

func StorageContainerDataPlaneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if client.StorageDomainSuffix == nil {
		return validation.IsURLWithPath(input, key)
	}

	if _, err := containers.ParseContainerID(v, *client.StorageDomainSuffix); err != nil {
		errors = append(errors, err)
	}

	return
}
