package keyvault

import (
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultSecret_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "rick-and-morty"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultSecret_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "rick-and-morty"),
				),
			},
			{
				Config:      testAccAzureRMKeyVaultSecret_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_key_vault_secret"),
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					testCheckAzureRMKeyVaultSecretDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_basic(data),
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
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "not_before_date", "2019-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr(data.ResourceName, "expiration_date", "2020-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultSecret_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "rick-and-morty"),
				),
			},
			{
				Config: testAccAzureRMKeyVaultSecret_basicUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "szechuan"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_updatingValueChangedExternally(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "rick-and-morty"),
					updateKeyVaultSecretValue(data.ResourceName, "mad-scientist"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAzureRMKeyVaultSecret_updateTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
				),
			},
			{
				Config:   testAccAzureRMKeyVaultSecret_updateTags(data),
				PlanOnly: true,
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultSecret_recovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_softDeleteRecovery(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "rick-and-morty"),
				),
			},
			{
				Config:  testAccAzureRMKeyVaultSecret_softDeleteRecovery(data, false),
				Destroy: true,
			},
			{
				// purge true here to make sure when we end the test there's no soft-deleted items left behind
				Config: testAccAzureRMKeyVaultSecret_softDeleteRecovery(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "rick-and-morty"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_withExternalAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultSecret_withExternalAccessPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKeyVaultSecret_withExternalAccessPolicyUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMKeyVaultSecretDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
	vaultClient := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_secret" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]
		vaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			// key vault's been deleted
			return nil
		}

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
			return fmt.Errorf("Error looking up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
		}

		ok, err = azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
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
			return fmt.Errorf("Error looking up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
		}

		ok, err = azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
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

func updateKeyVaultSecretValue(resourceName, value string) resource.TestCheckFunc {
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
			return fmt.Errorf("Error looking up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
		}

		ok, err = azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		updated := keyvault.SecretSetParameters{
			Value: utils.String(value),
		}
		if _, err = client.SetSecret(ctx, vaultBaseUrl, name, updated); err != nil {
			return fmt.Errorf("updating secret: %+v", err)
		}
		return nil
	}
}

func testAccAzureRMKeyVaultSecret_basic(data acceptance.TestData) string {
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

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

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
  key_vault_id = azurerm_key_vault.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultSecret_withExternalAccessPolicy(data acceptance.TestData) string {
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

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  key_permissions = [
    "create",
    "get",
  ]
  secret_permissions = [
    "set",
    "get",
    "delete",
  ]
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id
  depends_on   = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultSecret_withExternalAccessPolicyUpdate(data acceptance.TestData) string {
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

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  key_permissions = [
    "create",
    "get",
  ]
  secret_permissions = [
    "set",
    "get",
    "delete",
    "recover",
  ]
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id
  depends_on   = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultSecret_updateTags(data acceptance.TestData) string {
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

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

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
  value        = "mad-scientist"
  key_vault_id = azurerm_key_vault.test.id

  tags = {
    Rick = "Morty"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultSecret_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultSecret_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_secret" "import" {
  name         = azurerm_key_vault_secret.test.name
  value        = azurerm_key_vault_secret.test.value
  key_vault_id = azurerm_key_vault_secret.test.key_vault_id
}
`, template)
}

func testAccAzureRMKeyVaultSecret_complete(data acceptance.TestData) string {
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

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

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
  key_vault_id    = azurerm_key_vault.test.id
  content_type    = "application/xml"
  not_before_date = "2019-01-01T01:02:03Z"
  expiration_date = "2020-01-01T01:02:03Z"

  tags = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultSecret_basicUpdated(data acceptance.TestData) string {
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

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

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
  name         = "secret-%s"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultSecret_softDeleteRecovery(data acceptance.TestData, purge bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = "%t"
      recover_soft_deleted_key_vaults = true
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-kvs-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  soft_delete_enabled = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
      "recover",
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
  key_vault_id = azurerm_key_vault.test.id
}
`, purge, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
