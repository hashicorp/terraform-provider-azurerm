// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type QuotaGroupDataSource struct{}

func TestAccQuotaGroupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_quota_group", "test")
	d := QuotaGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrSet(data.ResourceName, "id"),
			),
		},
	})
}

func TestAccQuotaGroupDataSource_withLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_quota_group", "test")
	d := QuotaGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.withLocation(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrSet(data.ResourceName, "id"),
				acceptance.TestCheckResourceAttrSet(data.ResourceName, "quota_request.#"),
			),
		},
	})
}

func (d QuotaGroupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_quota_group" "test" {
  name                = azurerm_quota_group.test.name
  management_group_id = azurerm_quota_group.test.management_group_id
}
`, QuotaGroupResource{}.basic(data))
}

func (d QuotaGroupDataSource) withLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_quota_group" "test" {
  name                   = azurerm_quota_group.test.name
  management_group_id    = azurerm_quota_group.test.management_group_id
  location               = "%s"
  resource_provider_name = "Microsoft.Compute"
}
`, QuotaGroupResource{}.complete(data), data.Locations.Primary)
}
