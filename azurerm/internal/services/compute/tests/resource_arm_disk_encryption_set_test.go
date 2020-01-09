package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDiskEncryptionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			{
				// This test step is temporary due to freezing of functions in keyVault.
				// TODO: After applying soft-delete and purge-protection in keyVault, this extra step can be removed.
				Config: testAccAzureRMDiskEncryptionSet_dependencies(data),
				Check: resource.ComposeTestCheckFunc(
					enableSoftDeleteAndPurgeProtectionForKeyVault("azurerm_key_vault.test"),
				),
			},
			{
				Config: testAccAzureRMDiskEncryptionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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

func enableSoftDeleteAndPurgeProtectionForKeyVault(resourceName string) resource.TestCheckFunc {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		vaultName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		vaultPatch := keyvault.VaultPatchParameters{
			Properties: &keyvault.VaultPatchProperties{
				EnableSoftDelete:      utils.Bool(true),
				EnablePurgeProtection: utils.Bool(true),
			},
		}
		log.Printf("[DEBUG] Enabling Soft Delete & Purge Protection on Key Vault %q (Resource Group %q)..", vaultName, resourceGroup)
		_, err := client.Update(ctx, resourceGroup, vaultName, vaultPatch)
		if err != nil {
			return fmt.Errorf("Bad: error updating KeyVault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
		}
		log.Printf("[DEBUG] Enabled Soft Delete & Purge Protection on Key Vault %q (Resource Group %q)..", vaultName, resourceGroup)

		return nil
	}
}

func testAccAzureRMDiskEncryptionSet_dependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

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
  name         = "examplekey"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDiskEncryptionSet_basic(data acceptance.TestData) string {
	template := testAccAzureRMDiskEncryptionSet_dependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_uri   = azurerm_key_vault_key.test.id
}
`, template, data.RandomInteger)
}
