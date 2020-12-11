package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
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
				Config: testAccAzureRMDiskEncryptionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDiskEncryptionSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskEncryptionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDiskEncryptionSet_requiresImport),
		},
	})
}

func TestAccAzureRMDiskEncryptionSet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskEncryptionSet_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDiskEncryptionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskEncryptionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDiskEncryptionSet_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDiskEncryptionSet_keyRotate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskEncryptionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			// we have to first grant the permission for DiskEncryptionSet to access the KeyVault
			{
				Config: testAccAzureRMDiskEncryptionSet_grantAccessToKeyVault(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskEncryptionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			// after the access is granted, we can rotate the key in DiskEncryptionSet
			{
				Config: testAccAzureRMDiskEncryptionSet_keyRotate(data),
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

		id, err := parse.DiskEncryptionSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DiskEncryptionSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Disk Encryption Set %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
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

		id, err := parse.DiskEncryptionSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.DiskEncryptionSetsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDiskEncryptionSet_dependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv-%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  soft_delete_enabled         = true
  purge_protection_enabled    = true
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

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

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMDiskEncryptionSet_basic(data acceptance.TestData) string {
	template := testAccAzureRMDiskEncryptionSet_dependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDiskEncryptionSet_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDiskEncryptionSet_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "import" {
  name                = azurerm_disk_encryption_set.test.name
  resource_group_name = azurerm_disk_encryption_set.test.resource_group_name
  location            = azurerm_disk_encryption_set.test.location
  key_vault_key_id    = azurerm_disk_encryption_set.test.key_vault_key_id

  identity {
    type = "SystemAssigned"
  }
}
`, template)
}

func testAccAzureRMDiskEncryptionSet_complete(data acceptance.TestData) string {
	template := testAccAzureRMDiskEncryptionSet_dependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Hello = "woRld"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDiskEncryptionSet_grantAccessToKeyVault(data acceptance.TestData) string {
	template := testAccAzureRMDiskEncryptionSet_dependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "wrapkey",
    "unwrapkey",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDiskEncryptionSet_keyRotate(data acceptance.TestData) string {
	template := testAccAzureRMDiskEncryptionSet_dependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "new" {
  name         = "newKey"
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

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "wrapkey",
    "unwrapkey",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.new.id

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}
