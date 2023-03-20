package springcloud_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudGatewaySsoResource struct{}

func TestAccSpringCloudGatewaySso_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway_sso", "test")
	r := SpringCloudGatewaySsoResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_id", "client_secret"),
	})
}

func TestAccSpringCloudGatewaySso_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway_sso", "test")

	r := SpringCloudGatewaySsoResource{}
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

func (r SpringCloudGatewaySsoResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
	if resp.Properties == nil || resp.Properties.SsoProperties == nil {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

func (r SpringCloudGatewaySsoResource) template(data acceptance.TestData) string {
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

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudGatewaySsoResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {
}

resource "azurerm_spring_cloud_gateway_sso" "test" {
  spring_cloud_gateway_id = azurerm_spring_cloud_gateway.test.id
  client_id               = "%[2]s"
  client_secret           = "%[3]s"
  issuer_uri              = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
  scope                   = ["read"]
}
`, template, clientId, clientSecret)
}

func (r SpringCloudGatewaySsoResource) requiresImport(data acceptance.TestData) string {
	template := r.template(data)
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {
}

resource "azurerm_spring_cloud_gateway_sso" "test" {
  spring_cloud_gateway_id = azurerm_spring_cloud_gateway.test.id
  client_id               = "%[2]s"
  client_secret           = "%[3]s"
  issuer_uri              = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
  scope                   = ["read"]
}

resource "azurerm_spring_cloud_gateway_sso" "import" {
  spring_cloud_gateway_id = azurerm_spring_cloud_gateway.test.id
  client_id               = "%[2]s"
  client_secret           = "%[3]s"
  issuer_uri              = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
  scope                   = ["read"]
}
`, template, clientId, clientSecret)
}
