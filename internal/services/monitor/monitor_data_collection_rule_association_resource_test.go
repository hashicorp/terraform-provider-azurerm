// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionruleassociations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorDataCollectionRuleAssociationResource struct{}

func (r MonitorDataCollectionRuleAssociationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := datacollectionruleassociations.ParseScopedDataCollectionRuleAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Monitor.DataCollectionRuleAssociationsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccMonitorDataCollectionRuleAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule_association", "test")
	r := MonitorDataCollectionRuleAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDataCollectionRuleAssociation_basicEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule_association", "test")
	r := MonitorDataCollectionRuleAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDataCollectionRuleAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule_association", "test")
	r := MonitorDataCollectionRuleAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMonitorDataCollectionRuleAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule_association", "test")
	r := MonitorDataCollectionRuleAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDataCollectionRuleAssociation_updateDataCollectionRuleId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule_association", "test")
	r := MonitorDataCollectionRuleAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDataCollectionRuleId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDataCollectionRuleAssociation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule_association", "test")
	r := MonitorDataCollectionRuleAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MonitorDataCollectionRuleAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_rule_association" "test" {
  name                    = "test-dcra-%[2]d"
  target_resource_id      = azurerm_linux_virtual_machine.test.id
  data_collection_rule_id = azurerm_monitor_data_collection_rule.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleAssociationResource) basicEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_endpoint" "test" {
  name                = "acctestmdcr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_monitor_data_collection_rule_association" "test" {
  target_resource_id          = azurerm_linux_virtual_machine.test.id
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleAssociationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_rule_association" "test" {
  name                    = "test-dcra-%[2]d"
  target_resource_id      = azurerm_linux_virtual_machine.test.id
  data_collection_rule_id = azurerm_monitor_data_collection_rule.test.id
  description             = "test dcra"
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleAssociationResource) updateDataCollectionRuleId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_rule" "test2" {
  name                = "acctestmdcr2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  destinations {
    azure_monitor_metrics {
      name = "test-destination-metrics"
    }
  }
  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination-metrics"]
  }
}

resource "azurerm_monitor_data_collection_rule_association" "test" {
  name                    = "test-dcra-%[2]d"
  target_resource_id      = azurerm_linux_virtual_machine.test.id
  data_collection_rule_id = azurerm_monitor_data_collection_rule.test2.id
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_data_collection_rule_association" "import" {
  name                    = azurerm_monitor_data_collection_rule_association.test.name
  target_resource_id      = azurerm_monitor_data_collection_rule_association.test.target_resource_id
  data_collection_rule_id = azurerm_monitor_data_collection_rule_association.test.data_collection_rule_id
}
`, r.basic(data))
}

func (r MonitorDataCollectionRuleAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-DCRA-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "network-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "nic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "machine-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_B1ls"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_password = "test-Password@7890"

  disable_password_authentication = false

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
  lifecycle {
    ignore_changes = [tags, identity]
  }
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctestmdcr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  destinations {
    azure_monitor_metrics {
      name = "test-destination-metrics"
    }
  }
  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination-metrics"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
