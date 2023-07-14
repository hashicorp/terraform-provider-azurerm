// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BackupProtectionPolicyFileShareDataSource struct{}

func TestAccDataSourceBackupPolicyFileShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
	})
}

func (BackupProtectionPolicyFileShareDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_policy_file_share" "test" {
  name                = azurerm_backup_policy_file_share.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyFileShareResource{}.basicDaily(data))
}
