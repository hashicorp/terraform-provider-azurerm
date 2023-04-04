package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppVolumeGroupDataSource struct{}

func TestAccNetAppVolumeGroupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_volume_group", "test")
	d := NetAppVolumeGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("application_type").HasValue("SAP-HANA"),
				check.That(data.ResourceName).Key("volume.1.volume_spec_name").HasValue("log"),
			),
		},
	})
}

func (d NetAppVolumeGroupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_volume_group" "test" {
  name                   = azurerm_netapp_volume_group.test.name
  resource_group_name    = azurerm_netapp_volume_group.test.resource_group_name
  account_name           = azurerm_netapp_volume_group.test.account_name
}
`, NetAppVolumeGroupResource{}.basic(data))
}
