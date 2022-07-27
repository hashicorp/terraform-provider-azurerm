package bot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BotConnectionResource struct{}

func testAccBotConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_connection", "test")
	r := BotConnectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
	})
}

func testAccBotConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_connection", "test")
	r := BotConnectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
		{
			Config: r.completeUpdateConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
	})
}

func (t BotConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ConnectionClient.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r BotConnectionResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = azurerm_bot_channels_registration.test.name
  location              = azurerm_bot_channels_registration.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "box"
  client_id             = data.azurerm_client_config.current.client_id
  client_secret         = "86546868-e7ed-429f-b0e5-3a1caea7db64"
}
`, r.template(data), data.RandomInteger)
}

func (r BotConnectionResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = azurerm_bot_channels_registration.test.name
  location              = azurerm_bot_channels_registration.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Salesforce"
  client_id             = data.azurerm_client_config.current.client_id
  client_secret         = "60a97b1d-0894-4c5a-9968-7d1d29d77aed"
  scopes                = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"

  parameters = {
    loginUri = "https://www.google.com"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BotConnectionResource) completeUpdateConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = azurerm_bot_channels_registration.test.name
  location              = azurerm_bot_channels_registration.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Salesforce"
  client_id             = azurerm_user_assigned_identity.test.client_id
  client_secret         = "32ea21cb-cb20-4df9-ad39-b55e985e9117"
  scopes                = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}"

  parameters = {
    loginUri = "https://www.terraform.io"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r BotConnectionResource) template(data acceptance.TestData) string {
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

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "F0"
  microsoft_app_id    = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
