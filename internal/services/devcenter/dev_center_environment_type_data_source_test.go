// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DevCenterEnvironmentTypeDataSource struct{}

func TestAccDevCenterEnvironmentTypeDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_center_environment_type", "test")
	r := DevCenterEnvironmentTypeDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dev_center_id").Exists(),
			),
		},
	})
}

func (d DevCenterEnvironmentTypeDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dev_center_environment_type" "test" {
  name          = azurerm_dev_center_environment_type.test.name
  dev_center_id = azurerm_dev_center_environment_type.test.dev_center_id
}
`, DevCenterEnvironmentTypeTestResource{}.basic(data))
}
