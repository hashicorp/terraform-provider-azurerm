package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageAccountStaticWebsiteDataSource struct{}

func TestAccStorageAccountStaticWebsiteDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account_static_website", "test")
	r := StorageAccountStaticWebsiteDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("error_404_document").HasValue("sadpanda_2.html"),
				check.That(data.ResourceName).Key("index_document").HasValue("index_2.html"),
			),
		},
	})
}

func (d StorageAccountStaticWebsiteDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_storage_account_static_website" "test" {
  storage_account_id = azurerm_storage_account_static_website.test.storage_account_id
}
`, d.template(data))
}

func (d StorageAccountStaticWebsiteDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account_static_website" "test" {
  storage_account_id = azurerm_storage_account.test.id
  error_404_document = "sadpanda_2.html"
  index_document     = "index_2.html"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
