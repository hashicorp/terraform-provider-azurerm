// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SqlImageOfferName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^SQL[A-Za-z0-9]*-WS[A-Za-z0-9]*$`), "should be in the form SQL<SQLversion>-WS<OSversion>, for example SQL2019-WS2019")(v, k)
}
