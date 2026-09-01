// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SnapshotName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\da-zA-Z][-_\da-zA-Z]{3,63}$`), "must be between 4 and 64 characters in length and start with letters or numbers and contains only letters, numbers, underscore or hyphens")(v, k)
}
