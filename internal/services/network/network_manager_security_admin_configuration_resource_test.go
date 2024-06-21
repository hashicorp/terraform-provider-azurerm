// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/securityadminconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerSecurityAdminConfigurationResource struct{}

func testAccNetworkManagerSecurityAdminConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_security_admin_configuration", "test")
	r := ManagerSecurityAdminConfigurationResource{}
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

func testAccNetworkManagerSecurityAdminConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_security_admin_configuration", "test")
	r := ManagerSecurityAdminConfigurationResource{}
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

func testAccNetworkManagerSecurityAdminConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_security_admin_configuration", "test")
	r := ManagerSecurityAdminConfigurationResource{}
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

func testAccNetworkManagerSecurityAdminConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_security_admin_configuration", "test")
	r := ManagerSecurityAdminConfigurationResource{}
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

func (r ManagerSecurityAdminConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securityadminconfigurations.ParseSecurityAdminConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.SecurityAdminConfigurations
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerSecurityAdminConfigurationResource) template(data acceptance.TestData) string {
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

resource "azurerm_virtual_network" "test" {
  name                    = "acctest-vnet-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  address_space           = ["10.0.0.0/16"]
  flow_timeout_in_minutes = 10
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ManagerSecurityAdminConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_security_admin_configuration" "test" {
  name               = "acctest-nmsac-%d"
  network_manager_id = azurerm_network_manager.test.id
}
`, template, data.RandomInteger)
}

func (r ManagerSecurityAdminConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_security_admin_configuration" "import" {
  name               = azurerm_network_manager_security_admin_configuration.test.name
  network_manager_id = azurerm_network_manager_security_admin_configuration.test.network_manager_id
}
`, config)
}

func (r ManagerSecurityAdminConfigurationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_security_admin_configuration" "test" {
  name                                          = "acctest-nmsac-%d"
  network_manager_id                            = azurerm_network_manager.test.id
  description                                   = "test"
  apply_on_network_intent_policy_based_services = ["None"]
}
`, template, data.RandomInteger)
}

func (r ManagerSecurityAdminConfigurationResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_security_admin_configuration" "test" {
  name                                          = "acctest-nmsac-%d"
  network_manager_id                            = azurerm_network_manager.test.id
  description                                   = "update"
  apply_on_network_intent_policy_based_services = ["AllowRulesOnly"]
}
`, template, data.RandomInteger)
}
