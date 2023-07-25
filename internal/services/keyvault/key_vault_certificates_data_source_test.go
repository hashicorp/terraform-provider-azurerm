// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KeyVaultCertificatesDataSource struct{}

func TestAccDataSourceKeyVaultCertificates_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificates", "test")
	r := KeyVaultCertificatesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("names.#").HasValue("31"),
				check.That(data.ResourceName).Key("certificates.#").HasValue("31"),
			),
		},
	})
}

func (KeyVaultCertificatesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate" "test2" {
  count = 30
  name  = "certificate-${count.index}"
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
  key_vault_id = azurerm_key_vault.test.id
}

data "azurerm_key_vault_certificates" "test" {
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_certificate.test, azurerm_key_vault_certificate.test2]
}
`, KeyVaultCertificateResource{}.basicGenerate(data))
}
