// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestackOutboundTrustCertificateResource struct{}

func TestAccLocalRulestackOutboundTrustCertificateResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association", "test")
	r := LocalRulestackOutboundTrustCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LocalRulestackOutboundTrustCertificateResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.Client.CertificateObjectLocalRulestack.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRulestackOutboundTrustCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association" "test" {
  certificate_id = azurerm_palo_alto_local_rulestack_certificate.test.id
}

`, r.template(data))
}

func (r LocalRulestackOutboundTrustCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PAN-%[1]d"
  location = "%[2]s"
}

resource "azurerm_palo_alto_local_rulestack" "test" {
  name                = "testAcc-palrs-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkeyvault%[3]s"
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

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%[3]s"
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

resource "azurerm_palo_alto_local_rulestack_certificate" "test" {
  name         = "testacc-palc-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  key_vault_certificate_id = azurerm_key_vault_certificate.test.versionless_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
