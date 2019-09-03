package azurerm

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceArmStorageAccountBlobContainerSas_basic(t *testing.T) {
	dataSourceName := "data.azurerm_storage_account_blob_container_sas.test"
	rInt := tf.AccRandTimeInt()
	rString := acctest.RandString(4)
	location := testLocation()
	utcNow := time.Now().UTC()
	startDate := utcNow.Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccountBlobContainerSas_basic(rInt, rString, location, startDate, endDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "https_only", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "start", startDate),
					resource.TestCheckResourceAttr(dataSourceName, "expiry", endDate),
					resource.TestCheckResourceAttr(dataSourceName, "ip_address", "168.1.5.65"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.read", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.add", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.create", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.write", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.delete", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.list", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "cache_control", "max-age=5"),
					resource.TestCheckResourceAttr(dataSourceName, "content_disposition", "inline"),
					resource.TestCheckResourceAttr(dataSourceName, "content_encoding", "deflate"),
					resource.TestCheckResourceAttr(dataSourceName, "content_language", "en-US"),
					resource.TestCheckResourceAttr(dataSourceName, "content_type", "application/json"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sas"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageAccountBlobContainerSas_basic(rInt int, rString string, location string, startDate string, endDate string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestsa-%d"
  location = "%s"
}

resource "azurerm_storage_account" "storage" {
  name                = "acctestsads%s"
  resource_group_name = "${azurerm_resource_group.rg.name}"

  location                 = "${azurerm_resource_group.rg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "container" {
  name                  = "sas-test"
  resource_group_name   = "${azurerm_resource_group.rg.name}"
  storage_account_name  = "${azurerm_storage_account.storage.name}"
  container_access_type = "private"
}

data "azurerm_storage_account_blob_container_sas" "test" {
  connection_string = "${azurerm_storage_account.storage.primary_connection_string}"
  container_name    = "${azurerm_storage_container.container.name}"
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
`, rInt, location, rString, startDate, endDate)
}

func TestAccDataSourceArmStorageAccountBlobContainerSas_permissionsString(t *testing.T) {
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
		result := buildContainerPermissionsString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}
