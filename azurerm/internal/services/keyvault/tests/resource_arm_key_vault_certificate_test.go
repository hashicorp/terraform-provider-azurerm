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

func TestAccAzureRMKeyVaultCertificate_basicImportPFX(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicImportPFX(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
				),
			},
			data.ImportStep("certificate"),
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_basicImportPFXClassic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicImportPFXClassic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
				),
			},
			data.ImportStep("certificate"),
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultCertificate_basicImportPFX(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
				),
			},
			{
				Config:      testAccAzureRMKeyVaultCertificate_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_key_vault_certificate"),
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicGenerate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					testCheckAzureRMKeyVaultCertificateDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicGenerate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists("azurerm_key_vault_certificate.test"),
					testCheckAzureRMKeyVaultDisappears("azurerm_key_vault.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_basicGenerate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicGenerate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secret_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_basicGenerateSans(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicGenerateSans(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_basicGenerateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicGenerateTags(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_basicExtendedKeyUsage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_basicExtendedKeyUsage(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.extended_key_usage.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.extended_key_usage.0", "1.3.6.1.5.5.7.3.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.extended_key_usage.1", "1.3.6.1.5.5.7.3.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.extended_key_usage.2", "1.3.6.1.4.1.311.21.10"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_emptyExtendedKeyUsage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	config := testAccAzureRMKeyVaultCertificate_emptyExtendedKeyUsage(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.extended_key_usage.#", "0"),
				),
			},
		},
	})
}

func testCheckAzureRMKeyVaultCertificateDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.ManagementClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]
		keyVaultId := rs.Primary.Attributes["key_vault_id"]

		ok, err := azure.KeyVaultExists(ctx, acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error checking if key vault %q for Certificate %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Certificate %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		// get the latest version
		resp, err := client.GetCertificate(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return fmt.Errorf("Bad: Get on keyVault certificate: %+v", err)
		}

		return fmt.Errorf("Key Vault Certificate still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultCertificateExists(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Error checking if key vault %q for Certificate %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Certificate %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.GetCertificate(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault Certificate %q (resource group: %q) does not exist", name, vaultBaseUrl)
			}

			return fmt.Errorf("Bad: Get on keyVault certificate: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultCertificateDisappears(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Error checking if key vault %q for Certificate %q in Vault at url %q exists: %v", keyVaultId, name, vaultBaseUrl, err)
		}
		if !ok {
			log.Printf("[DEBUG] Certificate %q Key Vault %q was not found in Key Vault at URI %q ", name, keyVaultId, vaultBaseUrl)
			return nil
		}

		resp, err := client.DeleteCertificate(ctx, vaultBaseUrl, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultCertificate_basicImportPFX(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

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
  name         = "acctestcert%s"
  key_vault_id = "${azurerm_key_vault.test.id}"

  certificate {
    contents = "${filebase64("testdata/keyvaultcert.pfx")}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultCertificate_basicImportPFXClassic(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

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

    storage_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name      = "acctestcert%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"

  certificate {
    contents = "${filebase64("testdata/keyvaultcert.pfx")}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultCertificate_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultCertificate_basicImportPFX(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate" "import" {
  name         = "${azurerm_key_vault_certificate.test.name}"
  key_vault_id = "${azurerm_key_vault.test.id}"

  certificate {
    contents = "${filebase64("testdata/keyvaultcert.pfx")}"
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
`, template)
}

func testAccAzureRMKeyVaultCertificate_basicGenerate(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]

    storage_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = "${azurerm_key_vault.test.id}"

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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultCertificate_basicGenerateSans(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]

    storage_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = "${azurerm_key_vault.test.id}"

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

      subject = "CN=hello-world"

      subject_alternative_names {
        emails    = ["mary@stu.co.uk"]
        dns_names = ["internal.contoso.com"]
        upns      = ["john@doe.com"]
      }

      validity_in_months = 12
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultCertificate_basicGenerateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
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
  name         = "acctestcert%s"
  key_vault_id = "${azurerm_key_vault.test.id}"

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

  tags = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultCertificate_basicExtendedKeyUsage(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]

    storage_permissions = [
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
      extended_key_usage = [
        "1.3.6.1.5.5.7.3.1",     # Server Authentication
        "1.3.6.1.5.5.7.3.2",     # Client Authentication
        "1.3.6.1.4.1.311.21.10", # Application Policies
      ]

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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMKeyVaultCertificate_emptyExtendedKeyUsage(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]

    storage_permissions = [
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
      extended_key_usage = []

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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
