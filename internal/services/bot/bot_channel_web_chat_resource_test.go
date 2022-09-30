package bot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2021-05-01-preview/botservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BotChannelWebChatResource struct{}

func testAccBotChannelWebChat_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccBotChannelWebChat_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccBotChannelWebChat_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccBotChannelWebChat_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_web_chat", "test")
	r := BotChannelWebChatResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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
  site_names          = ["TestSite"]
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
  site_names          = ["TestSite"]
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
  site_names          = ["TestSite2", "TestSite3"]
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}
