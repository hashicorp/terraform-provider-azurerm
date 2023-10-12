// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedPrivateEndpointResource struct{}

func TestAccDataFactoryManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}

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

func TestAccDataFactoryManagedPrivateEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}

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

func TestAccDataFactoryManagedPrivateEndpoint_privateServiceLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateServiceLink(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagedPrivateEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	virtualNetworkId := managedprivateendpoints.NewManagedVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.ManagedVirtualNetworkName)
	iter, err := client.DataFactory.ManagedPrivateEndpoints.ListByFactoryComplete(ctx, virtualNetworkId)
	if err != nil {
		return nil, fmt.Errorf("listing %s: %+v", id, err)
	}

	for _, item := range iter.Items {
		if item.Name != nil && *item.Name == id.ManagedPrivateEndpointName {
			return pointer.To(true), nil
		}
	}

	return pointer.To(false), nil
}

func (r ManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory_managed_private_endpoint" "test" {
  name               = "acctestEndpoint%d"
  data_factory_id    = azurerm_data_factory.test.id
  target_resource_id = azurerm_storage_account.test.id
  subresource_name   = "blob"
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ManagedPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_data_factory_managed_private_endpoint" "import" {
  name               = azurerm_data_factory_managed_private_endpoint.test.name
  data_factory_id    = azurerm_data_factory_managed_private_endpoint.test.data_factory_id
  target_resource_id = azurerm_data_factory_managed_private_endpoint.test.target_resource_id
  subresource_name   = azurerm_data_factory_managed_private_endpoint.test.subresource_name
}
`, config)
}

func (r ManagedPrivateEndpointResource) privateServiceLink(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                                           = "acctsub-%[2]d"
  resource_group_name                            = azurerm_resource_group.test.name
  virtual_network_name                           = azurerm_virtual_network.test.name
  address_prefixes                               = ["10.0.2.0/24"]
  enforce_private_link_endpoint_network_policies = true
  enforce_private_link_service_network_policies  = true
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%[2]d"
    subnet_id = azurerm_subnet.test.id
    primary   = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}

resource "azurerm_data_factory_managed_private_endpoint" "test" {
  name               = "acctestEndpoint%[2]d"
  data_factory_id    = azurerm_data_factory.test.id
  target_resource_id = azurerm_private_link_service.test.id
  fqdns              = ["a.a.a.a.a"]
}
`, template, data.RandomInteger)
}

func (r ManagedPrivateEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-adf-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                            = "acctestdf%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  managed_virtual_network_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
