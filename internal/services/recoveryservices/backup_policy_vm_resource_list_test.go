// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccBackupProtectionPolicyVM_listByRecoveryVaultID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}
	listResourceAddress := "azurerm_backup_policy_vm.list"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQueryByRecoveryVaultID(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// a Recovery Services Vault also contains built-in default VM policies, so we expect at least the 3 we created
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
		},
	})
}

// provision multiple Backup Policy VM resources for testing
func (r BackupProtectionPolicyVMResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_backup_policy_vm" "test" {
  count = 3

  name                = "acctest-%[2]d-${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, r.template(data), data.RandomInteger)
}

// define the list query for testing by recovery vault ID
func (r BackupProtectionPolicyVMResource) basicQueryByRecoveryVaultID() string {
	return `
list "azurerm_backup_policy_vm" "list" {
  provider = azurerm
  config {
    recovery_vault_id = azurerm_recovery_services_vault.test.id
  }
}
`
}
