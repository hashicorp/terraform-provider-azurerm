package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAzureRMStorageAccountCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountCustomerManagedKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Delete the encryption settings resource and verify it is gone
				Config: testAccAzureRMStorageAccountCustomerManagedKey_template(data),
				Check: resource.ComposeTestCheckFunc(
					// Then ensure the encryption settings on the storage account
					// have been reverted to their default state
					testCheckAzureRMStorageAccountExistsWithDefaultSettings("azurerm_storage_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccountCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountCustomerManagedKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageAccountCustomerManagedKey_requiresImport),
		},
	})
}

func TestAccAzureRMStorageAccountCustomerManagedKey_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountCustomerManagedKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccountCustomerManagedKey_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccountCustomerManagedKey_testKeyVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccountCustomerManagedKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccountCustomerManagedKey_autoKeyRotation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountCustomerManagedKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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

				if encryption.KeySource != storage.KeySourceMicrosoftStorage {
					return fmt.Errorf("%s keySource not set to default(storage.KeySourceMicrosoftStorage): %s", resourceName, encryption.KeySource)
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

func testAccAzureRMStorageAccountCustomerManagedKey_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccountCustomerManagedKey_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_id       = azurerm_key_vault.test.id
  key_name           = azurerm_key_vault_key.first.name
  key_version        = azurerm_key_vault_key.first.version
}
`, template)
}

func testAccAzureRMStorageAccountCustomerManagedKey_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccountCustomerManagedKey_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "import" {
  storage_account_id = azurerm_storage_account_customer_managed_key.test.storage_account_id
  key_vault_id       = azurerm_storage_account_customer_managed_key.test.key_vault_id
  key_name           = azurerm_storage_account_customer_managed_key.test.key_name
  key_version        = azurerm_storage_account_customer_managed_key.test.key_version
}
`, template)
}

func testAccAzureRMStorageAccountCustomerManagedKey_updated(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccountCustomerManagedKey_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_id       = azurerm_key_vault.test.id
  key_name           = azurerm_key_vault_key.second.name
  key_version        = azurerm_key_vault_key.second.version
}
`, template)
}

func testAccAzureRMStorageAccountCustomerManagedKey_autoKeyRotation(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccountCustomerManagedKey_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_id       = azurerm_key_vault.test.id
  key_name           = azurerm_key_vault_key.first.name
}
`, template)
}

func testAccAzureRMStorageAccountCustomerManagedKey_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  soft_delete_enabled      = true
  purge_protection_enabled = true
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
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_key" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
