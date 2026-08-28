// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type BotChannelDirectLineSpeechResource struct{}

func TestAccBotChannelDirectLineSpeech_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cognitive_service_location", "cognitive_service_access_key"), // not returned from API
	})
}

func TestAccBotChannelDirectLineSpeech_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

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

func TestAccBotChannelDirectLineSpeech_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

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

func TestAccBotChannelDirectLineSpeech_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAccessKeyComplete(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Test that the CustomizeDiff functions as expected
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionReplace),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cognitive_service_location", "cognitive_service_access_key"), // not returned from API
		{
			Config: r.complete(data),
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

func (r BotChannelDirectLineSpeechResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameDirectLineSpeechChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return pointer.To(resp.Properties != nil), nil
}

func (r BotChannelDirectLineSpeechResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_direct_line_speech" "test" {
  bot_name                     = azurerm_bot_channels_registration.test.name
  location                     = azurerm_bot_channels_registration.test.location
  resource_group_name          = azurerm_resource_group.test.name
  cognitive_service_location   = azurerm_cognitive_account.test.location
  cognitive_service_access_key = azurerm_cognitive_account.test.primary_access_key
}
`, r.template(data))
}

func (r BotChannelDirectLineSpeechResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_direct_line_speech" "import" {
  bot_name                     = azurerm_bot_channel_direct_line_speech.test.bot_name
  location                     = azurerm_bot_channel_direct_line_speech.test.location
  resource_group_name          = azurerm_bot_channel_direct_line_speech.test.resource_group_name
  cognitive_service_location   = azurerm_bot_channel_direct_line_speech.test.cognitive_service_location
  cognitive_service_access_key = azurerm_bot_channel_direct_line_speech.test.cognitive_service_access_key
}
`, r.basic(data))
}

func (r BotChannelDirectLineSpeechResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_direct_line_speech" "test" {
  bot_name                   = azurerm_bot_channels_registration.test.name
  location                   = azurerm_bot_channels_registration.test.location
  resource_group_name        = azurerm_resource_group.test.name
  cognitive_account_id       = azurerm_cognitive_account.test.id
  custom_speech_model_id     = "a9316355-7b04-4468-9f6e-114419e6c9cc"
  custom_voice_deployment_id = "58dd86d4-31e3-4cf7-9b17-ee1d3dd77695"
}
`, r.template(data))
}

func (r BotChannelDirectLineSpeechResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account" "test2" {
  name                = "acctest-cogacct2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}

resource "azurerm_bot_channel_direct_line_speech" "test" {
  bot_name                   = azurerm_bot_channels_registration.test.name
  location                   = azurerm_bot_channels_registration.test.location
  resource_group_name        = azurerm_resource_group.test.name
  cognitive_account_id       = azurerm_cognitive_account.test2.id
  custom_speech_model_id     = "cf7a4202-9be3-4195-9619-5a747260626d"
  custom_voice_deployment_id = "b815f623-c217-4327-b765-f6e0fd7dceef"
}
`, r.template(data), data.RandomInteger)
}

func TestAccBotChannelDirectLineSpeech_withAccessKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAccessKeyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cognitive_service_location", "cognitive_service_access_key"), // not returned from API
	})
}

func (r BotChannelDirectLineSpeechResource) withAccessKeyComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_direct_line_speech" "test" {
  bot_name                     = azurerm_bot_channels_registration.test.name
  location                     = azurerm_bot_channels_registration.test.location
  resource_group_name          = azurerm_resource_group.test.name
  cognitive_service_location   = azurerm_cognitive_account.test.location
  cognitive_service_access_key = azurerm_cognitive_account.test.primary_access_key
  custom_speech_model_id       = "a9316355-7b04-4468-9f6e-114419e6c9cc"
  custom_voice_deployment_id   = "58dd86d4-31e3-4cf7-9b17-ee1d3dd77695"
}
`, r.template(data))
}

func (BotChannelDirectLineSpeechResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account" "test" {
  name                = "acctest-cogacct-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger)
}
