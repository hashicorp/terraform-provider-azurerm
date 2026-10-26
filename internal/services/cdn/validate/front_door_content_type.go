// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// FrontDoorContentType validates that the value is a well formed MIME type (e.g. `type/subtype`) as defined in RFC 6838,
// Front Door accepts any content type for compression so this intentionally does not restrict the value to a fixed list.
func FrontDoorContentType(v any, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z\d][a-zA-Z\d!#$&^_.+-]{0,126}/[a-zA-Z\d][a-zA-Z\d!#$&^_.+-]{0,126}$`), "must be a valid content type in the format `type/subtype`, e.g. `text/html`")(v, k)
}
