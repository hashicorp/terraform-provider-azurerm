// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DevCenterNetworkConnectionDataSource struct{}

func TestAccDevCenterNetworkConnectionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_center_network_connection", "test")
	r := DevCenterNetworkConnectionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("domain_join_type").Exists(),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
	})
}

func (d DevCenterNetworkConnectionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dev_center_network_connection" "test" {
  name                = azurerm_dev_center_network_connection.test.name
  resource_group_name = azurerm_dev_center_network_connection.test.resource_group_name
}
`, DevCenterNetworkConnectionTestResource{}.basic(data))
}
