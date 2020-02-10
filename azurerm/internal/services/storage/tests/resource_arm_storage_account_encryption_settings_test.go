package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMStorageAccountEncryptionSettings_basic(t *testing.T) {
	parentResourceName := "azurerm_storage_account.testsa"
	resourceName := "azurerm_storage_account_encryption_settings.custom"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccountEncryptionSettings_basic(ri, rs, location)
	postConfig := testAccAzureRMStorageAccountEncryptionSettings_basicDelete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroyed,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountEncryptionSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_blob_encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_file_encryption", "true"),
				),
			},
			{
				// Delete the encryption settings resource and verify it is gone
				// Whilst making sure the encryption settings on the storage account
				// have been reverted to their default state
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExistsWithDefaultSettings(parentResourceName),
					testCheckAzureRMStorageAccountEncryptionSettingsDestroyed(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccountEncryptionSettings_blobEncryptionDisable(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_STORAGE_ENCRYPTION_DISABLE")
	if !exists {
		t.Skip("`TF_ACC_STORAGE_ENCRYPTION_DISABLE` isn't specified - skipping since disabling encryption is generally disabled")
	}

	resourceName := "azurerm_storage_account_encryption_settings.custom"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccAzureRMStorageAccountEncryptionSettings_blobEncryptionDisabled(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountEncryptionSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_blob_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccountEncryptionSettings_fileEncryptionDisable(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_STORAGE_ENCRYPTION_DISABLE")
	if !exists {
		t.Skip("`TF_ACC_STORAGE_ENCRYPTION_DISABLE` isn't specified - skipping since disabling encryption is generally disabled")
	}

	resourceName := "azurerm_storage_account_encryption_settings.custom"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccAzureRMStorageAccountEncryptionSettings_fileEncryptionDisabled(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountEncryptionSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_file_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccountEncryptionSettings_disappears(t *testing.T) {
	parentResourceName := "azurerm_storage_account.testsa"
	resourceName := "azurerm_storage_account_encryption_settings.custom"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccountEncryptionSettings_basic(ri, rs, testLocation())
	postConfig := testAccAzureRMStorageAccountEncryptionSettings_basicDelete(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroyed,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountEncryptionSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_blob_encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_file_encryption", "true"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExistsWithDefaultSettings(parentResourceName),
					testCheckAzureRMStorageAccountEncryptionSettingsDestroyed(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMStorageAccountExistsWithDefaultSettings(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

		resp, err := conn.GetProperties(ctx, resourceGroup, storageAccount)
		if err != nil {
			return fmt.Errorf("Bad: Get on storageServiceClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: StorageAccount %q (resource group: %q) does not exist", storageAccount, resourceGroup)
		}

		if props := resp.AccountProperties; props != nil {
			if encryption := props.Encryption; encryption != nil {
				if services := encryption.Services; services != nil {
					if !*services.Blob.Enabled {
						return fmt.Errorf("enable_blob_encryption not set to default: %s", resourceName)
					}
					if !*services.File.Enabled {
						return fmt.Errorf("enable_file_encryption not set to default: %s", resourceName)
					}
				}

				if encryption.KeySource != storage.MicrosoftStorage {
					return fmt.Errorf("%s keySource not set to default(storage.MicrosoftStorage): %s", resourceName, encryption.KeySource)
				}
			} else {
				return fmt.Errorf("storage account encryption properties not found: %s", resourceName)
			}
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountEncryptionSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if storageAccountId := rs.Primary.Attributes["storage_account_id"]; storageAccountId == "" {
			return fmt.Errorf("Unable to read storageAccountId: %s", resourceName)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDestroyed(s *terraform.State) error {
	ctx := testAccProvider.Meta().(*ArmClient).StopContext
	conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetProperties(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Storage Account still exists:\n%#v", resp.AccountProperties)
		}
	}

	return nil
}

func testCheckAzureRMStorageAccountEncryptionSettingsDestroyed(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return nil
		}

		return fmt.Errorf("Found: %s", resourceName)
	}
}

func testAccAzureRMStorageAccountEncryptionSettings_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "production"
  }
}

resource "azurerm_storage_account_encryption_settings" "custom" {
  storage_account_id     = "${azurerm_storage_account.testsa.id}"
  enable_blob_encryption = true
  enable_file_encryption = true
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccountEncryptionSettings_basicDelete(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccountEncryptionSettings_fileEncryptionDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "production"
  }
}

resource "azurerm_storage_account_encryption_settings" "custom" {
  storage_account_id     = "${azurerm_storage_account.testsa.id}
  enable_file_encryption = false"
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccountEncryptionSettings_blobEncryptionDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "production"
  }
}

resource "azurerm_storage_account_encryption_settings" "custom" {
  storage_account_id     = "${azurerm_storage_account.testsa.id}
  enable_blob_encryption = false"
}
`, rInt, location, rString)
}
