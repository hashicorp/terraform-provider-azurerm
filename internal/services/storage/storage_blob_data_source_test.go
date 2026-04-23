// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type StorageBlobDataSource struct{}

func TestAccDataSourceStorageBlob_basicDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("Deprecated test skipping in 5.0")
	}
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}

	data := acceptance.BuildTestData(t, "data.azurerm_storage_blob", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageBlobDataSource{}.basicWithDataSourceDeprecated(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("encryption_scope").HasValue(fmt.Sprintf("acctestEScontainer%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("type").HasValue("Block"),
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.k1").HasValue("v1"),
				check.That(data.ResourceName).Key("metadata.k2").HasValue("v2"),
			),
		},
	})
}

func TestAccDataSourceStorageBlob_basic(t *testing.T) {
	if !features.FivePointOh() {
		t.Skip("5.0 test skipping in 4.x")
	}
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}

	data := acceptance.BuildTestData(t, "data.azurerm_storage_blob", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageBlobDataSource{}.basicWithDataSource(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("encryption_scope").HasValue(fmt.Sprintf("acctestEScontainer%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("type").HasValue("Block"),
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.k1").HasValue("v1"),
				check.That(data.ResourceName).Key("metadata.k2").HasValue("v2"),
			),
		},
	})
}

func (d StorageBlobDataSource) basicDeprecated(data acceptance.TestData, fileName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "blobdstest-%[1]s"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsadsc%[1]s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[3]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_container" "test" {
  name                  = "containerdstest-%[1]s"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  encryption_scope       = azurerm_storage_encryption_scope.test.name
  type                   = "Block"
  source                 = "%[4]s"

  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}
`, data.RandomString, data.Locations.Primary, data.RandomInteger, fileName)
}

func (d StorageBlobDataSource) basic(data acceptance.TestData, fileName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "blobdstest-%[1]s"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsadsc%[1]s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[3]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_container" "test" {
  name                  = "containerdstest-%[1]s"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_container_id   = azurerm_storage_container.test.id
  encryption_scope       = azurerm_storage_encryption_scope.test.name
  type                   = "Block"
  source                 = "%[4]s"

  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}
`, data.RandomString, data.Locations.Primary, data.RandomInteger, fileName)
}

func (d StorageBlobDataSource) basicWithDataSourceDeprecated(data acceptance.TestData, fileName string) string {
	config := d.basicDeprecated(data, fileName)
	return fmt.Sprintf(`
%s

data "azurerm_storage_blob" "test" {
  name                   = azurerm_storage_blob.test.name
  storage_account_name   = azurerm_storage_blob.test.storage_account_name
  storage_container_name = azurerm_storage_blob.test.storage_container_name
}
`, config)
}

func (d StorageBlobDataSource) basicWithDataSource(data acceptance.TestData, fileName string) string {
	config := d.basic(data, fileName)
	return fmt.Sprintf(`
%s

data "azurerm_storage_blob" "test" {
  name                   = azurerm_storage_blob.test.name
    storage_container_id = azurerm_storage_blob.test.storage_container_id
}
`, config)
}
