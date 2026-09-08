// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DataConnectionName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`^[\s]+$`), "must not consist of whitespaces only"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`), "may only contain letters, digits, whitespaces, dashes and dots"),
		validation.StringLenBetween(0, 40),
	)(v, k)
}

func EntityName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`^[\s]+$`), "must not consist of whitespaces only"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_\s.-]+$`), "may only contain letters, digits, underscores, spaces, dashes and dots"),
		validation.StringLenBetween(0, 1024),
	)(v, k)
}

func ClusterName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\-]+$`), "must begin with a letter and may only contain alphanumeric characters"),
		validation.StringLenBetween(4, 22),
	)(v, k)
}

func DatabaseName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`^[\s]+$`), "must not consist of whitespaces only"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\s._-]+$`), "may only contain alphanumeric characters, whitespaces, dashes, underscores and dots"),
		validation.StringLenBetween(0, 260),
	)(v, k)
}

func DatabasePrincipalAssignmentName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`^[\s]+$`), "must not consist of whitespaces only"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`), "may only contain alphanumeric characters, whitespaces, dashes and dots"),
		validation.StringLenBetween(0, 260),
	)(v, k)
}
