// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkSiteDataSource struct{}

func TestAccMobileNetworkSiteDataSource_complete(t *testing.T) {
	t.Skipf("Skipping since Mobile Network is deprecated and will be removed in 5.0")
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")

	d := MobileNetworkSiteDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (r MobileNetworkSiteDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_mobile_network_site" "test" {
  name              = azurerm_mobile_network_site.test.name
  mobile_network_id = azurerm_mobile_network_site.test.mobile_network_id
}
`, MobileNetworkSiteResource{}.complete(data))
}
