package storage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
)

type StorageAccountBlobContainerSASDataSource struct{}

func TestAccDataSourceStorageAccountBlobContainerSas_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account_blob_container_sas", "test")
	utcNow := time.Now().UTC()
	startDate := utcNow.Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: StorageAccountBlobContainerSASDataSource{}.basic(data, startDate, endDate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("https_only").HasValue("true"),
				check.That(data.ResourceName).Key("start").HasValue(startDate),
				check.That(data.ResourceName).Key("expiry").HasValue(endDate),
				check.That(data.ResourceName).Key("ip_address").HasValue("168.1.5.65"),
				check.That(data.ResourceName).Key("permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.read").HasValue("true"),
				check.That(data.ResourceName).Key("permissions.0.add").HasValue("true"),
				check.That(data.ResourceName).Key("permissions.0.create").HasValue("false"),
				check.That(data.ResourceName).Key("permissions.0.write").HasValue("false"),
				check.That(data.ResourceName).Key("permissions.0.delete").HasValue("true"),
				check.That(data.ResourceName).Key("permissions.0.list").HasValue("true"),
				check.That(data.ResourceName).Key("cache_control").HasValue("max-age=5"),
				check.That(data.ResourceName).Key("content_disposition").HasValue("inline"),
				check.That(data.ResourceName).Key("content_encoding").HasValue("deflate"),
				check.That(data.ResourceName).Key("content_language").HasValue("en-US"),
				check.That(data.ResourceName).Key("content_type").HasValue("application/json"),
				check.That(data.ResourceName).Key("sas").Exists(),
			),
		},
	})
}

func (d StorageAccountBlobContainerSASDataSource) basic(data acceptance.TestData, startDate string, endDate string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "storage" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.rg.name

  location                 = azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "container" {
  name                  = "sas-test"
  storage_account_name  = azurerm_storage_account.storage.name
  container_access_type = "private"
}

data "azurerm_storage_account_blob_container_sas" "test" {
  connection_string = azurerm_storage_account.storage.primary_connection_string
  container_name    = azurerm_storage_container.container.name
  https_only        = true

  ip_address = "168.1.5.65"

  start  = "%s"
  expiry = "%s"

  permissions {
    read   = true
    add    = true
    create = false
    write  = false
    delete = true
    list   = true
  }

  cache_control       = "max-age=5"
  content_disposition = "inline"
  content_encoding    = "deflate"
  content_language    = "en-US"
  content_type        = "application/json"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, startDate, endDate)
}

func TestAccDataSourceStorageAccountBlobContainerSas_permissionsString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"read": true}, "r"},
		{map[string]interface{}{"add": true}, "a"},
		{map[string]interface{}{"create": true}, "c"},
		{map[string]interface{}{"write": true}, "w"},
		{map[string]interface{}{"delete": true}, "d"},
		{map[string]interface{}{"list": true}, "l"},
		{map[string]interface{}{"add": true, "write": true, "read": true, "delete": true}, "rawd"},
	}

	for _, test := range testCases {
		result := storage.BuildContainerPermissionsString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}
