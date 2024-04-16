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

type BotChannelEmailResource struct{}

func TestAccBotChannelEmail_basic(t *testing.T) {
	if ok := skipEmailChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_EMAIL`, AND `ARM_TEST_EMAIL_PASSWORD` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_email", "test")
	r := BotChannelEmailResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("email_password"),
	})
}

func TestAccBotChannelEmail_update(t *testing.T) {
	if ok := skipEmailChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_EMAIL`, AND `ARM_TEST_EMAIL_PASSWORD` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_email", "test")
	r := BotChannelEmailResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("email_password"),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("email_password"),
	})
}

func TestAccBotChannelEmail_magicCode(t *testing.T) {
	if os.Getenv("ARM_TEST_BOT_RESOURCE_GROUP_NAME") == "" || os.Getenv("ARM_TEST_BOT_NAME") == "" || os.Getenv("ARM_TEST_EMAIL") == "" || os.Getenv("ARM_TEST_MAGIC_CODE") == "" {
		t.Skip("Skipping as one of `ARM_TEST_BOT_RESOURCE_GROUP_NAME`, `ARM_TEST_BOT_LOCATION`, `ARM_TEST_BOT_NAME`, `ARM_TEST_EMAIL`, AND `ARM_TEST_MAGIC_CODE` was not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_email", "test")
	r := BotChannelEmailResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.magicCode(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("email_password", "magic_code"),
	})
}

func (t BotChannelEmailResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameEmailChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelEmailResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_email" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  email_address       = "%s"
  email_password      = "%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_EMAIL"), os.Getenv("ARM_TEST_EMAIL_PASSWORD"))
}

func (BotChannelEmailResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_email" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  email_address       = "%s"
  email_password      = "%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_EMAIL"), os.Getenv("ARM_TEST_EMAIL_PASSWORD"))
}

func (BotChannelEmailResource) magicCode() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_bot_channel_email" "test" {
  bot_name            = "%s"
  location            = "global"
  resource_group_name = "%s"
  email_address       = "%s"
  magic_code          = "%s"
}
`, os.Getenv("ARM_TEST_BOT_NAME"), os.Getenv("ARM_TEST_BOT_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_EMAIL"), os.Getenv("ARM_TEST_MAGIC_CODE"))
}

func skipEmailChannel() bool {
	if os.Getenv("ARM_TEST_EMAIL") == "" || os.Getenv("ARM_TEST_EMAIL_PASSWORD") == "" {
		return true
	}

	return false
}
