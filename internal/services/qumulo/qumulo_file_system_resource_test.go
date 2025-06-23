// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package qumulo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2024-06-19/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FileSystemTestResource struct{}

func TestAccQumuloFileSystem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_qumulo_file_system", "test")
	r := FileSystemTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "email"),
	})
}

func TestAccQumuloFileSystem_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_qumulo_file_system", "test")
	r := FileSystemTestResource{}

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

func TestAccQumuloFileSystem_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_qumulo_file_system", "test")
	r := FileSystemTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "email"),
	})
}

func TestAccQumuloFileSystem_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_qumulo_file_system", "test")
	r := FileSystemTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "email"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "email"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "email"),
	})
}

func (r FileSystemTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := filesystems.ParseFileSystemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Qumulo.FileSystemsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r FileSystemTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_qumulo_file_system" "test" {
  name                = "acctest-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_password      = ")^X#ZX#JRyIY}t9"
  storage_sku         = "Cold_LRS"
  subnet_id           = azurerm_subnet.test.id
  email               = "test@test.com"
  zone                = "1"
}
`, r.template(data), data.RandomString)
}

func (r FileSystemTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_qumulo_file_system" "test" {
  name                = "acctest-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_password      = ")^X#ZX#JRyIY}t9"
  publisher_id        = "qumulo1584033880660"
  storage_sku         = "Cold_LRS"
  subnet_id           = azurerm_subnet.test.id
  email               = "test@test.com"
  zone                = "1"
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
}
`, r.template(data), data.RandomString)
}

func (r FileSystemTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_qumulo_file_system" "import" {
  name                = azurerm_qumulo_file_system.test.name
  resource_group_name = azurerm_qumulo_file_system.test.resource_group_name
  location            = azurerm_qumulo_file_system.test.location
  admin_password      = azurerm_qumulo_file_system.test.admin_password
  storage_sku         = azurerm_qumulo_file_system.test.storage_sku
  subnet_id           = azurerm_qumulo_file_system.test.subnet_id
  email               = azurerm_qumulo_file_system.test.email
  zone                = azurerm_qumulo_file_system.test.zone
}
`, r.basic(data))
}

func (r FileSystemTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_qumulo_file_system" "test" {
  name                = "acctest-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_password      = ")^X#ZX#JRyIY}t9"
  offer_id            = "qumulo-saas-mpp"
  plan_id             = "azure-native-qumulo-v3"
  publisher_id        = "qumulo1584033880660"
  storage_sku         = "Cold_LRS"
  subnet_id           = azurerm_subnet.test.id
  email               = "test@test.com"
  zone                = "1"
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
}
`, r.template(data), data.RandomString)
}

func (r FileSystemTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-qumulo-%s"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[1]s"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%[1]s"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Qumulo.Storage/fileSystems"
    }
  }
}
`, data.RandomString, data.Locations.Primary)
}
