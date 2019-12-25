package tests

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDiskEncryptionSet_basic(t *testing.T) {
	resourceName := "azurerm_disk_encryption_set.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	vaultName := fmt.Sprintf("vault%d", ri)
	keyName := fmt.Sprintf("key-%s", rs)
	desName := fmt.Sprintf("acctestdes-%d", ri)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			// This test step is temporary due to freezing of functions in keyVault.
			// After applying soft-delete and purge-protection in keyVault, this extra step can be removed.
			{
				Config:  testAccPrepareKeyvaultAndKey(resourceGroup, location, vaultName, keyName),
				Destroy: false,
				Check:   resource.ComposeTestCheckFunc(),
			},
			{
				PreConfig: func() { enableSoftDeleteAndPurgeProtectionForKeyvault(resourceGroup, vaultName) },
				Config:    testAccAzureRMDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "active_key.0.source_vault_id"),
					resource.TestCheckResourceAttrSet(resourceName, "active_key.0.key_url"),
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

func testCheckAzureRMDiskEncryptionSetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Disk Encryption Set not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DiskEncryptionSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Disk Encryption Set %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on Compute.DiskEncryptionSetsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDiskEncryptionSetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_disk_encryption_set" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.DiskEncryptionSetsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func enableSoftDeleteAndPurgeProtectionForKeyvault(resourceGroup, vaultName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		armClient := acceptance.AzureProvider.Meta().(*clients.Client)
		client := armClient.KeyVault.VaultsClient
		ctx := armClient.StopContext
		vaultPatch := keyvault.VaultPatchParameters{
			Properties: &keyvault.VaultPatchProperties{
				EnableSoftDelete:      utils.Bool(true),
				EnablePurgeProtection: utils.Bool(true),
			},
		}
		log.Printf("[LOG] Updating")
		_, err := client.Update(ctx, resourceGroup, vaultName, vaultPatch)
		if err != nil {
			return fmt.Errorf("Bad: error when updating Keyvault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
		}
		return nil
	}
}

func testAccPrepareKeyvaultAndKey(resourceGroup, location, vaultName, keyName string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.service_principal_object_id

    key_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, resourceGroup, location, vaultName, keyName)
}

func testAccAzureRMDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.service_principal_object_id

    key_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  active_key {
    source_vault_id = azurerm_key_vault.test.id
    key_url         = azurerm_key_vault_key.test.id
  }
}
`, resourceGroup, location, vaultName, keyName, desName)
}
