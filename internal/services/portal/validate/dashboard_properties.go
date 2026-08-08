// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/dashboards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func DashboardProperties(v interface{}, k string) (warnings []string, errs []error) {
	value := v.(string)

	if len(value) == 0 {
		errs = append(errs, fmt.Errorf("%q must not be empty", k))
		return warnings, errs
	}

	if !features.FivePointOh() {
		var dashboardProperties dashboard.DashboardProperties
		if err := json.Unmarshal([]byte(value), &dashboardProperties); err == nil {
			if dashboardProperties.Lenses == nil {
				errs = append(errs, errors.New("`lenses` is required in JSON payload"))
				return warnings, errs
			}
			for _, v := range *dashboardProperties.Lenses {
				if v.Parts == nil {
					errs = append(errs, errors.New("`lenses.parts` is required in JSON payload"))
					return warnings, errs
				}
			}
			return warnings, errs
		}
	}

	var dashboardProperties dashboards.DashboardPropertiesWithProvisioningState

	if err := json.Unmarshal([]byte(value), &dashboardProperties); err != nil {
		errs = append(errs, fmt.Errorf("parsing JSON: %+v", err))
		return warnings, errs
	}

	if dashboardProperties.Lenses == nil {
		errs = append(errs, errors.New("`lenses` is required in JSON payload"))
		return warnings, errs
	}

	for _, v := range *dashboardProperties.Lenses {
		if v.Parts == nil {
			errs = append(errs, errors.New("`lenses.parts` is required in JSON payload"))
			return warnings, errs
		}
	}
	return warnings, errs
}
