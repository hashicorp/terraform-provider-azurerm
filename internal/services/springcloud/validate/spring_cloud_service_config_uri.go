// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// the config server URI should be started with http://, https://, git@, or ssh://
func ConfigServerURI(i interface{}, k string) ([]string, []error) {
	return validation.StringStartsWithOneOf("http://", "https://", "git@", "ssh://")(i, k)
}
