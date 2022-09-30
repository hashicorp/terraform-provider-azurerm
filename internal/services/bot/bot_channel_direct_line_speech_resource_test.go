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

type BotChannelDirectLineSpeechResource struct{}

func testAccBotChannelDirectLineSpeech_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cognitive_service_location", "cognitive_service_access_key"), // not returned from API
	})
}

func testAccBotChannelDirectLineSpeech_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

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

func testAccBotChannelDirectLineSpeech_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cognitive_service_location", "cognitive_service_access_key"), // not returned from API
	})
}

func testAccBotChannelDirectLineSpeech_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_direct_line_speech", "test")
	r := BotChannelDirectLineSpeechResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cognitive_service_location", "cognitive_service_access_key"), // not returned from API
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cognitive_service_location", "cognitive_service_access_key"), // not returned from API
	})
}

func (r BotChannelDirectLineSpeechResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameDirectLineSpeechChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelDirectLineSpeechResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "test" {
  name                = "acctest-cogacct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}

resource "azurerm_bot_channel_direct_line_speech" "test" {
  bot_name                     = azurerm_bot_channels_registration.test.name
  location                     = azurerm_bot_channels_registration.test.location
  resource_group_name          = azurerm_resource_group.test.name
  cognitive_service_location   = azurerm_cognitive_account.test.location
  cognitive_service_access_key = azurerm_cognitive_account.test.primary_access_key
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger)
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

func (BotChannelDirectLineSpeechResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "test" {
  name                = "acctest-cogacct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}

resource "azurerm_bot_channel_direct_line_speech" "test" {
  bot_name                     = azurerm_bot_channels_registration.test.name
  location                     = azurerm_bot_channels_registration.test.location
  resource_group_name          = azurerm_resource_group.test.name
  cognitive_service_location   = azurerm_cognitive_account.test.location
  cognitive_service_access_key = azurerm_cognitive_account.test.primary_access_key
  custom_speech_model_id       = "a9316355-7b04-4468-9f6e-114419e6c9cc"
  custom_voice_deployment_id   = "58dd86d4-31e3-4cf7-9b17-ee1d3dd77695"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger)
}

func (BotChannelDirectLineSpeechResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "test" {
  name                = "acctest-cogacct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-dls-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test2" {
  name                = "acctest-cogacct-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}

resource "azurerm_bot_channel_direct_line_speech" "test" {
  bot_name                     = azurerm_bot_channels_registration.test.name
  location                     = azurerm_bot_channels_registration.test.location
  resource_group_name          = azurerm_resource_group.test.name
  cognitive_service_location   = azurerm_cognitive_account.test2.location
  cognitive_service_access_key = azurerm_cognitive_account.test2.primary_access_key
  custom_speech_model_id       = "cf7a4202-9be3-4195-9619-5a747260626d"
  custom_voice_deployment_id   = "b815f623-c217-4327-b765-f6e0fd7dceef"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), data.RandomInteger, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}
