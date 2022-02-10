package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BackupProtectedFileShareDataSource struct{}

func TestAccDataSourceBackupProtectedFileShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_file_share", "test")
	r := BackupProtectedFileShareDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source_file_share_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("source_storage_account_id").Exists(),
				check.That(data.ResourceName).Key("backup_policy_id").Exists(),
			),
		},
	})
}

func (BackupProtectedFileShareDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_protected_file_share" "test" {
  source_file_share_name    = azurerm_storage_share.test.name
  resource_group_name       = azurerm_resource_group.test.name
  recovery_vault_name       = azurerm_recovery_services_vault.test.name
  source_storage_account_id = azurerm_backup_container_storage_account.test.storage_account_id
}
`, BackupProtectedFileShareResource{}.basic(data))
}
