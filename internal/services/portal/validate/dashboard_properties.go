// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
)

func DashboardProperties(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return warnings, errors
	}

	var dashboardProperties dashboard.DashboardProperties
	if err := json.Unmarshal([]byte(value), &dashboardProperties); err != nil {
		errors = append(errors, fmt.Errorf("parsing JSON: %+v", err))
		return warnings, errors
	}

	if dashboardProperties.Lenses == nil {
		errors = append(errors, fmt.Errorf("`lenses` is required in JSON payload"))
		return warnings, errors
	}

	return warnings, errors
}
