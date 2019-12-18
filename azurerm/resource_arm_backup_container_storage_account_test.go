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
)

func TestAccAzureRMBackupProtectionContainerStorageAccount_basic(t *testing.T) {
	resourceGroupName := "azurerm_resource_group.testrg"
	vaultName := "azurerm_recovery_services_vault.testvlt"
	storageAccountName := "azurerm_storage_account.testsa"
	resourceName := "azurerm_backup_container_storage_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectionContainerStorageAccount(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBackupProtectionContainerStorageAccount(resourceGroupName, vaultName, storageAccountName, resourceName),
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

func testAccAzureRMBackupProtectionContainerStorageAccount(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-%d"
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

func testCheckAzureRMBackupProtectionContainerStorageAccount(resourceGroupStateName, vaultStateName, storageAccountName, resourceStateName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		resourceGroupState, ok := s.RootModule().Resources[resourceGroupStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceGroupStateName)
		}
		vaultState, ok := s.RootModule().Resources[vaultStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", vaultStateName)
		}
		storageState, ok := s.RootModule().Resources[storageAccountName]
		if !ok {
			return fmt.Errorf("Not found: %s", storageAccountName)
		}
		protectionContainerState, ok := s.RootModule().Resources[resourceStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceStateName)
		}

		resourceGroupName := resourceGroupState.Primary.Attributes["name"]
		vaultName := vaultState.Primary.Attributes["name"]
		storageAccountID := storageState.Primary.Attributes["id"]
		resourceStorageID := protectionContainerState.Primary.Attributes["storage_account_id"]

		if storageAccountID != resourceStorageID {
			return fmt.Errorf("Bad: Container resource's storage_account_id %q does not match storage account resource's ID %q", storageAccountID, resourceStorageID)
		}

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
		client := testAccProvider.Meta().(*ArmClient).RecoveryServices.BackupProtectionContainersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
