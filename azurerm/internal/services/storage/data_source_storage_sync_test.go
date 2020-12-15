package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageSync_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_sync", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageSync_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "incoming_traffic_policy"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tags.%"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageSync_basic(data acceptance.TestData) string {
	basic := testAccAzureRMStorageSync_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_sync" "test" {
  name                = azurerm_storage_sync.test.name
  resource_group_name = azurerm_storage_sync.test.resource_group_name
}
`, basic)
}
