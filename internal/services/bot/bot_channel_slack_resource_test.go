// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type BotChannelSlackResource struct{}

func TestAccBotChannelSlack_basic(t *testing.T) {
	if ok := skipSlackChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_SLACK_CLIENT_ID`, `ARM_TEST_SLACK_CLIENT_SECRET`, or `ARM_TEST_SLACK_VERIFICATION_TOKEN` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_slack", "test")
	r := BotChannelSlackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "verification_token", "landing_page_url", "client_id"),
	})
}

func TestAccBotChannelSlack_complete(t *testing.T) {
	if ok := skipSlackChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_SLACK_CLIENT_ID`, `ARM_TEST_SLACK_CLIENT_SECRET`, `ARM_TEST_SLACK_VERIFICATION_TOKEN`, `ARM_TEST_SLACK_SIGNING_SECRET` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_slack", "test")
	r := BotChannelSlackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "verification_token", "landing_page_url", "signing_secret", "client_id"),
	})
}

func TestAccBotChannelSlack_update(t *testing.T) {
	if ok := skipSlackChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_SLACK_CLIENT_ID`, `ARM_TEST_SLACK_CLIENT_SECRET`, or `ARM_TEST_SLACK_VERIFICATION_TOKEN` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_slack", "test")
	r := BotChannelSlackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "verification_token", "landing_page_url", "client_id"),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret", "verification_token", "landing_page_url", "client_id"),
	})
}

func (t BotChannelSlackResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelSlackResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_slack" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "%s"
  client_secret       = "%s"
  verification_token  = "%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_SLACK_CLIENT_ID"), os.Getenv("ARM_TEST_SLACK_CLIENT_SECRET"), os.Getenv("ARM_TEST_SLACK_VERIFICATION_TOKEN"))
}

func (BotChannelSlackResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_slack" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "%s"
  client_secret       = "%s"
  verification_token  = "%s"
  landing_page_url    = "http://example.com"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_SLACK_CLIENT_ID"), os.Getenv("ARM_TEST_SLACK_CLIENT_SECRET"), os.Getenv("ARM_TEST_SLACK_VERIFICATION_TOKEN"))
}

func (BotChannelSlackResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_slack" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "%s"
  client_secret       = "%s"
  verification_token  = "%s"
  signing_secret      = "%s"
  landing_page_url    = "http://example.com"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_SLACK_CLIENT_ID"), os.Getenv("ARM_TEST_SLACK_CLIENT_SECRET"), os.Getenv("ARM_TEST_SLACK_VERIFICATION_TOKEN"), os.Getenv("ARM_TEST_SLACK_SIGNING_SECRET"))
}

func skipSlackChannel() bool {
	if os.Getenv("ARM_TEST_SLACK_CLIENT_ID") == "" || os.Getenv("ARM_TEST_SLACK_CLIENT_SECRET") == "" || os.Getenv("ARM_TEST_SLACK_VERIFICATION_TOKEN") == "" || os.Getenv("ARM_TEST_SLACK_SIGNING_SECRET") == "" {
		return true
	}

	return false
}
