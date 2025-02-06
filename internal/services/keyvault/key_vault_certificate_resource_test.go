// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultCertificateResource struct{}

func TestAccKeyVaultCertificate_basicImportPFX(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicImportPFX(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("resource_manager_id").Exists(),
				check.That(data.ResourceName).Key("resource_manager_versionless_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.secret_properties.0.content_type").HasValue("application/x-pkcs12"),
				check.That(data.ResourceName).Key("versionless_id").HasValue(fmt.Sprintf("https://acctestkeyvault%s.vault.azure.net/certificates/acctestcert%s", data.RandomString, data.RandomString)),
				check.That(data.ResourceName).Key("versionless_secret_id").HasValue(fmt.Sprintf("https://acctestkeyvault%s.vault.azure.net/secrets/acctestcert%s", data.RandomString, data.RandomString)),
			),
		},
		data.ImportStep("certificate"),
	})
}

func TestAccKeyVaultCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicImportPFX(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicGenerate,
			TestResource: r,
		}),
	})
}

func TestAccKeyVaultCertificate_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerate(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.destroyParentKeyVault, "azurerm_key_vault.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccKeyVaultCertificate_basicGenerate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerate(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccKeyVaultCertificate_updateLifeTime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicGenerateUpdateLifetimeAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_basicGenerateUnknownIssuer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerateUnknownIssuer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_softDeleteRecovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.softDeleteRecovery(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
			),
		},
		{
			Config: r.softDeleteCertificate(data, false),
		},
		{
			Config: r.softDeleteRecovery(data, true),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerateSans(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerateTags(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccKeyVaultCertificate_updateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerateTags(data),
		},
		{
			Config: r.updateTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
				check.That(data.ResourceName).Key("tags.hello").DoesNotExist(),
			),
		},
	})
}

func TestAccKeyVaultCertificate_updateCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerateCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_size").HasValue("2048"),
			),
		},
		{
			Config: r.updateCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_size").HasValue("4096"),
			),
		},
	})
}

func TestAccKeyVaultCertificate_basicGenerateEllipticCurve(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerateEllipticCurve(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
				check.That(data.ResourceName).Key("certificate_attribute.0.created").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.curve").HasValue("P-256K"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_type").HasValue("EC"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_size").HasValue("256"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_basicGenerateEllipticCurveAutoKeySize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerateEllipticCurveAutoKeySize(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
				check.That(data.ResourceName).Key("certificate_attribute.0.created").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.curve").HasValue("P-521"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_type").HasValue("EC"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_size").HasValue("521"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_basicExtendedKeyUsage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicExtendedKeyUsage(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyExtendedKeyUsage(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withExternalAccessPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withExternalAccessPolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultCertificate_purge(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicGenerate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret_id").Exists(),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
				check.That(data.ResourceName).Key("certificate_attribute.0.created").Exists(),
			),
		},
		{
			Config:  r.basicGenerate(data),
			Destroy: true,
		},
	})
}

func TestAccKeyVaultCertificate_unorderedKeyUsage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.unorderedKeyUsage(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccKeyVaultCertificate_updatedImportedCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicImportPFX(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
			),
		},
		data.ImportStep("certificate"),
		{
			Config: r.basicImportPFX_ECDSA(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
			),
		},
		data.ImportStep("certificate"),
	})
}

func (t KeyVaultCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.KeyVault
	subscriptionId := clients.Account.SubscriptionId

	id, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := client.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil || keyVaultIdRaw == nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return nil, err
	}
	ok, err := client.Exists(ctx, *keyVaultId)
	if err != nil || !ok {
		return nil, fmt.Errorf("checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	cert, err := client.ManagementClient.GetCertificate(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Key Vault Certificate: %+v", err)
	}

	return utils.Bool(cert.ID != nil), nil
}

func (KeyVaultCertificateResource) destroyParentKeyVault(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ok, err := KeyVaultResource{}.Destroy(ctx, client, state)
	if err != nil {
		return err
	}

	if ok == nil || !*ok {
		return fmt.Errorf("deleting parent key vault failed")
	}

	return nil
}

func (KeyVaultCertificateResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	keyVaultId, err := commonids.ParseKeyVaultID(state.Attributes["key_vault_id"])
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
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicImportPEM_ECDSA(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/ecdsa.pem")
    password = ""
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicImportPFX_ECDSA(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/ecdsa.pfx")
    password = ""
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicImportPFX_RSA_bundle(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/rsa_bundle.pfx")
    password = ""
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicImportPEM_RSA(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/rsa_single.pem")
    password = ""
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicImportPEM_RSA_bundle(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/rsa_bundle.pem")
    password = ""
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
        "keyEncipherment",
        "keyCertSign",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicGenerateCertificate(data acceptance.TestData) string {
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

func (r KeyVaultCertificateResource) basicGenerateUpdateLifetimeAction(data acceptance.TestData) string {
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
        "keyEncipherment",
        "keyCertSign",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) updateCertificate(data acceptance.TestData) string {
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
      key_size   = 4096
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
        lifetime_percentage = 30
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

func (r KeyVaultCertificateResource) basicGenerateEllipticCurve(data acceptance.TestData) string {
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
      curve      = "P-256K"
      exportable = true
      key_size   = 256
      key_type   = "EC"
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
        "digitalSignature",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) basicGenerateEllipticCurveAutoKeySize(data acceptance.TestData) string {
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
      curve      = "P-521"
      exportable = true
      key_type   = "EC"
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
        "digitalSignature",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
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

func (r KeyVaultCertificateResource) softDeleteCertificate(data acceptance.TestData, purge bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_certificates_on_destroy = %t
      recover_soft_deleted_key_vaults = true
    }
  }
}

%s`, purge, r.template(data))
}

func (r KeyVaultCertificateResource) softDeleteRecovery(data acceptance.TestData, purge bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_certificates_on_destroy = %t
      recover_soft_deleted_key_vaults            = true
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
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
  ]

  key_permissions = [
    "Create",
  ]

  secret_permissions = [
    "Set",
  ]

  storage_permissions = [
    "Set",
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
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Backup",
    "Create",
    "Delete",
    "Get",
    "Recover",
    "Purge",
    "Update",
  ]

  key_permissions = [
    "Create",
  ]

  secret_permissions = [
    "Set",
  ]

  storage_permissions = [
    "Set",
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
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Import",
      "Purge",
      "Recover",
      "Update",
      "List",
    ]

    key_permissions = [
      "Create",
    ]

    secret_permissions = [
      "Get",
      "Set",
    ]

    storage_permissions = [
      "Set",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r KeyVaultCertificateResource) updateTags(data acceptance.TestData) string {
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
    ENV = "Test"
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultCertificateResource) unorderedKeyUsage(data acceptance.TestData) string {
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
        "digitalSignature",
        "cRLSign",
        "dataEncipherment",
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
