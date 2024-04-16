// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppPoolDataSource struct{}

func TestAccDataSourceNetAppPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_pool", "test")
	r := NetAppPoolDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("account_name").Exists(),
				check.That(data.ResourceName).Key("service_level").Exists(),
				check.That(data.ResourceName).Key("size_in_tb").Exists(),
				check.That(data.ResourceName).Key("encryption_type").Exists(),
			),
		},
	})
}

func (NetAppPoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_pool" "test" {
  resource_group_name = azurerm_netapp_pool.test.resource_group_name
  account_name        = azurerm_netapp_pool.test.account_name
  name                = azurerm_netapp_pool.test.name
}
`, NetAppPoolResource{}.basic(data))
}
