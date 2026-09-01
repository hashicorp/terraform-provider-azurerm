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

func TestAccNetworkWatcherFlowLog_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	checkedFields := map[string]struct{}{
		"name":                 {},
		"network_watcher_name": {},
		"resource_group_name":  {},
		"subscription_id":      {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_network_watcher_flow_log.test", checkedFields),
			customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_watcher_flow_log.test", tfjsonpath.New("name"), tfjsonpath.New("id")),
			customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_watcher_flow_log.test", tfjsonpath.New("network_watcher_name"), tfjsonpath.New("id")),
			customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_watcher_flow_log.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("id")),
			customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_watcher_flow_log.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
