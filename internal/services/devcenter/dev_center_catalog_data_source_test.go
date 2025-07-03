// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DevCenterCatalogDataSource struct{}

func TestAccDevCenterCatalogDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_center_catalog", "test")
	r := DevCenterCatalogDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dev_center_id").Exists(),
				check.That(data.ResourceName).Key("catalog_github.#").Exists(),
			),
		},
	})
}

func (d DevCenterCatalogDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dev_center_catalog" "test" {
  name          = azurerm_dev_center_catalog.test.name
  dev_center_id = azurerm_dev_center_catalog.test.dev_center_id
}
`, DevCenterCatalogsResource{}.basic(data))
}
