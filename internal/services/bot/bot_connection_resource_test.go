// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func TestAccBotConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_connection", "test")
	r := BotConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
	})
}

func TestAccBotConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_connection", "test")
	r := BotConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateConfig(data, "/subscriptions/${data.azurerm_client_config.current.subscription_id}", "https://www.google.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "service_provider_name"),
		{
			Config: r.updateConfig(data, "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}", "https://www.terraform.io"),
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
  client_id             = azuread_application_registration.test.client_id
  client_secret         = "86546868-e7ed-429f-b0e5-3a1caea7db64"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger)
}

func (r BotConnectionResource) updateConfig(data acceptance.TestData, scope string, loginUrl string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = azurerm_bot_channels_registration.test.name
  location              = azurerm_bot_channels_registration.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Salesforce"
  client_id             = azuread_application_registration.test.client_id
  client_secret         = "60a97b1d-0894-4c5a-9968-7d1d29d77aed"
  scopes                = "%s"

  parameters = {
    loginUri = "%s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger, scope, loginUrl)
}
