package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageSyncGroupDataSource struct{}

func TestAccDataSourceStorageSyncGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_sync_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: StorageSyncGroupDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_sync_id").Exists(),
			),
		},
	})
}

func (d StorageSyncGroupDataSource) basic(data acceptance.TestData) string {
	basic := StorageSyncGroupResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_sync_group" "test" {
  name            = azurerm_storage_sync_group.test.name
  storage_sync_id = azurerm_storage_sync.test.id
}
`, basic)
}
