// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DevCenterNetworkConnectionTestResource struct{}

func TestAccDevCenterNetworkConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_network_connection", "test")
	r := DevCenterNetworkConnectionTestResource{}

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

func TestAccDevCenterNetworkConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_network_connection", "test")
	r := DevCenterNetworkConnectionTestResource{}

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

func TestAccDevCenterNetworkConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_network_connection", "test")
	r := DevCenterNetworkConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("domain_password"),
	})
}

func TestAccDevCenterNetworkConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_network_connection", "test")
	r := DevCenterNetworkConnectionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("domain_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("domain_password"),
	})
}

func (r DevCenterNetworkConnectionTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkconnections.ParseNetworkConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.NetworkConnections.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DevCenterNetworkConnectionTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_network_connection" "test" {
  name                = "acctest-dcnc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  domain_join_type    = "AzureADJoin"
  subnet_id           = azurerm_subnet.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterNetworkConnectionTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_network_connection" "import" {
  name                = azurerm_dev_center_network_connection.test.name
  resource_group_name = azurerm_dev_center_network_connection.test.resource_group_name
  location            = azurerm_dev_center_network_connection.test.location
  domain_join_type    = azurerm_dev_center_network_connection.test.domain_join_type
  subnet_id           = azurerm_dev_center_network_connection.test.subnet_id
}
`, r.basic(data))
}

func (r DevCenterNetworkConnectionTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "network" {
  name     = "acctestrg-dcncn-%d"
  location = "%s"
}

resource "azurerm_dev_center_network_connection" "test" {
  name                = "acctest-dcnc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  domain_join_type    = "HybridAzureADJoin"
  subnet_id           = azurerm_subnet.test.id
  domain_name         = "never.gonna.shut.you.down"
  domain_username     = "tfuser@microsoft.com"
  domain_password     = "P@ssW0RD7890"
  organization_unit   = "OU=Sales,DC=Fabrikam,DC=com"

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DevCenterNetworkConnectionTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "network" {
  name     = "acctestrg-dcncn-%d"
  location = "%s"
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_dev_center_network_connection" "test" {
  name                = "acctest-dcnc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  domain_join_type    = "HybridAzureADJoin"
  subnet_id           = azurerm_subnet.test2.id
  domain_name         = "never2.gonna.shut.you.down"
  domain_username     = "tfuser2@microsoft.com"
  domain_password     = "P@ssW0RD7891"
  organization_unit   = "OU=SaleStores,DC=Fabrikam,DC=com"

  tags = {
    ENV = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DevCenterNetworkConnectionTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-dcnc-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
