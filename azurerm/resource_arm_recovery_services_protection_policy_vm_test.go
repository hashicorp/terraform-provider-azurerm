package azurerm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				Config:      testAccAzureRMRecoveryServicesProtectionPolicyVm_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_recovery_services_protection_policy_vm"),
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateDailyToWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeeklyToDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeeklyToPartial(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).recoveryServices.ProtectionPoliciesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
		client := testAccProvider.Meta().(*ArmClient).recoveryServices.ProtectionPoliciesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt int, location string) string {
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
}
`, rInt, location, strconv.Itoa(rInt)[12:17])
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(rInt int, location string) string {
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
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_requiresImport(rInt int, location string) string {
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
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(rInt, location))
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

func testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(rInt int, location string) string {
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
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
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

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(rInt int, location string) string {
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
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
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

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(rInt int, location string) string {
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
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
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

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(rInt int, location string) string {
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
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
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
