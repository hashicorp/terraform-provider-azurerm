// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package discovery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DiscoveryBookshelfPrivateLinkResourcesDataSource struct{}

func TestAccDiscoveryBookshelfPrivateLinkResourcesDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_discovery_bookshelf_private_link_resources", "test")
	d := DiscoveryBookshelfPrivateLinkResourcesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("private_link_resources.#").Exists(),
				check.That(data.ResourceName).Key("private_link_resources.0.group_id").Exists(),
			),
		},
	})
}

func (d DiscoveryBookshelfPrivateLinkResourcesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_discovery_bookshelf_private_link_resources" "test" {
  bookshelf_id = azurerm_discovery_bookshelf.test.id
}
`, DiscoveryBookshelfResource{}.basic(data))
}
