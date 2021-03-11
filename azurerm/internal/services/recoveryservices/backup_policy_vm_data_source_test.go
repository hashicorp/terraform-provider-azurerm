package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type BackupProtectionPolicyVMDataSource struct {
}

func TestAccDataSourceBackupPolicyVm_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (BackupProtectionPolicyVMDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_policy_vm" "test" {
  name                = azurerm_backup_policy_vm.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyVMResource{}.basicDaily(data))
}
