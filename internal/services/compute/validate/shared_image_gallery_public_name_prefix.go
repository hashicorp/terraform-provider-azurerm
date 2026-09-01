// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SharedImageGalleryPrefix(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9]{5,16}$"), "must be 5 to 16 characters long, and can only contain alphanumeric")(v, k)
}
