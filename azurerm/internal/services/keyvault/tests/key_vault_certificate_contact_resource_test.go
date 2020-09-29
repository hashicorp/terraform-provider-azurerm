package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultCertificateContact_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateContact_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultCertificateContact_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateContact_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMKeyVaultCertificateContact_requiresImport),
		},
	})
}

//
func TestAccAzureRMKeyVaultCertificateContact_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateContact_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultCertificateContact_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateContact_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKeyVaultCertificateContact_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKeyVaultCertificateContact_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultCertificateContact_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateContact_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
					testCheckAzureRMKeyVaultCertificateContactDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificateContact_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateContact_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateContactExists(data.ResourceName),
					testCheckAzureRMKeyVaultDisappears("azurerm_key_vault.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMKeyVaultCertificateContactDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
	vaultClient := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_certificate_contact" {
			continue
		}

		keyVaultId := rs.Primary.Attributes["key_vault_id"]
		vaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			// deleted, this is fine.
			return nil
		}

		ok, err := azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to check if key vault %q for Certificate Contact in Vault at url %q exists: %v", keyVaultId, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Certificate Contact of Key Vault %q was not found in Key Vault at URI %q ", keyVaultId, vaultBaseUrl)
			return nil
		}

		// get the latest version
		resp, err := client.GetCertificateContacts(ctx, vaultBaseUrl)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return fmt.Errorf("bad: Get on keyVault certificate contact: %+v", err)
		}

		return fmt.Errorf("Key Vault Certificate Contact still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultCertificateContactExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
		vaultClient := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		keyVaultId := rs.Primary.Attributes["key_vault_id"]
		vaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to look up key vault url from id %q: %+v", keyVaultId, err)
		}

		ok, err = azure.KeyVaultExists(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to check if key vault %q for Certificate Contact at url %q exists: %v", keyVaultId, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Key Vault %q was not found in Key Vault at URI %q ", keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.GetCertificateContacts(ctx, vaultBaseUrl)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Certificate Contact for key vault %q does not exist", vaultBaseUrl)
			}

			return fmt.Errorf("bad: Get on keyVault certificate Contact: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultCertificateContactDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
		vaultClient := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		keyVaultId := rs.Primary.Attributes["key_vault_id"]
		vaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to look up key vault URI from id %q: %+v", keyVaultId, err)
		}

		ok, err = azure.KeyVaultExists(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to check if key vault %q in Vault at url %q exists: %v", keyVaultId, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Key Vault %q was not found in Key Vault at URI %q ", keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.DeleteCertificateContacts(ctx, vaultBaseUrl)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("bad: Delete on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultCertificateContact_basic(data acceptance.TestData) string {
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

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "delete",
      "import",
      "get",
      "ManageContacts",
    ]
  }
}

resource "azurerm_key_vault_certificate_contact" "test" {
  key_vault_id = azurerm_key_vault.test.id
  contact {
    email = "example@example.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMKeyVaultCertificateContact_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultCertificateContact_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate_contact" "import" {
  key_vault_id = azurerm_key_vault_certificate_contact.test.key_vault_id

  dynamic "contact" {
    for_each = azurerm_key_vault_certificate_contact.test.contact
    content {
      email = lookup(contact.value, "email", null)
    }
  }
}
`, template)
}

func testAccAzureRMKeyVaultCertificateContact_complete(data acceptance.TestData) string {
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

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "delete",
      "import",
      "get",
      "ManageContacts",
    ]
  }
}

resource "azurerm_key_vault_certificate_contact" "test" {
  key_vault_id = azurerm_key_vault.test.id
  contact {
    email = "example@example.com"
  }

  contact {
    email = "example1@example.com"
    name  = "example"
    phone = "01234567890"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
