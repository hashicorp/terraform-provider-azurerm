// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageSyncName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[0-9a-zA-Z-_. ]*[0-9a-zA-Z-_]$"), "can only consist of letters, numbers, spaces, and any of the following characters: '.-_' and that does not end with characters: '. '")(v, k)
}
