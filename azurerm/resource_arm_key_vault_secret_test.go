package azurerm

import (
	"fmt"
	"log"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultSecret_basic(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_basic(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "rick-and-morty"),
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

func TestAccAzureRMKeyVaultSecret_basicClassic(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_basicClasic(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "rick-and-morty"),
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

func TestAccAzureRMKeyVaultSecret_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_basic(rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "rick-and-morty"),
				),
			},
			{
				Config:      testAccAzureRMKeyVaultSecret_requiresImport(rs, location),
				ExpectError: acceptance.RequiresImportError("azurerm_key_vault_secret"),
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_disappears(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_basic(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					testCheckAzureRMKeyVaultSecretDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_basic(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists("azurerm_key_vault_secret.test"),
					testCheckAzureRMKeyVaultDisappears("azurerm_key_vault.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_complete(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_complete(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "not_before_date", "2019-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr(resourceName, "expiration_date", "2020-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
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

func TestAccAzureRMKeyVaultSecret_update(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_basic(rs, acceptance.Location())
	updatedConfig := testAccAzureRMKeyVaultSecret_basicUpdated(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "rick-and-morty"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "szechuan"),
				),
			},
		},
	})
}

func testCheckAzureRMKeyVaultSecretDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_secret" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]

		ok, err := azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		// get the latest version
		resp, err := client.GetSecret(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Key Vault Secret still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultSecretExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]

		ok, err := azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.GetSecret(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault Secret %q (resource group: %q) does not exist", name, vaultBaseUrl)
			}

			return fmt.Errorf("Bad: Get on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultSecretDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]

		ok, err := azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.DeleteSecret(ctx, vaultBaseUrl, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultSecret_basic(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "premium"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "rick-and-morty"
  key_vault_id = "${azurerm_key_vault.test.id}"
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultSecret_basicClasic(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "premium"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name      = "secret-%s"
  value     = "rick-and-morty"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultSecret_requiresImport(rString string, location string) string {
	template := testAccAzureRMKeyVaultSecret_basic(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_secret" "import" {
  name         = "${azurerm_key_vault_secret.test.name}"
  value        = "${azurerm_key_vault_secret.test.value}"
  key_vault_id = "${azurerm_key_vault_secret.test.key_vault_id}"
}
`, template)
}

func testAccAzureRMKeyVaultSecret_complete(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "premium"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name            = "secret-%s"
  value           = "<rick><morty /></rick>"
  key_vault_id    = "${azurerm_key_vault.test.id}"
  content_type    = "application/xml"
  not_before_date = "2019-01-01T01:02:03Z"
  expiration_date = "2020-01-01T01:02:03Z"

  tags = {
    "hello" = "world"
  }
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultSecret_basicUpdated(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "premium"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name      = "secret-%s"
  value     = "szechuan"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
}
`, rString, location, rString, rString)
}
