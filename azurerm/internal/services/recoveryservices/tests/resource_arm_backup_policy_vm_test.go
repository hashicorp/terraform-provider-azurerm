package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBackupProtectionPolicyVM_basicDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.frequency", "Daily"),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.time", "23:00"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_daily.0.count", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_withInstantRestoreRetentionRangeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDailyWithInstantRestoreRetentionRange(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMBackupProtectionPolicyVM_requiresImport),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_basicWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_completeDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_completeDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_completeWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_completeWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.frequency", "Weekly"),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.time", "23:00"),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_weekly.0.count", "42"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_weekly.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.count", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.weeks.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.count", "77"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.weeks.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.months.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_updateDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_completeDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.frequency", "Daily"),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.time", "23:00"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_daily.0.count", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_weekly.0.count", "42"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_weekly.0.weekdays.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.count", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.weekdays.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.weeks.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.count", "77"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.weekdays.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.weeks.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.months.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_updateWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_completeWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.frequency", "Weekly"),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.time", "23:00"),
					resource.TestCheckResourceAttr(data.ResourceName, "backup.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_weekly.0.count", "42"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_weekly.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.count", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_monthly.0.weeks.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.count", "77"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.weekdays.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.weeks.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_yearly.0.months.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_updateDailyToWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_updateWeeklyToDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_basicDaily(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyVM_updateWeeklyToPartial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_completeWeekly(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBackupProtectionPolicyVM_completeWeeklyPartial(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMBackupProtectionPolicyVmExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMBackupProtectionPolicyVmDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_backup_policy_vm" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		policyName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Recovery Services Vault Policy still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMBackupProtectionPolicyVmExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectionPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		policyName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Recovery Services Vault %q Policy: %q", vaultName, policyName)
		}

		resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Recovery Services Vault Policy %q (resource group: %q) was not found: %+v", policyName, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on recoveryServicesVaultsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMBackupProtectionPolicyVM_template(data acceptance.TestData) string {
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

func testAccAzureRMBackupProtectionPolicyVM_basicDaily(data acceptance.TestData) string {
	template := testAccAzureRMBackupProtectionPolicyVM_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMBackupProtectionPolicyVM_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMBackupProtectionPolicyVM_basicDaily(data)
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
`, template)
}

func testAccAzureRMBackupProtectionPolicyVM_basicWeekly(data acceptance.TestData) string {
	template := testAccAzureRMBackupProtectionPolicyVM_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMBackupProtectionPolicyVM_completeDaily(data acceptance.TestData) string {
	template := testAccAzureRMBackupProtectionPolicyVM_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMBackupProtectionPolicyVM_completeWeekly(data acceptance.TestData) string {
	template := testAccAzureRMBackupProtectionPolicyVM_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMBackupProtectionPolicyVM_completeWeeklyPartial(data acceptance.TestData) string {
	template := testAccAzureRMBackupProtectionPolicyVM_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMBackupProtectionPolicyVM_basicDailyWithInstantRestoreRetentionRange(data acceptance.TestData) string {
	template := testAccAzureRMBackupProtectionPolicyVM_template(data)
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
`, template, data.RandomInteger)
}
