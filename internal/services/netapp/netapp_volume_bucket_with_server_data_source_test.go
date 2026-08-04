// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppVolumeBucketWithServerDataSource struct{}

func TestAccDataSourceNetAppVolumeBucketWithServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_volume_bucket_with_server", "test")
	d := NetAppVolumeBucketWithServerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("path").HasValue("/"),
				check.That(data.ResourceName).Key("permissions").HasValue("ReadOnly"),
				check.That(data.ResourceName).Key("server.0.fqdn").Exists(),
			),
		},
	})
}

func (d NetAppVolumeBucketWithServerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_volume_bucket_with_server" "test" {
  name             = azurerm_netapp_volume_bucket_with_server.test.name
  netapp_volume_id = azurerm_netapp_volume.test.id
}
`, NetAppVolumeBucketWithServerResource{}.basic(data))
}
