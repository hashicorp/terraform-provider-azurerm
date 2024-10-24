// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppBackupVaultDataSource struct{}

func TestAccNetAppBackupVaultDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_backup_vault", "test")
	d := NetAppBackupVaultDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (d NetAppBackupVaultDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_backup_vault" "test" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  name                = azurerm_netapp_backup_vault.test.name
}
`, NetAppBackupVaultResource{}.basic(data))
}
