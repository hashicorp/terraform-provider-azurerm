package springcloud_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudGatewayResource struct{}

func TestAccSpringCloudGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sso.0.client_id", "sso.0.client_secret", "sensitive_environment_variables.%", "sensitive_environment_variables.NEW_RELIC_APP_NAME"),
	})
}

func TestAccSpringCloudGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	r := SpringCloudGatewayResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sso.0.client_id", "sso.0.client_secret", "sensitive_environment_variables.%", "sensitive_environment_variables.NEW_RELIC_APP_NAME"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudGatewayResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudGatewayID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.GatewayClient.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudGatewayResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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
    credentials_allowed = false
    allowed_headers     = ["*"]
    allowed_methods     = ["PUT"]
    allowed_origins     = ["test.com"]
    exposed_headers     = ["x-test-header"]
    max_age_seconds     = 86400
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
