package firewall_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FirewallPolicyTLSCertificateResource struct {
}

func TestAccFirewallPolicyTLSCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyTLSCertificateResource{}

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

func TestAccFirewallPolicyTLSCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyTLSCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyTLSCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyTLSCertificateResource{}

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

func (FirewallPolicyTLSCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	var id, err = parse.FirewallPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Firewall.FirewallPolicyClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.FirewallPolicyPropertiesFormat != nil), nil
}

func (FirewallPolicyTLSCertificateResource) basic(data acceptance.TestData) string {
	r := FirewallPolicyTLSCertificateResource{}
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test" {
  name                     = "acctest-networkfw-Policy-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  threat_intelligence_mode = "Off"
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
    fqdns        = ["foo.com", "bar.com"]
  }
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  intrusion_detection {
    mode = "Alert"
    signature_overrides {
      state = "Alert"
      id    = "1"
    }
    traffic_bypass {
      name                  = "Name bypass traffic settings"
      description           = "Description bypass traffic settings"
      protocol              = "ANY"
      destination_addresses = ["*"]
      destination_ports     = ["*"]
      source_ip_groups      = ["*"]
      destination_ip_groups = ["*"]
    }
  }
  identity {
    type = "SystemAssigned"
  }
  tls_certificate {
    key_vault_secret_id = azurerm_key_vault_certificate.test.secret_id
    name                = azurerm_key_vault_certificate.test.name
  }
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger)
}

func (FirewallPolicyTLSCertificateResource) requiresImport(data acceptance.TestData) string {
	r := FirewallPolicyTLSCertificateResource{}
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "import" {
  name                = azurerm_firewall_policy.test.name
  resource_group_name = azurerm_firewall_policy.test.resource_group_name
  location            = azurerm_firewall_policy.test.location
}
`, template)
}

func (FirewallPolicyTLSCertificateResource) inherit(data acceptance.TestData) string {
	r := FirewallPolicyTLSCertificateResource{}
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy" "test-parent" {
  name                = "acctest-networkfw-Policy-%d-parent"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_policy_id      = azurerm_firewall_policy.test-parent.id
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
    fqdns        = ["foo.com", "bar.com"]
  }
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (FirewallPolicyTLSCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-networkfw-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                            = "tlskv%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  tenant_id                       = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id   = azurerm_key_vault.test.id
  application_id = data.azurerm_client_config.current.client_id
  tenant_id      = data.azurerm_client_config.current.tenant_id
  object_id      = data.azurerm_client_config.current.object_id

  key_permissions = [
    "backup",
    "create",
    "delete",
    "get",
    "import",
    "list",
    "purge",
    "recover",
    "restore",
    "update"
  ]

  certificate_permissions = [
    "backup",
    "create",
    "get",
    "list",
    "import",
    "purge",
    "delete",
    "recover",
  ]

  secret_permissions = [
    "get",
    "list",
    "set",
    "purge",
    "delete",
    "recover"
  ]

}

resource "azurerm_key_vault_certificate" "test" {
  name         = "AzureFirewallPolicyTLSCertificateCertificate"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/cert_key.pem")
  }

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
    secret_properties {
      content_type = "application/x-pem-file"
    }
    x509_certificate_properties {
      # Server Authentication = 1.3.6.1.5.5.7.3.1
      # Client Authentication = 1.3.6.1.5.5.7.3.2
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject_alternative_names {
        dns_names = ["api.pluginsdk.io"]
      }
      subject            = "CN=api.pluginsdk.io"
      validity_in_months = 1
    }
  }

}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
