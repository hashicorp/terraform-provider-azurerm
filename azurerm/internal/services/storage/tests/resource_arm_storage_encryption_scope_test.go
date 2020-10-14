package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageEncryptionScope_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageEncryptionScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageEncryptionScope_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageEncryptionScope_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageEncryptionScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageEncryptionScope_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageEncryptionScope_requiresImport),
		},
	})
}

func TestAccAzureRMStorageEncryptionScope_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageEncryptionScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageEncryptionScope_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageEncryptionScope_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageEncryptionScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageEncryptionScope_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageEncryptionScope_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageEncryptionScope_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageEncryptionScopeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("storage Sync Service not found: %s", resourceName)
		}

		id, err := parsers.StorageEncryptionScopeID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.SyncServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Storage Sync Service (Storage Sync Service Name %q / Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("bad: Get on StorageEncryptionScopesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageEncryptionScopeDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.SyncServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_encryption_scope" {
			continue
		}

		id, err := parsers.StorageEncryptionScopeID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on StorageEncryptionScopesClient: %+v", err)
			}
		}

		return nil
	}
	return nil
}

func testAccAzureRMStorageEncryptionScope_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageEncryptionScope_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageEncryptionScope_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "enscope%s"
  storage_account_id = azurerm_storage_account.test.id
}
`, template, data.RandomString)
}

func testAccAzureRMStorageEncryptionScope_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageEncryptionScope_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_encryption_scope" "import" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, template)
}

func testAccAzureRMStorageEncryptionScope_complete(data acceptance.TestData) string {
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
  name     = "acctestRG-storage-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%[3]s"
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

  key_permissions = ["get", "unwrapkey", "wrapkey"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
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
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_encryption_scope" "test" {
  name               = "enscope%[3]s"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.KeyVault"
  key_vault_key_id   = azurerm_key_vault_key.first.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
