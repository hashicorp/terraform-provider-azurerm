package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data.ResourceName, data.RandomInteger),
			},
			data.RequiresImportErrorStep(testAccAzureRMRecoveryServicesProtectionPolicyVm_requiresImport),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data.ResourceName, data.RandomInteger),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data.ResourceName, data.RandomInteger),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateDailyToWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data.ResourceName, data.RandomInteger),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeeklyToDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data.ResourceName, data.RandomInteger),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeeklyToPartial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_protection_policy_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(data.ResourceName, data.RandomInteger),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(data),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(data.ResourceName, data.RandomInteger),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_recovery_services_protection_policy_vm" {
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

func testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName string) resource.TestCheckFunc {
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

func testAccAzureRMRecoveryServicesProtectionPolicyVm_base(data acceptance.TestData) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesProtectionPolicyVm_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

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

func testAccAzureRMRecoveryServicesProtectionPolicyVm_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(data)
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "import" {
  name                = "${azurerm_recovery_services_protection_policy_vm.test.name}"
  resource_group_name = "${azurerm_recovery_services_protection_policy_vm.test.resource_group_name}"
  recovery_vault_name = "${azurerm_recovery_services_protection_policy_vm.test.recovery_vault_name}"

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

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "10"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesProtectionPolicyVm_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

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

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Weekly"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "2"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesProtectionPolicyVm_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

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

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "10"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.count", "7"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.count", "77"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.months.#", "2"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesProtectionPolicyVm_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

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

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Weekly"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.count", "7"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.count", "77"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.months.#", "2"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesProtectionPolicyVm_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday", "Saturday"]
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

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Weekly"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.count", "7"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.count", "77"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weekdays.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weeks.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.months.#", "1"),
	)
}
