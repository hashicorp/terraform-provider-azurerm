package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultKey_basicEC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_basicEC(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
				),
			},
			data.ImportStep("key_size"),
		},
	})
}

func TestAccAzureRMKeyVaultKey_basicECClassic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_basicECClassic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
				),
			},
			data.ImportStep("key_size"),
		},
	})
}

func TestAccAzureRMKeyVaultKey_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultKey_basicEC(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMKeyVaultKey_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_key_vault_key"),
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_basicECHSM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_basicECHSM(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_curveEC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_curveEC(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultKey_basicRSA(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_basicRSA(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
				),
			},
			data.ImportStep("key_size"),
		},
	})
}

func TestAccAzureRMKeyVaultKey_basicRSAHSM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_basicRSAHSM(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
				),
			},
			data.ImportStep("key_size"),
		},
	})
}

func TestAccAzureRMKeyVaultKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_complete(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
			data.ImportStep("key_size"),
		},
	})
}

func TestAccAzureRMKeyVaultKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	config := testAccAzureRMKeyVaultKey_basicRSA(data)
	updatedConfig := testAccAzureRMKeyVaultKey_basicUpdated(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_opts.#", "6"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_opts.0", "decrypt"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_opts.#", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_opts.0", "encrypt"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultKey_basicEC(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(data.ResourceName),
					testCheckAzureRMKeyVaultKeyDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultKey_basicEC(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists("azurerm_key_vault_key.test"),
					testCheckAzureRMKeyVaultDisappears("azurerm_key_vault.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMKeyVaultKeyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_key" {
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
		resp, err := client.GetKey(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Key Vault Key still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultKeyExists(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Error checking if key vault %q for Key %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.GetKey(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault Key %q (resource group: %q) does not exist", name, vaultBaseUrl)
			}

			return fmt.Errorf("Bad: Get on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultKeyDisappears(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Error checking if key vault %q for Key %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.DeleteKey(ctx, vaultBaseUrl, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultKey_basicEC(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
      "delete",
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

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = "${azurerm_key_vault.test.id}"
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultKey_basicECClassic(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
      "delete",
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

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "EC"
  key_size  = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultKey_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultKey_basicEC(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "import" {
  name         = "${azurerm_key_vault_key.test.name}"
  key_vault_id = "${azurerm_key_vault.test.id}"
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}
`, template)
}

func testAccAzureRMKeyVaultKey_basicRSA(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = "${azurerm_key_vault.test.id}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultKey_basicRSAHSM(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
      "delete",
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

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = "${azurerm_key_vault.test.id}"
  key_type     = "RSA-HSM"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultKey_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
      "delete",
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

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = "${azurerm_key_vault.test.id}"
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

  tags = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultKey_basicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "RSA"
  key_size  = 2048

  key_opts = [
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultKey_curveEC(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
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

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "EC"
  curve     = "P-521"

  key_opts = [
    "sign",
    "verify",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultKey_basicECHSM(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
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

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "EC-HSM"
  curve     = "P-521"

  key_opts = [
    "sign",
    "verify",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
