// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DevCenterDevBoxDefinitionDataSource struct{}

func TestAccDevCenterDevBoxDefinitionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_center_dev_box_definition", "test")
	r := DevCenterDevBoxDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dev_center_id").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("image_reference_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
			),
		},
	})
}

func (d DevCenterDevBoxDefinitionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dev_center_dev_box_definition" "test" {
  name          = azurerm_dev_center_dev_box_definition.test.name
  dev_center_id = azurerm_dev_center_dev_box_definition.test.dev_center_id
}
`, DevCenterDevBoxDefinitionTestResource{}.basic(data))
}
