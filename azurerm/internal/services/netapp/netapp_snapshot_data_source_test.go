package netapp_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type NetAppSnapshotDataSource struct {
}

func TestAccDataSourceNetAppSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (NetAppSnapshotDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_snapshot" "test" {
  resource_group_name = azurerm_netapp_snapshot.test.resource_group_name
  account_name        = azurerm_netapp_snapshot.test.account_name
  pool_name           = azurerm_netapp_snapshot.test.pool_name
  volume_name         = azurerm_netapp_snapshot.test.volume_name
  name                = azurerm_netapp_snapshot.test.name
}
`, NetAppSnapshotResource{}.basic(data))
}
