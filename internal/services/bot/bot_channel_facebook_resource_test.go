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

type BotChannelFacebookResource struct{}

func TestAccBotChannelFacebook_basic(t *testing.T) {
	skipFacebookChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_facebook", "test")
	r := BotChannelFacebookResource{}

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

func TestAccBotChannelFacebook_requiresImport(t *testing.T) {
	skipFacebookChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_facebook", "test")
	r := BotChannelFacebookResource{}

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

func TestAccBotChannelFacebook_update(t *testing.T) {
	skipFacebookChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_facebook", "test")
	r := BotChannelFacebookResource{}

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

func (t BotChannelFacebookResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameFacebookChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelFacebookResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_facebook" "test" {
  bot_name                    = azurerm_bot_channels_registration.test.name
  location                    = azurerm_bot_channels_registration.test.location
  resource_group_name         = azurerm_resource_group.test.name
  facebook_application_id     = "%s"
  facebook_application_secret = "%s"

  page {
    id           = "%s"
    access_token = "%s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_ID"), os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_SECRET"), os.Getenv("ARM_TEST_PAGE_ID"), os.Getenv("ARM_TEST_PAGE_ACCESS_TOKEN"))
}

func (r BotChannelFacebookResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_facebook" "import" {
  bot_name                    = azurerm_bot_channel_facebook.test.bot_name
  location                    = azurerm_bot_channel_facebook.test.location
  resource_group_name         = azurerm_bot_channel_facebook.test.resource_group_name
  facebook_application_id     = azurerm_bot_channel_facebook.test.facebook_application_id
  facebook_application_secret = azurerm_bot_channel_facebook.test.facebook_application_secret

  page {
    id           = "%s"
    access_token = "%s"
  }
}
`, r.basic(data), os.Getenv("ARM_TEST_PAGE_ID"), os.Getenv("ARM_TEST_PAGE_ACCESS_TOKEN"))
}

func (BotChannelFacebookResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_facebook" "test" {
  bot_name                    = azurerm_bot_channels_registration.test.name
  location                    = azurerm_bot_channels_registration.test.location
  resource_group_name         = azurerm_resource_group.test.name
  facebook_application_id     = "%s"
  facebook_application_secret = "%s"

  page {
    id           = "%s"
    access_token = "%s"
  }
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_ID2"), os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_SECRET2"), os.Getenv("ARM_TEST_PAGE_ID2"), os.Getenv("ARM_TEST_PAGE_ACCESS_TOKEN2"))
}

func skipFacebookChannel(t *testing.T) {
	if os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_ID") == "" || os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_SECRET") == "" || os.Getenv("ARM_TEST_PAGE_ID") == "" || os.Getenv("ARM_TEST_PAGE_ACCESS_TOKEN") == "" || os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_ID2") == "" || os.Getenv("ARM_TEST_FACEBOOK_APPLICATION_SECRET2") == "" || os.Getenv("ARM_TEST_PAGE_ID2") == "" || os.Getenv("ARM_TEST_PAGE_ACCESS_TOKEN2") == "" {
		t.Skip("Skipping as one of `ARM_TEST_FACEBOOK_APPLICATION_ID`, `ARM_TEST_FACEBOOK_APPLICATION_SECRET`, `ARM_TEST_PAGE_ID`, `ARM_TEST_PAGE_ACCESS_TOKEN`, `ARM_TEST_FACEBOOK_APPLICATION_ID2`, `ARM_TEST_FACEBOOK_APPLICATION_SECRET2`, `ARM_TEST_PAGE_ID2` and `ARM_TEST_PAGE_ACCESS_TOKEN2` was not specified")
	}
}
