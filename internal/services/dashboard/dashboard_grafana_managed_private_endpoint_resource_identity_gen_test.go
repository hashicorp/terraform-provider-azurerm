// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dashboard_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDashboardGrafanaManagedPrivateEndpoint_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana_managed_private_endpoint", "test")
	r := DashboardGrafanaManagedPrivateEndpointResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"grafana_name":        {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_dashboard_grafana_managed_private_endpoint.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_dashboard_grafana_managed_private_endpoint.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dashboard_grafana_managed_private_endpoint.test", tfjsonpath.New("grafana_name"), tfjsonpath.New("grafana_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dashboard_grafana_managed_private_endpoint.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("grafana_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dashboard_grafana_managed_private_endpoint.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("grafana_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
