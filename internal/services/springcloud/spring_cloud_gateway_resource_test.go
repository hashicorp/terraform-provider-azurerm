// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudGatewayResource struct{}

func TestAccSpringCloudGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
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

func TestAccSpringCloudGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
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

func TestAccSpringCloudGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sso.0.client_id", "sso.0.client_secret", "sensitive_environment_variables.%", "sensitive_environment_variables.NEW_RELIC_APP_NAME"),
	})
}

func TestAccSpringCloudGateway_clientAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sso.0.client_id", "sso.0.client_secret", "sensitive_environment_variables.%", "sensitive_environment_variables.NEW_RELIC_APP_NAME"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudGateway_responseCache(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.responseCachePerInstance(data, "10MB", "30s"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.responseCachePerRoute(data, "900KB", "5m"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudGateway_apms(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.elasticApm(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudGateway_apmsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticApm(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dynatraceApm(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleApms(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.noApms(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudGatewayResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appplatform.ParseGatewayID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.AppPlatformClient.GatewaysGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(true), nil
}

func (r SpringCloudGatewayResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudGatewayResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, template)
}

func (r SpringCloudGatewayResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_gateway" "import" {
  name                    = azurerm_spring_cloud_gateway.test.name
  spring_cloud_service_id = azurerm_spring_cloud_gateway.test.spring_cloud_service_id
}
`, config)
}

func (r SpringCloudGatewayResource) complete(data acceptance.TestData, clientId, clientSecret string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
data "azurerm_client_config" "current" {
}

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id

  https_only                               = false
  public_network_access_enabled            = false
  instance_count                           = 2
  application_performance_monitoring_types = ["ApplicationInsights", "NewRelic"]

  api_metadata {
    description       = "test description"
    documentation_url = "https://www.test.com/docs"
    server_url        = "https://www.test.com"
    title             = "test title"
    version           = "1.0"
  }

  cors {
    credentials_allowed     = false
    allowed_headers         = ["*"]
    allowed_methods         = ["PUT"]
    allowed_origins         = ["test.com"]
    allowed_origin_patterns = ["test*.com"]
    exposed_headers         = ["x-test-header"]
    max_age_seconds         = 86400
  }

  environment_variables = {
    APPLICATIONINSIGHTS_SAMPLE_RATE = "10"
  }

  sensitive_environment_variables = {
    NEW_RELIC_APP_NAME = "scg-asa"
  }

  quota {
    cpu    = "1"
    memory = "2Gi"
  }

  sso {
    client_id     = "%s"
    client_secret = "%s"
    issuer_uri    = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
    scope         = ["read"]
  }
}
`, template, clientId, clientSecret)
}

func (r SpringCloudGatewayResource) apmsTemplate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_dynatrace_application_performance_monitoring" "test" {
  name                    = "acctest-dapm-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  tenant                  = "test-tenant"
  tenant_token            = "dt0s01.AAAAAAAAAAAAAAAAAAAAAAAA.BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
  connection_point        = "https://example.live.dynatrace.com:443"
}

resource "azurerm_spring_cloud_elastic_application_performance_monitoring" "test" {
  name                    = "acctest-eapm-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  application_packages    = ["org.example", "org.another.example"]
  service_name            = "test-service-name"
  server_url              = "http://127.0.0.1:8200"
}
`, template, data.RandomInteger)
}

func (r SpringCloudGatewayResource) dynatraceApm(data acceptance.TestData) string {
	template := r.apmsTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_gateway" "test" {
  name                                   = "default"
  spring_cloud_service_id                = azurerm_spring_cloud_service.test.id
  application_performance_monitoring_ids = [azurerm_spring_cloud_dynatrace_application_performance_monitoring.test.id]
}
`, template)
}

func (r SpringCloudGatewayResource) elasticApm(data acceptance.TestData) string {
	template := r.apmsTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_gateway" "test" {
  name                                   = "default"
  spring_cloud_service_id                = azurerm_spring_cloud_service.test.id
  application_performance_monitoring_ids = [azurerm_spring_cloud_elastic_application_performance_monitoring.test.id]
}
`, template)
}

func (r SpringCloudGatewayResource) multipleApms(data acceptance.TestData) string {
	template := r.apmsTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_gateway" "test" {
  name                                   = "default"
  spring_cloud_service_id                = azurerm_spring_cloud_service.test.id
  application_performance_monitoring_ids = [azurerm_spring_cloud_dynatrace_application_performance_monitoring.test.id, azurerm_spring_cloud_elastic_application_performance_monitoring.test.id]
}
`, template)
}

func (r SpringCloudGatewayResource) noApms(data acceptance.TestData) string {
	template := r.apmsTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, template)
}

func (r SpringCloudGatewayResource) clientAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Resource Provider"
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Set",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.test.object_id

    secret_permissions = [
      "Get",
      "List",
    ]

    certificate_permissions = [
      "Get",
      "List",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest-cert-%[2]d"
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

      subject            = "CN=contoso.com"
      validity_in_months = 12
    }
  }
}

resource "azurerm_spring_cloud_certificate" "test" {
  name                     = "acctest-scc-%[2]d"
  resource_group_name      = azurerm_spring_cloud_service.test.resource_group_name
  service_name             = azurerm_spring_cloud_service.test.name
  key_vault_certificate_id = azurerm_key_vault_certificate.test.id
  exclude_private_key      = true
}

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  client_authorization {
    certificate_ids      = [azurerm_spring_cloud_certificate.test.id]
    verification_enabled = true
  }
}
`, template, data.RandomIntOfLength(10))
}

func (r SpringCloudGatewayResource) responseCachePerRoute(data acceptance.TestData, size, timeToLive string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id

  local_response_cache_per_route {
    size         = "%[2]s"
    time_to_live = "%[3]s"
  }
}
`, template, size, timeToLive)
}

func (r SpringCloudGatewayResource) responseCachePerInstance(data acceptance.TestData, size, timeToLive string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id

  local_response_cache_per_instance {
    size         = "%[2]s"
    time_to_live = "%[3]s"
  }
}
`, template, size, timeToLive)
}
