// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkDDoSProtectionPlanResource struct{}

// NOTE: this is a test group to avoid each test case to run in parallel, since Azure only allows one DDoS Protection
// Plan per region.
func TestAccNetworkDDoSProtectionPlan(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"normal": {
			"basic":          testAccNetworkDDoSProtectionPlan_basic,
			"requiresImport": testAccNetworkDDoSProtectionPlan_requiresImport,
			"withTags":       testAccNetworkDDoSProtectionPlan_withTags,
			"disappears":     testAccNetworkDDoSProtectionPlan_disappears,
		},
		"datasource": {
			"basic": testAccNetworkDDoSProtectionPlanDataSource_basic,
		},
	}

	for group, steps := range testCases {
		t.Run(group, func(t *testing.T) {
			for name, tc := range steps {
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccNetworkDDoSProtectionPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")
	r := NetworkDDoSProtectionPlanResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_ids.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkDDoSProtectionPlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")
	r := NetworkDDoSProtectionPlanResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_ddos_protection_plan"),
		},
	})
}

func testAccNetworkDDoSProtectionPlan_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")
	r := NetworkDDoSProtectionPlanResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withUpdatedTagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Staging"),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkDDoSProtectionPlan_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")
	r := NetworkDDoSProtectionPlanResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicConfig,
			TestResource: r,
		}),
	})
}

func (t NetworkDDoSProtectionPlanResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.DdosProtectionPlans.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NetworkDDoSProtectionPlanResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := client.Network.DdosProtectionPlans.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (NetworkDDoSProtectionPlanResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r NetworkDDoSProtectionPlanResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_ddos_protection_plan" "import" {
  name                = azurerm_network_ddos_protection_plan.test.name
  location            = azurerm_network_ddos_protection_plan.test.location
  resource_group_name = azurerm_network_ddos_protection_plan.test.resource_group_name
}
`, r.basicConfig(data))
}

func (NetworkDDoSProtectionPlanResource) withTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (NetworkDDoSProtectionPlanResource) withUpdatedTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
