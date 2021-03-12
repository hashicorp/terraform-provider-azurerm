package bot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BotConnectionResource struct {
}

func testAccBotConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_connection", "test")
	r := BotConnectionResource{}

	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
	})
}

func testAccBotConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_connection", "test")
	r := BotConnectionResource{}

	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.completeConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
		{
			Config: r.completeUpdateConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
	})
}

func (t BotConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func (BotConnectionResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = azurerm_bot_channels_registration.test.name
  location              = azurerm_bot_channels_registration.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "box"
  client_id             = "test"
  client_secret         = "secret"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger)
}

func (BotConnectionResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = azurerm_bot_channels_registration.test.name
  location              = azurerm_bot_channels_registration.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Salesforce"
  client_id             = "test"
  client_secret         = "secret"
  scopes                = "testscope"

  parameters = {
    loginUri = "www.example.com"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger)
}

func (BotConnectionResource) completeUpdateConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = azurerm_bot_channels_registration.test.name
  location              = azurerm_bot_channels_registration.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Salesforce"
  client_id             = "test2"
  client_secret         = "secret2"
  scopes                = "testscope2"

  parameters = {
    loginUri = "www.example2.com"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger)
}
