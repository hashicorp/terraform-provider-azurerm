package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMBackupProtectionContainerStorageAccount_basic(t *testing.T) {
	resourceName := "azurerm_backup_container_storage_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectionContainerStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionContainerStorageAccount_basic(ri, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBackupProtectionContainerStorageAccountExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMBackupProtectionContainerStorageAccount_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "testvlt" {
  name                = "acctest-vault-%d"
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"
  sku                 = "Standard"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_backup_container_storage_account" "test" {
  resource_group_name  = "${azurerm_resource_group.testrg.name}"
  recovery_vault_name  = "${azurerm_recovery_services_vault.testvlt.name}"
  storage_account_id   = "${azurerm_storage_account.testsa.id}"
}
`, rInt, location, rInt, rString)
}

func testCheckAzureRMBackupProtectionContainerStorageAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		state, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroupName := state.Primary.Attributes["resource_group_name"]
		vaultName := state.Primary.Attributes["recovery_vault_name"]
		storageAccountID := state.Primary.Attributes["storage_account_id"]

		parsedStorageAccountID, err := azure.ParseAzureResourceID(storageAccountID)
		if err != nil {
			return fmt.Errorf("Bad: Unable to parse storage_account_id '%s': %+v", storageAccountID, err)
		}
		accountName, hasName := parsedStorageAccountID.Path["storageAccounts"]
		if !hasName {
			return fmt.Errorf("Bad: Parsed storage_account_id '%s' doesn't contain 'storageAccounts'", storageAccountID)
		}

		containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroup, accountName)

		// Ensure container exists in API
		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.BackupProtectionContainersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, vaultName, resourceGroupName, "Azure", containerName)
		if err != nil {
			return fmt.Errorf("Bad: Get on protection container: %+v", err)
		}

		if resp.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: container: %q does not exist", containerName)
		}

		return nil
	}
}

func testCheckAzureRMBackupProtectionContainerStorageAccountDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_backup_container_storage_account" {
			continue
		}

		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		storageAccountID := rs.Primary.Attributes["storage_account_id"]

		parsedStorageAccountID, err := azure.ParseAzureResourceID(storageAccountID)
		if err != nil {
			return fmt.Errorf("Bad: Unable to parse storage_account_id '%s': %+v", storageAccountID, err)
		}
		accountName, hasName := parsedStorageAccountID.Path["storageAccounts"]
		if !hasName {
			return fmt.Errorf("Bad: Parsed storage_account_id '%s' doesn't contain 'storageAccounts'", storageAccountID)
		}

		containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroup, accountName)

		// Ensure container exists in API
		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.BackupProtectionContainersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, vaultName, resourceGroupName, "Azure", containerName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Backup Container Storage Account still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}
