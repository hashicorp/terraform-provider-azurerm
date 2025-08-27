// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppAccountDataSource struct{}

func TestAccDataSourceNetAppAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_account", "test")
	r := NetAppAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccDataSourceNetAppAccount_systemAssignedManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_account", "test")
	r := NetAppAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.systemAssignedManagedIdentityConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
	})
}

func (r NetAppAccountDataSource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_account" "test" {
  resource_group_name = azurerm_netapp_account.test.resource_group_name
  name                = azurerm_netapp_account.test.name
}
`, NetAppAccountResource{}.basicConfig(data))
}

func (r NetAppAccountDataSource) systemAssignedManagedIdentityConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_account" "test" {
  resource_group_name = azurerm_netapp_account.test.resource_group_name
  name                = azurerm_netapp_account.test.name
}
`, NetAppAccountResource{}.systemAssignedManagedIdentity(data))
}
