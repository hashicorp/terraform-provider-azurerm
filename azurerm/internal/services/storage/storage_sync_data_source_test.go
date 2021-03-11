package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageSyncDataSource struct{}

func TestAccDataSourceStorageSync_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_sync", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: StorageSyncDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("incoming_traffic_policy").Exists(),
				check.That(data.ResourceName).Key("tags.%").Exists(),
			),
		},
	})
}

func (d StorageSyncDataSource) basic(data acceptance.TestData) string {
	basic := StorageSyncResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_sync" "test" {
  name                = azurerm_storage_sync.test.name
  resource_group_name = azurerm_storage_sync.test.resource_group_name
}
`, basic)
}
