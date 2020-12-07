package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageSyncGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_sync_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageSyncGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncGroupExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageSyncGroup_basic(data acceptance.TestData) string {
	basic := testAccAzureRMStorageSyncGroup_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_sync_group" "test" {
  name            = azurerm_storage_sync_group.test.name
  storage_sync_id = azurerm_storage_sync.test.id
}
`, basic)
}
