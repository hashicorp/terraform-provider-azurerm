package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAzureRMStorageAccountCustomerManagedKey_basic(t *testing.T) {
	parentResourceName := "azurerm_storage_account.testsa"
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "custom")
	preConfig := testAccAzureRMStorageAccountCustomerManagedKey_basic(data)
	postConfig := testAccAzureRMStorageAccountCustomerManagedKey_basicDelete(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroyed,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			{
				// Delete the encryption settings resource and verify it is gone
				// Whilst making sure the encryption settings on the storage account
				// have been reverted to their default state
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExistsWithDefaultSettings(parentResourceName),
					testCheckAzureRMStorageAccountCustomerManagedKeyDestroyed(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccountCustomerManagedKey_disappears(t *testing.T) {
	parentResourceName := "azurerm_storage_account.testsa"
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "custom")
	preConfig := testAccAzureRMStorageAccountCustomerManagedKey_basic(data)
	postConfig := testAccAzureRMStorageAccountCustomerManagedKey_basicDelete(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroyed,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExistsWithDefaultSettings(parentResourceName),
					testCheckAzureRMStorageAccountCustomerManagedKeyDestroyed(data.ResourceName),
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
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.AccountsClient

		resp, err := conn.GetProperties(ctx, resourceGroup, storageAccount, "")
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

func testCheckAzureRMStorageAccountCustomerManagedKeyExists(resourceName string) resource.TestCheckFunc {
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
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.AccountsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetProperties(ctx, resourceGroup, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Storage Account still exists:\n%#v", resp.AccountProperties)
		}
	}

	return nil
}

func testCheckAzureRMStorageAccountCustomerManagedKeyDestroyed(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return nil
		}

		return fmt.Errorf("Found: %s", resourceName)
	}
}

func testAccAzureRMStorageAccountCustomerManagedKey_basic(data acceptance.TestData) string {
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

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account_customer_managed_key" "custom" {
  storage_account_id     = "${azurerm_storage_account.testsa.id}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(4))
}

func testAccAzureRMStorageAccountCustomerManagedKey_basicDelete(data acceptance.TestData) string {
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

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(4))
}
