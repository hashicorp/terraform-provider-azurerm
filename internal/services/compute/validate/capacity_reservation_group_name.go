// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CapacityReservationGroupName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^_\W]([\w-._]{0,62}[\w_])?$`), `The Capacity Reservation Group Name must be between 1 and 64 characters long. It cannot contain special characters \/"[]:|<>+=;,?*@&, whitespace, or begin with '_' or end with '.' or '-'`)
}
