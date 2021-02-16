package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultCertificateResource struct {
}

func TestAccKeyVaultCertificate_basicImportPFX(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicImportPFX(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
			),
		},
		data.ImportStep("certificate"),
	})
}

func TestAccKeyVaultCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicImportPFX(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_certificate"),
		},
	})
}

func TestAccKeyVaultCertificate_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicGenerate,
			TestResource: r,
		}),
	})
}

func TestAccKeyVaultCertificate_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicGenerate(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.destroyParentKeyVault, "azurerm_key_vault.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccKeyVaultCertificate_basicGenerate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicGenerate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
				check.That(data.ResourceName).Key("certificate_attribute.0.created").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_basicGenerateUnknownIssuer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicGenerateUnknownIssuer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_softDeleteRecovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.softDeleteRecovery(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
			),
		},
		{
			Config:  r.softDeleteRecovery(data, false),
			Destroy: true,
		},
		{
			Config: r.softDeleteRecovery(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
			),
		},
	})
}

func TestAccKeyVaultCertificate_basicGenerateSans(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicGenerateSans(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns.#").HasValue("1"),
			),
		},
	})
}

func TestAccKeyVaultCertificate_basicGenerateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicGenerateTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_basicExtendedKeyUsage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicExtendedKeyUsage(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.extended_key_usage.#").HasValue("3"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.extended_key_usage.0").HasValue("1.3.6.1.5.5.7.3.1"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.extended_key_usage.1").HasValue("1.3.6.1.5.5.7.3.2"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.extended_key_usage.2").HasValue("1.3.6.1.4.1.311.21.10"),
			),
		},
	})
}

func TestAccKeyVaultCertificate_emptyExtendedKeyUsage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyExtendedKeyUsage(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.extended_key_usage.#").HasValue("0"),
			),
		},
	})
}

func TestAccKeyVaultCertificate_withExternalAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withExternalAccessPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withExternalAccessPolicyUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t KeyVaultCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	keyVaultsClient := clients.KeyVault
	client := clients.KeyVault.ManagementClient

	id, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, clients.Resource, id.KeyVaultBaseUrl)
	if err != nil || keyVaultIdRaw == nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return nil, err
	}
	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil || !ok {
		return nil, fmt.Errorf("checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	cert, err := client.GetCertificate(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Key Vault Certificate: %+v", err)
	}

	return utils.Bool(cert.ID != nil), nil
}

func (KeyVaultCertificateResource) destroyParentKeyVault(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	ok, err := KeyVaultResource{}.Destroy(ctx, client, state)
	if err != nil {
		return err
	}

	if ok == nil || !*ok {
		return fmt.Errorf("deleting parent key vault failed")
	}

	return nil
}

func (KeyVaultCertificateResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	keyVaultId, err := parse.VaultID(state.Attributes["key_vault_id"])
	if err != nil {
		return nil, err
	}

	vaultBaseUrl, err := client.KeyVault.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return nil, fmt.Errorf("looking up base uri for Secret %q from id %q: %+v", name, keyVaultId, err)
	}

	if _, err := client.KeyVault.ManagementClient.DeleteCertificate(ctx, *vaultBaseUrl, name); err != nil {
		return nil, fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
	}

	return utils.Bool(true), nil
}

func (r KeyVaultCertificateResource) basicImportPFX(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/keyvaultcert.pfx")
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
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate" "import" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/keyvaultcert.pfx")
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
`, r.basicImportPFX(data))
}

func (r KeyVaultCertificateResource) basicGenerate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicGenerateUnknownIssuer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Unknown"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "EmailContacts"
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
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicGenerateSans(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicGenerateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicExtendedKeyUsage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) emptyExtendedKeyUsage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) softDeleteRecovery(data acceptance.TestData, purge bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = "%t"
      recover_soft_deleted_key_vaults = true
    }
  }
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
`, purge, r.template(data), data.RandomString)
}

func (KeyVaultCertificateResource) withExternalAccessPolicy(data acceptance.TestData) string {
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
  name                       = "acctestkeyvault%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_enabled        = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "create",
    "delete",
    "get",
    "purge",
    "recover",
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

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
  depends_on = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KeyVaultCertificateResource) withExternalAccessPolicyUpdate(data acceptance.TestData) string {
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
  name                       = "acctestkeyvault%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_enabled        = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "backup",
    "create",
    "delete",
    "get",
    "recover",
    "purge",
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

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

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
  depends_on = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KeyVaultCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkeyvault%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_enabled        = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "import",
      "purge",
      "recover",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
      "get",
    ]

    storage_permissions = [
      "set",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
