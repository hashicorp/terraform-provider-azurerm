// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SiteRecoveryFabricDataSource struct{}

func TestAccDataSourceSiteRecoveryFabric_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_fabric", "test")
	r := SiteRecoveryFabricDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
			),
		},
	})
}

func (SiteRecoveryFabricDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_site_recovery_fabric" "test" {
  name                = azurerm_site_recovery_fabric.test.name
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
}
`, SiteRecoveryFabricResource{}.basic(data))
}
