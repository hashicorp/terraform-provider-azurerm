// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FrontDoorCustomDomainHostName(i interface{}, k string) (_ []string, errors []error) {
	// an FQDN of at least two labels, where the first label may be the wildcard `*` and every
	// label is alphanumeric plus non-leading/non-trailing hyphens, at most 63 characters each
	label := `[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?`
	fqdn := fmt.Sprintf(`^(?:\*|%[1]s)(?:\.%[1]s)+$`, label)

	return validation.All(
		validation.StringLenBetween(1, 253),
		validation.StringMatch(regexp.MustCompile(fqdn), "must be a valid fully qualified domain name"),
	)(i, k)
}
