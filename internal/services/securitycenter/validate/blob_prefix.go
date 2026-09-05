// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BlobPrefix(v interface{}, k string) (warnings []string, errors []error) {
	values := strings.Split(v.(string), "/")
	containerName := values[0]

	if _, stringLenBetweenErrors := validation.StringLenBetween(1, 489)(v, k); stringLenBetweenErrors != nil {
		errors = append(errors, stringLenBetweenErrors...)
	}
	if _, stringDoesNotContainAnyErrors := validation.StringDoesNotContainAny("\\?#%&+:*\"|")(v, k); len(stringDoesNotContainAnyErrors) > 0 {
		errors = append(errors, stringDoesNotContainAnyErrors...)
	}
	if !regexp.MustCompile(`^[0-9a-z-]+$`).MatchString(containerName) {
		errors = append(errors, fmt.Errorf("only lowercase alphanumeric characters and hyphens allowed in %q container name: %q", k, containerName))
	}
	if len(containerName) > 63 {
		errors = append(errors, fmt.Errorf("%q container name must be at most 63 characters: %q", k, containerName))
	}
	if len(values) > 1 && len(containerName) < 3 {
		errors = append(errors, fmt.Errorf("%q container name must be at least 3 characters: %q", k, containerName))
	}
	if regexp.MustCompile(`^-`).MatchString(containerName) {
		errors = append(errors, fmt.Errorf("%q container name cannot begin with a hyphen: %q", k, containerName))
	}

	return warnings, errors
}
