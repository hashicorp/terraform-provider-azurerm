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

type BackupProtectionPolicyFileShareResource struct {
}

func TestAccBackupProtectionPolicyFileShare_basicDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

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

func TestAccBackupProtectionPolicyFileShare_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

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

func TestAccBackupProtectionPolicyFileShare_WeeklyRetention(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.WeeklyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_WeeklyRetentionImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.WeeklyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBackupProtectionPolicyFileShare_MonthlyRetention(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.MonthlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_MonthlyRetentionImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.MonthlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBackupProtectionPolicyFileShare_YearlyRetention(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.YearlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_YearlyRetentionImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.YearlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBackupProtectionPolicyFileShare_completeDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

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

func TestAccBackupProtectionPolicyFileShare_updateDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

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
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_updateDailyRetentionToWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.WeeklyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_updateWeeklyRetentionToMonthly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.WeeklyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.MonthlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_updateMonthlyRetentionToYearly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.MonthlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.YearlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyFileShare_updateYearlyRetentionToDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.YearlyRetention(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
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

func TestAccBackupProtectionPolicyFileShare_updateDailyToPartial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeDaily(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeDailyPartial(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
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

func (BackupProtectionPolicyFileShareResource) template(data acceptance.TestData) string {
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

func (r BackupProtectionPolicyFileShareResource) basicDaily(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
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

func (r BackupProtectionPolicyFileShareResource) requiresImport(data acceptance.TestData) string {
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
`, r.basicDaily(data))
}

func (r BackupProtectionPolicyFileShareResource) WeeklyRetention(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
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
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
  }

}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyFileShareResource) WeeklyRetentionImport(data acceptance.TestData) string {
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

  retention_weekly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
  }

}
`, r.WeeklyRetention(data))
}

func (r BackupProtectionPolicyFileShareResource) MonthlyRetention(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
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

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyFileShareResource) MonthlyRetentionImport(data acceptance.TestData) string {
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

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

}
`, r.MonthlyRetention(data))
}

func (r BackupProtectionPolicyFileShareResource) YearlyRetention(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
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

  retention_yearly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }

}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyFileShareResource) YearlyRetentionImport(data acceptance.TestData) string {
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

  retention_yearly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }

}
`, r.YearlyRetention(data))
}

func (r BackupProtectionPolicyFileShareResource) completeDaily(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
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
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyFileShareResource) completeDailyPartial(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
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
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First"]
  }

  retention_yearly {
    count    = 7
    weekdays = ["Sunday"]
    weeks    = ["Last"]
    months   = ["January"]
  }
}
`, r.template(data), data.RandomInteger)
}
