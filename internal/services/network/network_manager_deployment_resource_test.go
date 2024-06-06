// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerDeploymentResource struct{}

func testAccNetworkManagerDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_deployment", "test")
	r := ManagerDeploymentResource{}
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

func testAccNetworkManagerDeployment_basicAdmin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_deployment", "test")
	r := ManagerDeploymentResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAdmin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerDeployment_withTriggers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_deployment", "test")
	r := ManagerDeploymentResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTriggers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// triggers is an arbitrary list(string) which
		// is not known at the backend API
		data.ImportStep("triggers"),
	})
}

func testAccNetworkManagerDeployment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_deployment", "test")
	r := ManagerDeploymentResource{}
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

func testAccNetworkManagerDeployment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_deployment", "test")
	r := ManagerDeploymentResource{}
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

func testAccNetworkManagerDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_deployment", "test")
	r := ManagerDeploymentResource{}
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
	})
}

func (r ManagerDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkManagerDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.NetworkManagers
	listParam := networkmanagers.NetworkManagerDeploymentStatusParameter{
		Regions:         &[]string{azure.NormalizeLocation(id.Location)},
		DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(id.ScopeAccess)},
	}
	networkManagerId := networkmanagers.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName)

	resp, err := client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, listParam)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("unexpected null model %s", *id)
	}

	return utils.Bool(resp.Model.Value != nil && len(*resp.Model.Value) != 0 && *(*resp.Model.Value)[0].ConfigurationIds != nil && len(*(*resp.Model.Value)[0].ConfigurationIds) != 0), nil
}

func (r ManagerDeploymentResource) template(data acceptance.TestData) string {
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
  name                = "acctest-nm-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["SecurityAdmin", "Connectivity"]
}

resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_virtual_network" "test" {
  name                    = "acctest-vnet-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  address_space           = ["10.0.0.0/16"]
  flow_timeout_in_minutes = 10
}

resource "azurerm_network_manager_connectivity_configuration" "test" {
  name                  = "acctest-nmcc-%[1]d"
  network_manager_id    = azurerm_network_manager.test.id
  connectivity_topology = "HubAndSpoke"
  applies_to_group {
    group_connectivity = "None"
    network_group_id   = azurerm_network_manager_network_group.test.id
  }
  hub {
    resource_id   = azurerm_virtual_network.test.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagerDeploymentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_deployment" "test" {
  network_manager_id = azurerm_network_manager.test.id
  location           = "eastus"
  scope_access       = "Connectivity"
  configuration_ids  = [azurerm_network_manager_connectivity_configuration.test.id]
}
`, template)
}

func (r ManagerDeploymentResource) basicAdmin(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_security_admin_configuration" "test" {
  name               = "acctest-nmsac-%d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_admin_rule_collection" "test" {
  name                            = "acctest-nmarc-%[2]d"
  security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
  network_group_ids               = [azurerm_network_manager_network_group.test.id]
}

resource "azurerm_network_manager_admin_rule" "test" {
  name                     = "acctest-nmar-%[2]d"
  admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  action                   = "Deny"
  description              = "test"
  direction                = "Inbound"
  priority                 = 1
  protocol                 = "Tcp"
  source_port_ranges       = ["80"]
  destination_port_ranges  = ["80"]
  source {
    address_prefix_type = "ServiceTag"
    address_prefix      = "Internet"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "*"
  }
}

resource "azurerm_network_manager_deployment" "test" {
  network_manager_id = azurerm_network_manager.test.id
  location           = "eastus"
  scope_access       = "SecurityAdmin"
  configuration_ids  = [azurerm_network_manager_security_admin_configuration.test.id]
  depends_on         = [azurerm_network_manager_admin_rule.test]
}
`, template, data.RandomInteger)
}

func (r ManagerDeploymentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_deployment" "import" {
  network_manager_id = azurerm_network_manager_deployment.test.network_manager_id
  location           = azurerm_network_manager_deployment.test.location
  scope_access       = azurerm_network_manager_deployment.test.scope_access
  configuration_ids  = azurerm_network_manager_deployment.test.configuration_ids
}
`, config)
}

func (r ManagerDeploymentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_connectivity_configuration" "test2" {
  name                  = "acctest-nmcc2-%d"
  network_manager_id    = azurerm_network_manager.test.id
  connectivity_topology = "HubAndSpoke"
  applies_to_group {
    group_connectivity = "None"
    network_group_id   = azurerm_network_manager_network_group.test.id
  }
  hub {
    resource_id   = azurerm_virtual_network.test.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}

resource "azurerm_network_manager_deployment" "test" {
  network_manager_id = azurerm_network_manager.test.id
  location           = "eastus"
  scope_access       = "Connectivity"
  configuration_ids = [
    azurerm_network_manager_connectivity_configuration.test.id,
    azurerm_network_manager_connectivity_configuration.test2.id
  ]
}
`, template, data.RandomInteger)
}

func (r ManagerDeploymentResource) withTriggers(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_security_admin_configuration" "test" {
  name               = "acctest-nmsac-%d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_admin_rule_collection" "test" {
  name                            = "acctest-nmarc-%[2]d"
  security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.test.id
  network_group_ids               = [azurerm_network_manager_network_group.test.id]
}

resource "azurerm_network_manager_admin_rule" "test" {
  name                     = "acctest-nmar-%[2]d"
  admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.test.id
  action                   = "Deny"
  description              = "test"
  direction                = "Inbound"
  priority                 = 1
  protocol                 = "Tcp"
  source_port_ranges       = ["80"]
  destination_port_ranges  = ["80"]
  source {
    address_prefix_type = "ServiceTag"
    address_prefix      = "Internet"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "*"
  }
}

resource "azurerm_network_manager_deployment" "test" {
  network_manager_id = azurerm_network_manager.test.id
  location           = "eastus"
  scope_access       = "SecurityAdmin"
  configuration_ids  = [azurerm_network_manager_security_admin_configuration.test.id]
  depends_on         = [azurerm_network_manager_admin_rule.test]
  triggers = {
    source_port_ranges = join(",", azurerm_network_manager_admin_rule.test.source_port_ranges)
  }
}

`, template, data.RandomInteger)
}

func (r ManagerDeploymentResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_connectivity_configuration" "test2" {
  name                  = "acctest-nmcc2-%d"
  network_manager_id    = azurerm_network_manager.test.id
  connectivity_topology = "HubAndSpoke"
  applies_to_group {
    group_connectivity = "None"
    network_group_id   = azurerm_network_manager_network_group.test.id
  }
  hub {
    resource_id   = azurerm_virtual_network.test.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}

resource "azurerm_network_manager_deployment" "test" {
  network_manager_id = azurerm_network_manager.test.id
  location           = "eastus"
  scope_access       = "Connectivity"
  configuration_ids = [
    azurerm_network_manager_connectivity_configuration.test2.id
  ]
}
`, template, data.RandomInteger)
}
