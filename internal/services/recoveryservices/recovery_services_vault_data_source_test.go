// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RecoveryServicesVaultDataSource struct{}

func TestAccDataSourceAzureRMRecoveryServicesVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("identity.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceAzureRMRecoveryServicesVault_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsNotEmpty(),
			),
		},
	})
}

func (RecoveryServicesVaultDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_recovery_services_vault" "test" {
  name                = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, RecoveryServicesVaultResource{}.basic(data))
}

func (RecoveryServicesVaultDataSource) basicWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_recovery_services_vault" "test" {
  name                = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, RecoveryServicesVaultResource{}.basicWithIdentity(data))
}
