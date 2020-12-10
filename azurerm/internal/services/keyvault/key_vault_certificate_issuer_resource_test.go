package keyvault_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultCertificateIssuer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateIssuerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateIssuer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateIssuerExists(data.ResourceName),
				),
			},
			data.ImportStep("password"),
		},
	})
}

func TestAccAzureRMKeyVaultCertificateIssuer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateIssuerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateIssuer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateIssuerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMKeyVaultCertificateIssuer_requiresImport),
		},
	})
}

func TestAccAzureRMKeyVaultCertificateIssuer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateIssuerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateIssuer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateIssuerExists(data.ResourceName),
				),
			},
			data.ImportStep("password"),
		},
	})
}

func TestAccAzureRMKeyVaultCertificateIssuer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateIssuerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateIssuer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateIssuerExists(data.ResourceName),
				),
			},
			data.ImportStep("password"),
			{
				Config: testAccAzureRMKeyVaultCertificateIssuer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateIssuerExists(data.ResourceName),
				),
			},
			data.ImportStep("password"),
		},
	})
}

func TestAccAzureRMKeyVaultCertificateIssuer_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateIssuerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateIssuer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateIssuerExists(data.ResourceName),
					testCheckAzureRMKeyVaultCertificateIssuerDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificateIssuer_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateIssuerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificateIssuer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateIssuerExists("azurerm_key_vault_certificate_issuer.test"),
					testCheckAzureRMKeyVaultDisappears("azurerm_key_vault.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMKeyVaultCertificateIssuerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
	vaultClient := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_certificate_issuer" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]
		vaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			// deleted, this is fine.
			return nil
		}

		ok, err := azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to check if key vault %q for Certificate Issuer %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Certificate Issuer %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		// get the latest version
		resp, err := client.GetCertificateIssuer(ctx, vaultBaseUrl, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return fmt.Errorf("Bad: Get on keyVault certificate issuer: %+v", err)
		}

		return fmt.Errorf("Key Vault Certificate Issuer still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultCertificateIssuerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
		vaultClient := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]
		vaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to look up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
		}

		ok, err = azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to check if key vault %q for Certificate Issuer %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Certificate Issuer %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.GetCertificateIssuer(ctx, vaultBaseUrl, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault Certificate Issuer %q (resource group: %q) does not exist", name, vaultBaseUrl)
			}

			return fmt.Errorf("Bad: Get on keyVault certificate issuer: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultCertificateIssuerDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
		vaultClient := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		name := rs.Primary.Attributes["name"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]
		vaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to look up base URI from id %q: %+v", keyVaultId, err)
		}

		ok, err = azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("failed to check if key vault %q for Certificate Issuer %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Certificate Issuer %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.DeleteCertificateIssuer(ctx, vaultBaseUrl, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultCertificateIssuer_basic(data acceptance.TestData) string {
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
      "manageissuers",
      "setissuers",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate_issuer" "test" {
  name          = "acctestKVCI-%d"
  key_vault_id  = azurerm_key_vault.test.id
  provider_name = "OneCertV2-PrivateCA"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMKeyVaultCertificateIssuer_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultCertificateIssuer_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate_issuer" "import" {
  name          = azurerm_key_vault_certificate_issuer.test.name
  key_vault_id  = azurerm_key_vault_certificate_issuer.test.key_vault_id
  org_id        = azurerm_key_vault_certificate_issuer.test.org_id
  account_id    = azurerm_key_vault_certificate_issuer.test.account_id
  password      = "test"
  provider_name = azurerm_key_vault_certificate_issuer.test.provider_name
}

`, template)
}

func testAccAzureRMKeyVaultCertificateIssuer_complete(data acceptance.TestData) string {
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
      "manageissuers",
      "setissuers",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate_issuer" "test" {
  name          = "acctestKVCI-%d"
  key_vault_id  = azurerm_key_vault.test.id
  account_id    = "test-account"
  password      = "test"
  provider_name = "DigiCert"

  org_id = "accTestOrg"
  admin {
    email_address = "admin@contoso.com"
    first_name    = "First"
    last_name     = "Last"
    phone         = "01234567890"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
