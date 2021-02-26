package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BackupProtectionPolicyFileShareResource struct {
}

func TestAccBackupProtectionPolicyFileShare_basicDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check:  checkAccBackupProtectionPolicyFileShare_basicDaily(data.ResourceName, data.RandomInteger),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check:  checkAccBackupProtectionPolicyFileShare_basicDaily(data.ResourceName, data.RandomInteger),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBackupProtectionPolicyFileShare_updateDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check:  checkAccBackupProtectionPolicyFileShare_basicDaily(data.ResourceName, data.RandomInteger),
		},
		data.ImportStep(),
		{
			Config: r.updateDaily(data),
			Check:  checkAccBackupProtectionPolicyFileShare_updateDaily(data.ResourceName, data.RandomInteger),
		},
		data.ImportStep(),
	})
}

func (t BackupProtectionPolicyFileShareResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	policyName := id.Path["backupPolicies"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup

	resp, err := clients.RecoveryServices.ProtectionPoliciesClient.Get(ctx, vaultName, resourceGroup, policyName)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service Protection Policy (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (BackupProtectionPolicyFileShareResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-RSV-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r BackupProtectionPolicyFileShareResource) basicDaily(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
  name                = "acctest-PFS-%d"
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
`, r.base(data), data.RandomInteger)
}

func (r BackupProtectionPolicyFileShareResource) updateDaily(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
  name                = "acctest-PFS-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:30"
  }

  retention_daily {
    count = 180
  }
}
`, r.base(data), data.RandomInteger)
}

func (BackupProtectionPolicyFileShareResource) requiresImport(data acceptance.TestData) string {
	template := BackupProtectionPolicyFileShareResource{}.basicDaily(data)
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "import" {
  name                = azurerm_backup_policy_file_share.test.name
  resource_group_name = azurerm_backup_policy_file_share.test.resource_group_name
  recovery_vault_name = azurerm_backup_policy_file_share.test.recovery_vault_name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, template)
}

func checkAccBackupProtectionPolicyFileShare_basicDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-PFS-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-backup-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-RSV-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "10"),
	)
}

func checkAccBackupProtectionPolicyFileShare_updateDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-PFS-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-backup-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-RSV-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:30"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "180"),
	)
}
