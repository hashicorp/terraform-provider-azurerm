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
	parentResourceName := "azurerm_storage_account.test"
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
	parentResourceName := "azurerm_storage_account.test"
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
data "azurerm_client_config" "current" {}

provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "testrg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = azurerm_resource_group.testrg.location
  resource_group_name = azurerm_resource_group.testrg.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  soft_delete_enabled      = true
  purge_protection_enabled = true

  tags = {
    environment = "testing"
  }
}

resource "azurerm_key_vault_key" "test" {
  name                       = "key%d"
  key_vault_id               = azurerm_key_vault.test.id
  key_vault_access_policy_id = azurerm_key_vault_access_policy.storage.id
  key_type                   = "RSA"
  key_size                   = 2048
  key_opts                   = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_key_vault_access_policy" "storage" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_storage_account.test.identity.0.principal_id

  key_permissions    = ["get", "create", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.service_principal_application_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.testrg.name
  location                 = azurerm_resource_group.testrg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

	identity {
    type = "SystemAssigned"
	}
	
  tags = {
    environment = "testing"
  }
}

resource "azurerm_storage_account_customer_managed_key" "custom" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_id       = azurerm_key_vault.test.id
  key_name           = azurerm_key_vault_key.test.name
  key_version        = azurerm_key_vault_key.test.version
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomStringOfLength(4))
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

  location                 = azurerm_resource_group.testrg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(4))
}
