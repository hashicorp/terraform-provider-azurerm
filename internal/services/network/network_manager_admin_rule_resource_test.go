// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/adminrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerAdminRuleResource struct{}

func testAccNetworkManagerAdminRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerAdminRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccNetworkManagerAdminRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerAdminRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_admin_rule", "test")
	r := ManagerAdminRuleResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagerAdminRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := adminrules.ParseRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.AdminRules
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerAdminRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["SecurityAdmin"]
}

resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-%d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_security_admin_configuration" "test" {
  name               = "acctest-nmsac-%d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_admin_rule_collection" "test" {
  name                            = "acctest-nmarc-%d"
  security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
  network_group_ids               = [azurerm_network_manager_network_group.test.id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ManagerAdminRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_network_manager_admin_rule" "test" {
  name                     = "acctest-nmar-%d"
  admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  action                   = "Deny"
  direction                = "Outbound"
  protocol                 = "Tcp"
  priority                 = 1
}
`, template, data.RandomInteger)
}

func (r ManagerAdminRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_admin_rule" "import" {
  name                     = azurerm_network_manager_admin_rule.test.name
  admin_rule_collection_id = azurerm_network_manager_admin_rule.test.admin_rule_collection_id
  action                   = azurerm_network_manager_admin_rule.test.action
  direction                = azurerm_network_manager_admin_rule.test.direction
  priority                 = azurerm_network_manager_admin_rule.test.priority
  protocol                 = azurerm_network_manager_admin_rule.test.protocol
}
`, config)
}

func (r ManagerAdminRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_admin_rule" "test" {
  name                     = "acctest-nmar-%d"
  admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  action                   = "Deny"
  description              = "test admin rule"
  direction                = "Outbound"
  priority                 = 1
  protocol                 = "Tcp"
  source_port_ranges       = ["80", "22", "443"]
  destination_port_ranges  = ["80", "22"]
  source {
    address_prefix_type = "ServiceTag"
    address_prefix      = "Internet"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "*"
  }
}
`, template, data.RandomInteger)
}

func (r ManagerAdminRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_admin_rule" "test" {
  name                     = "acctest-nmar-%d"
  admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  action                   = "Allow"
  description              = "test"
  direction                = "Inbound"
  priority                 = 1234
  protocol                 = "Ah"
  source_port_ranges       = ["80", "1024-65535"]
  destination_port_ranges  = ["80"]
  source {
    address_prefix_type = "ServiceTag"
    address_prefix      = "ActionGroup"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "10.1.0.1"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "10.0.0.0/24"
  }
}
`, template, data.RandomInteger)
}
