package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BackupProtectionPolicyVMResource struct {
}

func TestAccBackupProtectionPolicyVM_basicDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_withInstantRestoreRetentionRangeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicDailyWithInstantRestoreRetentionRange(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBackupProtectionPolicyVM_basicWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicWeekly(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeWeekly(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Weekly"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("backup.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("42"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("77"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("42"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("77"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicWeekly(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeWeekly(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Weekly"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("backup.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("42"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("77"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateDailyToWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWeekly(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeeklyToDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicWeekly(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeeklyToPartial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeWeekly(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeWeeklyPartial(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BackupProtectionPolicyVMResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func (BackupProtectionPolicyVMResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) basicDaily(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
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

func (r BackupProtectionPolicyVMResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "import" {
  name                = azurerm_backup_policy_vm.test.name
  resource_group_name = azurerm_backup_policy_vm.test.resource_group_name
  recovery_vault_name = azurerm_backup_policy_vm.test.recovery_vault_name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, r.basicDaily(data))
}

func (r BackupProtectionPolicyVMResource) basicWeekly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday"]
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) completeDaily(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) completeWeekly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) completeWeeklyPartial(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday"]
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday"]
    weeks    = ["Last"]
    months   = ["January"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) basicDailyWithInstantRestoreRetentionRange(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                           = "acctest-BPVM-%d"
  resource_group_name            = azurerm_resource_group.test.name
  recovery_vault_name            = azurerm_recovery_services_vault.test.name
  instant_restore_retention_days = 5
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
