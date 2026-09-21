// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SubscriptionName(i interface{}, k string) (warnings []string, errs []error) {
	return validation.All(
		validation.StringLenBetween(1, 64),
		validation.StringDoesNotContainAny("<>;|"),
	)(i, k)
}
