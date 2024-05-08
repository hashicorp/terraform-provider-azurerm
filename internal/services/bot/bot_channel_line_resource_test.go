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

type BotChannelLineResource struct{}

func TestAccBotChannelLine_basic(t *testing.T) {
	skipLineChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_line", "test")
	r := BotChannelLineResource{}

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

func TestAccBotChannelLine_requiresImport(t *testing.T) {
	skipLineChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_line", "test")
	r := BotChannelLineResource{}

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

func TestAccBotChannelLine_complete(t *testing.T) {
	skipLineChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_line", "test")
	r := BotChannelLineResource{}

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

func TestAccBotChannelLine_update(t *testing.T) {
	skipLineChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_line", "test")
	r := BotChannelLineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BotChannelLineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameLineChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelLineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_line" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name

  line_channel {
    access_token = "%s"
    secret       = "%s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN"), os.Getenv("ARM_TEST_CHANNEL_SECRET"))
}

func (r BotChannelLineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_line" "import" {
  bot_name            = azurerm_bot_channel_line.test.bot_name
  location            = azurerm_bot_channel_line.test.location
  resource_group_name = azurerm_bot_channel_line.test.resource_group_name

  line_channel {
    access_token = "%s"
    secret       = "%s"
  }
}
`, r.basic(data), os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN"), os.Getenv("ARM_TEST_CHANNEL_SECRET"))
}

func (BotChannelLineResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_line" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name

  line_channel {
    access_token = "%s"
    secret       = "%s"
  }

  line_channel {
    access_token = "%s"
    secret       = "%s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN"), os.Getenv("ARM_TEST_CHANNEL_SECRET"), os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN2"), os.Getenv("ARM_TEST_CHANNEL_SECRET2"))
}

func (BotChannelLineResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_line" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name

  line_channel {
    access_token = "%s"
    secret       = "%s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN2"), os.Getenv("ARM_TEST_CHANNEL_SECRET2"))
}

func skipLineChannel(t *testing.T) {
	if os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN") == "" || os.Getenv("ARM_TEST_CHANNEL_SECRET") == "" || os.Getenv("ARM_TEST_CHANNEL_ACCESS_TOKEN2") == "" || os.Getenv("ARM_TEST_CHANNEL_SECRET2") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CHANNEL_ACCESS_TOKEN`, `ARM_TEST_CHANNEL_SECRET`, `ARM_TEST_CHANNEL_ACCESS_TOKEN2`, `ARM_TEST_CHANNEL_SECRET2` was not specified")
	}
}
