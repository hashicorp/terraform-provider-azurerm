// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DevCenterProjectPoolDataSource struct{}

func TestAccDevCenterProjectPoolDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_center_project_pool", "test")
	r := DevCenterProjectPoolDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("dev_center_project_id").Exists(),
				check.That(data.ResourceName).Key("dev_box_definition_name").Exists(),
				check.That(data.ResourceName).Key("local_administrator_enabled").Exists(),
				check.That(data.ResourceName).Key("dev_center_attached_network_name").Exists(),
			),
		},
	})
}

func (d DevCenterProjectPoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dev_center_project_pool" "test" {
  name                  = azurerm_dev_center_project_pool.test.name
  dev_center_project_id = azurerm_dev_center_project_pool.test.dev_center_project_id
}
`, DevCenterProjectPoolTestResource{}.basic(data))
}
