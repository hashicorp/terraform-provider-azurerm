package azurerm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: These tests fail because enabling backup on file shares with no content
func TestAccAzureRMBackupProtectedFileShare_basic(t *testing.T) {
	resourceName := "azurerm_backup_protected_file_share.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectedFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectedFileShare_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBackupProtectedFileShareExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{ //vault cannot be deleted unless we unregister all backups
				Config: testAccAzureRMBackupProtectedFileShare_base(ri, acceptance.Location()),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccAzureRMBackupProtectedFileShare_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_backup_protected_file_share.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectedFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectedFileShare_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBackupProtectedFileShareExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
				),
			},
			{
				Config:      testAccAzureRMBackupProtectedFileShare_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_backup_protected_file_share"),
			},
			{ //vault cannot be deleted unless we unregister all backups
				Config: testAccAzureRMBackupProtectedFileShare_base(ri, acceptance.Location()),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccAzureRMBackupProtectedFileShare_updateBackupPolicyId(t *testing.T) {
	protectedFileShareResourceName := "azurerm_backup_protected_file_share.test"
	fBackupPolicyResourceName := "azurerm_backup_policy_file_share.test1"
	sBackupPolicyResourceName := "azurerm_backup_policy_file_share.test2"

	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectedFileShareDestroy,
		Steps: []resource.TestStep{
			{ // Create resources and link first backup policy id
				Config: testAccAzureRMBackupProtectedFileShare_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(protectedFileShareResourceName, "backup_policy_id", fBackupPolicyResourceName, "id"),
				),
			},
			{ // Modify backup policy id to the second one
				// Set Destroy false to prevent error from cleaning up dangling resource
				Config: testAccAzureRMBackupProtectedFileShare_updatePolicy(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(protectedFileShareResourceName, "backup_policy_id", sBackupPolicyResourceName, "id"),
				),
			},
			{ // Remove protected items first before the associated policies are deleted
				Config: testAccAzureRMBackupProtectedFileShare_base(ri, acceptance.Location()),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testCheckAzureRMBackupProtectedFileShareDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_backup_protected_file_share" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		storageID := rs.Primary.Attributes["source_storage_account_id"]
		fileShareName := rs.Primary.Attributes["source_file_share_name"]

		parsedStorageID, err := azure.ParseAzureResourceID(storageID)
		if err != nil {
			return fmt.Errorf("[ERROR] Unable to parse source_storage_account_id '%s': %+v", storageID, err)
		}
		accountName, hasName := parsedStorageID.Path["storageAccounts"]
		if !hasName {
			return fmt.Errorf("[ERROR] parsed source_storage_account_id '%s' doesn't contain 'storageAccounts'", storageID)
		}

		protectedItemName := fmt.Sprintf("AzureFileShare;%s", fileShareName)
		containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageID.ResourceGroup, accountName)

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectedItemsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure Backup Protected File Share still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMBackupProtectedFileShareExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Azure Backup Protected File Share: %q", resourceName)
		}

		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		storageID := rs.Primary.Attributes["source_storage_account_id"]
		fileShareName := rs.Primary.Attributes["source_file_share_name"]

		parsedStorageID, err := azure.ParseAzureResourceID(storageID)
		if err != nil {
			return fmt.Errorf("[ERROR] Unable to parse source_storage_account_id '%s': %+v", storageID, err)
		}
		accountName, hasName := parsedStorageID.Path["storageAccounts"]
		if !hasName {
			return fmt.Errorf("[ERROR] parsed source_storage_account_id '%s' doesn't contain 'storageAccounts'", storageID)
		}

		protectedItemName := fmt.Sprintf("AzureFileShare;%s", fileShareName)
		containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageID.ResourceGroup, accountName)

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectedItemsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Azure Backup Protected File Share %q (resource group: %q) was not found: %+v", protectedItemName, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on recoveryServicesVaultsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMBackupProtectedFileShare_base(rInt int, location string) string {
	rstr := strconv.Itoa(rInt)
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]s"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "acctest-ss-%[1]d"
  storage_account_name = "${azurerm_storage_account.test.name}"
  metadata             = {}
  
  lifecycle {
	ignore_changes = [metadata] // Ignore changes Azure Backup makes to the metadata
  }
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-VAULT-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_backup_policy_file_share" "test1" {
  name                = "acctest-PFS-%[1]d"
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
`, rInt, location, rstr[len(rstr)-5:])
}

func testAccAzureRMBackupProtectedFileShare_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_container_storage_account" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  storage_account_id  = "${azurerm_storage_account.test.id}"
}

resource "azurerm_backup_protected_file_share" "test" {
  resource_group_name       = "${azurerm_resource_group.test.name}"
  recovery_vault_name       = "${azurerm_recovery_services_vault.test.name}"
  source_storage_account_id = "${azurerm_backup_container_storage_account.test.storage_account_id}"
  source_file_share_name    = "${azurerm_storage_share.test.name}"
  backup_policy_id          = "${azurerm_backup_policy_file_share.test1.id}"
}
`, testAccAzureRMBackupProtectedFileShare_base(rInt, location))
}

func testAccAzureRMBackupProtectedFileShare_updatePolicy(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test2" {
  name                = "acctest-%d-Secondary"
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

resource "azurerm_backup_container_storage_account" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  storage_account_id  = "${azurerm_storage_account.test.id}"
}

resource "azurerm_backup_protected_file_share" "test" {
  resource_group_name       = "${azurerm_resource_group.test.name}"
  recovery_vault_name       = "${azurerm_recovery_services_vault.test.name}"
  source_storage_account_id = "${azurerm_backup_container_storage_account.test.storage_account_id}"
  source_file_share_name    = "${azurerm_storage_share.test.name}"
  backup_policy_id          = "${azurerm_backup_policy_file_share.test2.id}"
}
`, testAccAzureRMBackupProtectedFileShare_base(rInt, location), rInt)
}

func testAccAzureRMBackupProtectedFileShare_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_file_share" "test_import" {
  resource_group_name       = "${azurerm_resource_group.test.name}"
  recovery_vault_name       = "${azurerm_recovery_services_vault.test.name}"
  source_storage_account_id = "${azurerm_storage_account.test.id}"
  source_file_share_name    = "${azurerm_storage_share.test.name}"
  backup_policy_id          = "${azurerm_backup_policy_file_share.test1.id}"
}
`, testAccAzureRMBackupProtectedFileShare_basic(rInt, location))
}
