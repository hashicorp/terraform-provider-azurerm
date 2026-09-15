// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Swagger says: Supported characters for the name are a-z, A-Z, 0-9 and _. The maximum name length is 80 characters.
// Confirmed with the service team, they gave me a regex: ^[^_\W][\w-._]{0,79}(?<![-.])$
// This means the name can contain a-z, A-Z, 0-9, underscore, dot or hyphen, and must not starts with a underscore or any other non-word characters (underscore is considered as word character)
// additionally, the name cannot end with hyphen or dot.
// Golang regex does not support "negative look ahead" (aka `?<!`), therefore I transformed this regex to the following regular expression.
func DiskEncryptionSetName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^_\W][\w-._]{0,78}\w$`), "must be between 1 - 80 characters long, and contains only a-z, A-Z, 0-9 and _")(v, k)
}
