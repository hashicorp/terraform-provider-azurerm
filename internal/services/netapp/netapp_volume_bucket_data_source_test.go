// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppVolumeBucketDataSource struct{}

func TestAccDataSourceNetAppVolumeBucket_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_volume_bucket", "test")
	d := NetAppVolumeBucketDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("path").HasValue("/"),
				check.That(data.ResourceName).Key("permissions").HasValue("ReadOnly"),
			),
		},
	})
}

func (d NetAppVolumeBucketDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_volume_bucket" "test" {
  name                = azurerm_netapp_volume_bucket.test.name
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.test.name
}
`, NetAppVolumeBucketResource{}.basic(data))
}
