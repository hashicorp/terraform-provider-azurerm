package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultCertificate_basicImportPFX(t *testing.T) {
	resourceName := "azurerm_key_vault_certificate.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultCertificate_basicImportPFX(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_data"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_basicGenerate(t *testing.T) {
	resourceName := "azurerm_key_vault_certificate.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultCertificate_basicGenerate(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "secret_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_data"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_basicGenerateTags(t *testing.T) {
	resourceName := "azurerm_key_vault_certificate.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultCertificate_basicGenerateTags(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_data"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func testCheckAzureRMKeyVaultCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).keyVaultManagementClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]

		// get the latest version
		resp, err := client.GetCertificate(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Key Vault Certificate still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultCertificateExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]

		client := testAccProvider.Meta().(*ArmClient).keyVaultManagementClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.GetCertificate(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault Certificate %q (resource group: %q) does not exist", name, vaultBaseUrl)
			}

			return fmt.Errorf("Bad: Get on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultCertificate_basicImportPFX(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "delete",
      "import",
      "get",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name      = "acctestcert%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"

  certificate {
    contents = "${base64encode(file("testdata/keyvaultcert.pfx"))}"
    password = ""
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultCertificate_basicGenerate(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update"
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name      = "acctestcert%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}

`, rString, location, rString, rString)
}
func testAccAzureRMKeyVaultCertificate_basicGenerateTags(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update"
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name      = "acctestcert%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }

  tags {
    "hello" = "world"
  }
}

`, rString, location, rString, rString)
}
