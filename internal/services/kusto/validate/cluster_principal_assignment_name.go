// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ClusterPrincipalAssignmentName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`^[\s]+$`), "must not consist of whitespaces only"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`), "may only contain alphanumeric characters, whitespaces, dashes and dots"),
		validation.StringLenBetween(0, 260),
	)(v, k)
}
