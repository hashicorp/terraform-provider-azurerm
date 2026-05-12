// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccRouteServerBgpConnection_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_server_bgp_connection", "test")
	r := RouteServerBgpConnectionResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"hub_name":            {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_route_server_bgp_connection.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_route_server_bgp_connection.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_route_server_bgp_connection.test", tfjsonpath.New("hub_name"), tfjsonpath.New("route_server_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_route_server_bgp_connection.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("route_server_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_route_server_bgp_connection.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("route_server_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
