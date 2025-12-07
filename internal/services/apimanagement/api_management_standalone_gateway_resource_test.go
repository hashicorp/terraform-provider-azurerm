// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementStandaloneGatewayResource struct{}

func TestAccApiManagementStandaloneGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_standalone_gateway", "test")
	r := ApiManagementStandaloneGatewayResource{}
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

func TestAccApiManagementStandaloneGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_standalone_gateway", "test")
	r := ApiManagementStandaloneGatewayResource{}
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

func TestAccApiManagementStandaloneGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_standalone_gateway", "test")
	r := ApiManagementStandaloneGatewayResource{}
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

func TestAccApiManagementStandaloneGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_standalone_gateway", "test")
	r := ApiManagementStandaloneGatewayResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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

func (r ApiManagementStandaloneGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apigateway.ParseGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiGatewayClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementStandaloneGatewayResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementStandaloneGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_standalone_gateway" "test" {
  name                = "acctest-aag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  sku {
    capacity = 1
    name     = "WorkspaceGatewayPremium"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementStandaloneGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_api_management_standalone_gateway" "import" {
  name                = azurerm_api_management_standalone_gateway.test.name
  resource_group_name = azurerm_api_management_standalone_gateway.test.resource_group_name
  location            = azurerm_api_management_standalone_gateway.test.location

  sku {
    capacity = 1
    name     = "WorkspaceGatewayPremium"
  }
}
`, r.basic(data))
}

func (r ApiManagementStandaloneGatewayResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "apim-delegation"
    service_delegation {
      name = "Microsoft.Web/serverFarms"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/action"
      ]
    }
  }
}

resource "azurerm_api_management_standalone_gateway" "test" {
  name                 = "acctest-aag-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  virtual_network_type = "External"
  backend_subnet_id    = azurerm_subnet.test.id

  sku {
    capacity = 1
    name     = "WorkspaceGatewayPremium"
  }

  tags = {
    key = "value"
  }

}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementStandaloneGatewayResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "apim-delegation"
    service_delegation {
      name = "Microsoft.Web/serverFarms"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/action"
      ]
    }
  }
}

resource "azurerm_api_management_standalone_gateway" "test" {
  name                 = "acctest-aag-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  virtual_network_type = "External"
  backend_subnet_id    = azurerm_subnet.test.id

  sku {
    capacity = 2
    name     = "WorkspaceGatewayPremium"
  }

  tags = {
    key = "value"
  }
}
`, r.template(data), data.RandomInteger)
}
