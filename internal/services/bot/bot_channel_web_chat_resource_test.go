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
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type BotChannelWebChatResource struct{}

func TestAccBotChannelWebChat_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

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

func TestAccBotChannelWebChat_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

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

func TestAccBotChannelWebChat_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

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

func TestAccBotChannelWebChat_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r BotChannelWebChatResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameWebChatChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelWebChatResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_web_chat" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name

  site {
    name = "TestSite"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}

func (r BotChannelWebChatResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_web_chat" "import" {
  bot_name            = azurerm_bot_channel_web_chat.test.bot_name
  location            = azurerm_bot_channel_web_chat.test.location
  resource_group_name = azurerm_bot_channel_web_chat.test.resource_group_name

  site {
    name = "TestSite"
  }
}
`, r.basic(data))
}

func (BotChannelWebChatResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_web_chat" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name

  site {
    name                        = "TestSite1"
    user_upload_enabled         = false
    endpoint_parameters_enabled = true
    storage_enabled             = false
  }

  site {
    name                        = "TestSite2"
    user_upload_enabled         = true
    endpoint_parameters_enabled = false
    storage_enabled             = true
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}
