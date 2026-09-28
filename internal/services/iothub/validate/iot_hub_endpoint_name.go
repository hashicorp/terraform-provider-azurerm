// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"

func IoTHubEndpointName(v interface{}, k string) ([]string, []error) {
	reservedNames := []string{
		"events",
		"operationsMonitoringEvents",
		"fileNotifications",
		"$default",
	}

	return validation.StringNotInSlice(reservedNames, false)(v, k)
}
