package storage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageBlobDataSource struct{}

func TestAccDataSourceStorageBlob_basic(t *testing.T) {
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
				check.That(data.ResourceName).Key("type").HasValue("Block"),
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.k1").HasValue("v1"),
				check.That(data.ResourceName).Key("metadata.k2").HasValue("v2"),
			),
		},
	})
}

func (d StorageBlobDataSource) basic(data acceptance.TestData, fileName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "blobdstest-%s"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsadsc%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "containerdstest-%s"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.vhd"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source                 = "%s"

  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}
`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString, fileName)
}

func (d StorageBlobDataSource) basicWithDataSource(data acceptance.TestData, fileName string) string {
	config := d.basic(data, fileName)
	return fmt.Sprintf(`
%s

data "azurerm_storage_blob" "test" {
  name                   = azurerm_storage_blob.test.name
  storage_account_name   = azurerm_storage_blob.test.storage_account_name
  storage_container_name = azurerm_storage_blob.test.storage_container_name
}
`, config)
}
