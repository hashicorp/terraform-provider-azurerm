// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type NetworkAnchorDataSource struct{}

func TestAccNetworkAnchorDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, fmt.Sprintf("data.%[1]s", oracle.NetworkAnchorDataSource{}.ResourceType()), "test")
	r := NetworkAnchorDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("resource_anchor_id").Exists(),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (d NetworkAnchorDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_network_anchor" "test" {
  name                = azurerm_oracle_network_anchor.test.name
  resource_group_name = azurerm_oracle_network_anchor.test.resource_group_name
}
`, NetworkAnchorResource{}.basic(data))
}
