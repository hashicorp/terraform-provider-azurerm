package azurerm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBackupProtectionPolicyFileShare_basicDaily(t *testing.T) {
	resourceName := "azurerm_backup_policy_file_share.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyFileShare_basicDaily(ri, acceptance.Location()),
				Check:  checkAccAzureRMBackupProtectionPolicyFileShare_basicDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyFileShare_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_backup_policy_file_share.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyFileShare_basicDaily(ri, acceptance.Location()),
				Check:  checkAccAzureRMBackupProtectionPolicyFileShare_basicDaily(resourceName, ri),
			},
			{
				Config:      testAccAzureRMBackupProtectionPolicyFileShare_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_backup_policy_file_share"),
			},
		},
	})
}

func TestAccAzureRMBackupProtectionPolicyFileShare_updateDaily(t *testing.T) {
	resourceName := "azurerm_backup_policy_file_share.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionPolicyFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionPolicyFileShare_basicDaily(ri, acceptance.Location()),
				Check:  checkAccAzureRMBackupProtectionPolicyFileShare_basicDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMBackupProtectionPolicyFileShare_updateDaily(ri, acceptance.Location()),
				Check:  checkAccAzureRMBackupProtectionPolicyFileShare_updateDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMBackupProtectionPolicyFileShareDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_backup_policy_file_share" {
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

func testCheckAzureRMBackupProtectionPolicyFileShareExists(resourceName string) resource.TestCheckFunc {
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

func testAccAzureRMBackupProtectionPolicyFileShare_base(rInt int, location string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%[1]d"
  location = "%[2]s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-RSV-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}
`, rInt, location, strconv.Itoa(rInt)[12:17])
}

func testAccAzureRMBackupProtectionPolicyFileShare_basicDaily(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
  name                = "acctest-PFS-%d"
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
`, testAccAzureRMBackupProtectionPolicyFileShare_base(rInt, location), rInt)
}

func testAccAzureRMBackupProtectionPolicyFileShare_updateDaily(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test" {
  name                = "acctest-PFS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

  backup {
    frequency = "Daily"
    time      = "23:30"
  }

  retention_daily {
    count = 180
  }
}
`, testAccAzureRMBackupProtectionPolicyFileShare_base(rInt, location), rInt)
}

func testAccAzureRMBackupProtectionPolicyFileShare_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "import" {
  name                = "${azurerm_backup_policy_file_share.test.name}"
  resource_group_name = "${azurerm_backup_policy_file_share.test.resource_group_name}"
  recovery_vault_name = "${azurerm_backup_policy_file_share.test.recovery_vault_name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, testAccAzureRMBackupProtectionPolicyFileShare_basicDaily(rInt, location))
}

func checkAccAzureRMBackupProtectionPolicyFileShare_basicDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMBackupProtectionPolicyFileShareExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-PFS-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-backup-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-RSV-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "10"),
	)
}

func checkAccAzureRMBackupProtectionPolicyFileShare_updateDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMBackupProtectionPolicyFileShareExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-PFS-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-backup-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-RSV-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:30"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "180"),
	)
}
