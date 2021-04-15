package frontdoor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FrontDoorCustomHttpsConfigurationResource struct {
}

func TestAccFrontDoorCustomHttpsConfiguration_CustomHttps(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_https_configuration", "test")
	r := FrontDoorCustomHttpsConfigurationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.CustomHttpsEnabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_https_provisioning_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("custom_https_configuration.0.certificate_source").HasValue("FrontDoor"),
			),
		},
		{
			Config: r.CustomHttpsDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_https_provisioning_enabled").HasValue("false"),
			),
		},
		{
			Config: r.CustomHttpsKeyVault(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_https_provisioning_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("custom_https_configuration.0.certificate_source").HasValue("AzureKeyVault"),
			),
		},
	})
}

func (FrontDoorCustomHttpsConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FrontendEndpointIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Frontdoor.FrontDoorsFrontendClient.Get(ctx, id.ResourceGroup, id.FrontDoorName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Frontend Endpoint %q (Front Door %q / Resource Group %q): %v", id.Name, id.FrontDoorName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.FrontendEndpointProperties != nil), nil
}

func (r FrontDoorCustomHttpsConfigurationResource) CustomHttpsEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_custom_https_configuration" "test" {
  frontend_endpoint_id              = azurerm_frontdoor.test.frontend_endpoints[local.endpoint_name]
  resource_group_name               = azurerm_resource_group.test.name
  custom_https_provisioning_enabled = true

  custom_https_configuration {
    certificate_source = "FrontDoor"
  }
}
`, r.template(data))
}

func (r FrontDoorCustomHttpsConfigurationResource) CustomHttpsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_custom_https_configuration" "test" {
  frontend_endpoint_id              = azurerm_frontdoor.test.frontend_endpoints[local.endpoint_name]
  resource_group_name               = azurerm_resource_group.test.name
  custom_https_provisioning_enabled = false
}
`, r.template(data))
}

func (r FrontDoorCustomHttpsConfigurationResource) CustomHttpsKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_frontdoor_custom_https_configuration" "test" {
  frontend_endpoint_id              = azurerm_frontdoor.test.frontend_endpoints[local.endpoint_name]
  custom_https_provisioning_enabled = true

  custom_https_configuration {
    certificate_source                      = "AzureKeyVault"
    azure_key_vault_certificate_vault_id    = azurerm_key_vault.test.id
    azure_key_vault_certificate_secret_name = azurerm_key_vault_certificate.test.name
  }
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-FD-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

data "azuread_service_principal" "front_door" {
  # from https://docs.microsoft.com/en-us/azure/frontdoor/front-door-custom-domain-https#prepare-your-azure-key-vault-account-and-certificate
  application_id = "ad0e1c7e-6d38-4ba4-9efd-0bc77ba9f037"
}

resource "azurerm_key_vault_access_policy" "front_door" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.front_door.object_id

  certificate_permissions = [
    "Get",
  ]

  secret_permissions = [
    "Get",
  ]
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "generated-cert"
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
        dns_names = ["internal.contoso.com", "domain.hello.world"]
      }

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (FrontDoorCustomHttpsConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%d"
  location = "%s"
}

locals {
  backend_name        = "backend-bing-custom"
  endpoint_name       = "frontend-endpoint-custom"
  health_probe_name   = "health-probe-custom"
  load_balancing_name = "load-balancing-setting-custom"
}

resource "azurerm_frontdoor" "test" {
  name                                         = "acctest-FD-%d"
  resource_group_name                          = azurerm_resource_group.test.name
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]

    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = local.backend_name
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name = local.health_probe_name
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name      = local.endpoint_name
    host_name = "acctest-FD-%d.azurefd.net"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
