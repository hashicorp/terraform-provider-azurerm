// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRuleResource struct{}

func TestAccPaloAltoLocalRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_rule", "test")

	r := LocalRuleResource{}

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

func TestAccPaloAltoLocalRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_rule", "test")

	r := LocalRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPaloAltoLocalRule_withDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_rule", "test")

	r := LocalRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDestination(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_rule", "test")

	r := LocalRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_rule", "test")

	r := LocalRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRule_updateProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_rule", "test")

	r := LocalRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicProtocol(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicProtocolPorts(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicProtocol(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LocalRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := localrules.ParseLocalRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.Client.LocalRules.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  priority     = 100
  action       = "Allow"
  protocol     = "application-default"

  applications = ["any"]

  destination {
    cidrs = ["any"]
  }

  source {
    cidrs = ["any"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LocalRuleResource) basicProtocol(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  priority     = 100
  action       = "Allow"

  applications = ["any"]

  protocol = "TCP:8080"

  destination {
    cidrs = ["any"]
  }

  source {
    cidrs = ["any"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LocalRuleResource) basicProtocolPorts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  priority     = 100
  action       = "Allow"

  applications = ["any"]

  protocol_ports = ["TCP:8080", "TCP:8081"]

  destination {
    cidrs = ["any"]
  }

  source {
    cidrs = ["any"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LocalRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_palo_alto_local_rulestack_rule" "import" {
  name         = azurerm_palo_alto_local_rulestack_rule.test.name
  rulestack_id = azurerm_palo_alto_local_rulestack_rule.test.rulestack_id
  priority     = azurerm_palo_alto_local_rulestack_rule.test.priority
  action       = "Allow"
  applications = azurerm_palo_alto_local_rulestack_rule.test.applications
  protocol     = azurerm_palo_alto_local_rulestack_rule.test.protocol

  destination {
    cidrs = azurerm_palo_alto_local_rulestack_rule.test.destination.0.cidrs
  }

  source {
    cidrs = azurerm_palo_alto_local_rulestack_rule.test.source.0.cidrs
  }
}
`, r.basic(data), data.RandomInteger)
}

func (r LocalRuleResource) withDestination(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  priority     = 100
  action       = "Allow"
  protocol     = "application-default"

  applications = ["any"]

  destination {
    countries = ["US", "GB"]
  }

  source {
    countries = ["US", "GB"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LocalRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  priority     = 100

  action        = "DenySilent"
  applications  = ["any"]
  audit_comment = "test audit comment"

  category {
    custom_urls = ["hacking"] // TODO - This is another resource type in PAN?
  }

  decryption_rule_type = "SSLOutboundInspection"
  description          = "Acceptance Test Rule - dated %[2]d"

  destination {
    countries                       = ["US", "GB"]
    local_rulestack_fqdn_list_ids   = [azurerm_palo_alto_local_rulestack_fqdn_list.test.id]
    local_rulestack_prefix_list_ids = [azurerm_palo_alto_local_rulestack_prefix_list.test.id]
  }

  logging_enabled = false

  inspection_certificate_id = azurerm_palo_alto_local_rulestack_certificate.test.id

  negate_destination = true
  negate_source      = true

  protocol = "TCP:8080"

  enabled = false

  source {
    countries                       = ["US", "GB"]
    local_rulestack_prefix_list_ids = [azurerm_palo_alto_local_rulestack_prefix_list.test.id]
  }

  tags = {
    "acctest" = "true"
    "foo"     = "bar"
  }

  depends_on = [
    azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association.test,
    azurerm_palo_alto_local_rulestack_outbound_untrust_certificate_association.test
  ]
}
`, r.templateWithCertsEnabled(data), data.RandomInteger)
}

func (r LocalRuleResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_rule" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  priority     = 100

  action        = "DenySilent"
  applications  = ["any"]
  audit_comment = "test audit comment"

  category {
    custom_urls = ["web-based-email", "social-networking"]
  }

  description = "Acceptance Test Rule - updated %[2]d"

  destination {
    countries = ["US", "GB"]
  }

  logging_enabled = false

  inspection_certificate_id = azurerm_palo_alto_local_rulestack_certificate.test.id

  negate_destination = false
  negate_source      = false

  protocol_ports = ["TCP:8080", "TCP:8081"]

  enabled = true

  source {
    countries = ["US", "GB"]
  }

  tags = {
    "acctest" = "true"
    "foo"     = "bar"
  }
}
`, r.templateWithCertsEnabled(data), data.RandomInteger)
}

func (r LocalRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PAN-%[1]d"
  location = "%[2]s"
}

resource "azurerm_palo_alto_local_rulestack_certificate" "test" {
  name         = "testacc-palc-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  self_signed  = true
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

resource "azurerm_key_vault_certificate" "trust" {
  name         = "acctesttrust%[3]s"
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

resource "azurerm_key_vault_certificate" "untrust" {
  name         = "acctestuntrust%[3]s"
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

resource "azurerm_palo_alto_local_rulestack_certificate" "trust" {
  name         = "testacc-palcT-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  key_vault_certificate_id = azurerm_key_vault_certificate.trust.versionless_id
}

resource "azurerm_palo_alto_local_rulestack_certificate" "untrust" {
  name         = "testacc-palcU-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  key_vault_certificate_id = azurerm_key_vault_certificate.untrust.versionless_id
}

resource "azurerm_palo_alto_local_rulestack" "test" {
  name                = "testAcc-palrs-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LocalRuleResource) templateWithCertsEnabled(data acceptance.TestData) string {
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

resource "azurerm_palo_alto_local_rulestack_certificate" "test" {
  name         = "testacc-palc-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id
  self_signed  = true
}

resource "azurerm_palo_alto_local_rulestack_fqdn_list" "test" {
  name         = "testacc-pafqdn-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  fully_qualified_domain_names = ["contoso.com", "test.example.com", "anothertest.example.com"]
}

resource "azurerm_palo_alto_local_rulestack_prefix_list" "test" {
  name         = "testacc-palr-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  prefix_list = ["10.0.0.0/8", "172.16.0.0/16"]
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

resource "azurerm_key_vault_certificate" "trust" {
  name         = "acctesttrust%[3]s"
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

resource "azurerm_key_vault_certificate" "untrust" {
  name         = "acctestuntrust%[3]s"
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

resource "azurerm_palo_alto_local_rulestack_certificate" "trust" {
  name         = "testacc-palcT-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  key_vault_certificate_id = azurerm_key_vault_certificate.trust.versionless_id
}

resource "azurerm_palo_alto_local_rulestack_certificate" "untrust" {
  name         = "testacc-palcU-%[1]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  key_vault_certificate_id = azurerm_key_vault_certificate.untrust.versionless_id
}

resource "azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association" "test" {
  certificate_id = azurerm_palo_alto_local_rulestack_certificate.trust.id
}

resource "azurerm_palo_alto_local_rulestack_outbound_untrust_certificate_association" "test" {
  certificate_id = azurerm_palo_alto_local_rulestack_certificate.untrust.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
