package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMKeyVaultSecret_validateName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hello",
			ExpectError: false,
		},
		{
			Input:       "hello-world",
			ExpectError: false,
		},
		{
			Input:       "hello-world-21",
			ExpectError: false,
		},
		{
			Input:       "hello_world_21",
			ExpectError: true,
		},
		{
			Input:       "Hello-World",
			ExpectError: false,
		},
		{
			Input:       "20202020",
			ExpectError: false,
		},
		{
			Input:       "ABC123!@£",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := validateKeyVaultSecretName(tc.Input, "")

		hasError := len(errors) > 0

		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Key Vault Secret Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}

func TestAccAzureRMKeyVaultSecret_parseID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    SecretID
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: SecretID{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		secretId, err := parseKeyVaultSecretID(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
			}

			return
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.KeyVaultBaseUrl != secretId.KeyVaultBaseUrl {
			t.Fatalf("Expected 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyVaultBaseUrl, secretId.KeyVaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Version != secretId.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, secretId.Version, tc.Input)
		}
	}
}

func TestAccAzureRMKeyVaultSecret_basic(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_basic(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "rick-and-morty"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_complete(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_complete(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultSecret_update(t *testing.T) {
	resourceName := "azurerm_key_vault_secret.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultSecret_basic(rs, testLocation())
	updatedConfig := testAccAzureRMKeyVaultSecret_basicUpdated(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	client := testAccProvider.Meta().(*ArmClient).keyVaultManagementClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_Secret" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]

		// get the latest version
		resp, err := client.GetSecret(vaultBaseUrl, name, "")
		if err != nil {
			if responseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Key Vault Secret still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultSecretExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]

		client := testAccProvider.Meta().(*ArmClient).keyVaultManagementClient

		resp, err := client.GetSecret(vaultBaseUrl, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on keyVaultManagementClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Key Vault Secret %q (resource group: %q) does not exist", name, vaultBaseUrl)
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

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "all",
    ]

    secret_permissions = [
      "all",
    ]
  }

  tags {
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

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "all",
    ]

    secret_permissions = [
      "all",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "<rick><morty /></rick>"
  vault_uri    = "${azurerm_key_vault.test.vault_uri}"
  content_type = "application/xml"
  tags {
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

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "all",
    ]

    secret_permissions = [
      "all",
    ]
  }

  tags {
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
