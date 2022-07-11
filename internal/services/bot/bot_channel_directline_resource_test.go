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

type BotChannelDirectlineResource struct{}

func testAccBotChannelDirectline_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_directline", "test")
	r := BotChannelDirectlineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccBotChannelDirectline_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_directline", "test")
	r := BotChannelDirectlineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccBotChannelDirectline_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_directline", "test")
	r := BotChannelDirectlineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BotChannelDirectlineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameBasicChannelChannelNameDirectLineChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelDirectlineResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(` 
%s

resource "azurerm_bot_channel_directline" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  site {
    name    = "test"
    enabled = true
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}

func (r BotChannelDirectlineResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(` 
%s

resource "azurerm_bot_channel_directline" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  site {
    name                            = "test"
    enabled                         = true
    v1_allowed                      = true
    v3_allowed                      = true
    enhanced_authentication_enabled = true
    trusted_origins                 = ["https://example.com"]
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}

func (r BotChannelDirectlineResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(` 
%s

resource "azurerm_bot_channel_directline" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  site {
    name    = "test"
    enabled = false
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data))
}
