package keyvault_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault`
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.1", "set"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.1", "set"),
				),
			},
			{
				Config:      testAccAzureRMKeyVaultAccessPolicy_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_key_vault_access_policy"),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test_with_application_id")
	resourceName2 := "azurerm_key_vault_access_policy.test_no_application_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultAccessPolicy_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.0", "create"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.1", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.1", "delete"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_permissions.0", "create"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_permissions.1", "delete"),

					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "key_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "key_permissions.1", "encrypt"),
					resource.TestCheckResourceAttr(resourceName2, "secret_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "secret_permissions.1", "delete"),
					resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.1", "delete"),
				),
			},
			data.ImportStep(),
			{
				ResourceName:      resourceName2,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.1", "set"),
				),
			},
			{
				Config: testAccAzureRMKeyVaultAccessPolicy_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.0", "list"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.1", "encrypt"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_nonExistentVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccAzureRMKeyVaultAccessPolicy_nonExistentVault(data),
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`Error retrieving Key Vault`),
			},
		},
	})
}

func testCheckAzureRMKeyVaultAccessPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resGroup := id.ResourceGroup
		vaultName := id.Path["vaults"]

		objectId := rs.Primary.Attributes["object_id"]
		applicationId := rs.Primary.Attributes["application_id"]

		resp, err := client.Get(ctx, resGroup, vaultName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault %q (resource group: %q) does not exist", vaultName, resGroup)
			}

			return fmt.Errorf("Bad: Get on keyVaultClient: %+v", err)
		}

		policy := keyvault.FindKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationId)

		if policy == nil {
			return fmt.Errorf("Bad: Key Vault Policy %q (resource group: %q, object_id: %s) does not exist", vaultName, resGroup, objectId)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultAccessPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultAccessPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "import" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault_access_policy.test.tenant_id
  object_id    = azurerm_key_vault_access_policy.test.object_id

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_multiple(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test_with_application_id" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "create",
    "get",
  ]

  secret_permissions = [
    "get",
    "delete",
  ]

  certificate_permissions = [
    "create",
    "delete",
  ]

  application_id = data.azurerm_client_config.current.client_id
  tenant_id      = data.azurerm_client_config.current.tenant_id
  object_id      = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_access_policy" "test_no_application_id" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = [
    "list",
    "delete",
  ]

  certificate_permissions = [
    "list",
    "delete",
  ]

  storage_permissions = [
    "backup",
    "delete",
    "deletesas",
    "get",
    "getsas",
    "list",
    "listsas",
    "purge",
    "recover",
    "regeneratekey",
    "restore",
    "set",
    "setsas",
    "update",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_update(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = []

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "premium"

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMKeyVaultAccessPolicy_nonExistentVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  # Must appear to be URL, but not actually exist - appending a string works
  key_vault_id = "${azurerm_key_vault.test.id}NOPE"

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
