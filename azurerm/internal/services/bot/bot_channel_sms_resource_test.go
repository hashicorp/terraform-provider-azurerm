package bot_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/botservice/mgmt/2021-03-01/botservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BotChannelSMSResource struct {
}

func testAccBotChannelSMS_basic(t *testing.T) {
	skipSMSChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_sms", "test")
	r := BotChannelSMSResource{}

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

func testAccBotChannelSMS_requiresImport(t *testing.T) {
	skipSMSChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_sms", "test")
	r := BotChannelSMSResource{}

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

func (t BotChannelSMSResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSmsChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelSMSResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_sms" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_sid         = "%s"
  auth_token          = "%s"
  phone_number        = "%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_ACCOUNT_SID"), os.Getenv("ARM_TEST_AUTH_TOKEN"), os.Getenv("ARM_TEST_PHONE_NUMBER"))
}

func (r BotChannelSMSResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_sms" "import" {
  bot_name            = azurerm_bot_channel_sms.test.bot_name
  location            = azurerm_bot_channel_sms.test.location
  resource_group_name = azurerm_bot_channel_sms.test.resource_group_name
  account_sid         = azurerm_bot_channel_sms.test.account_sid
  auth_token          = azurerm_bot_channel_sms.test.auth_token
  phone_number        = azurerm_bot_channel_sms.test.phone_number
}
`, r.basic(data))
}

func skipSMSChannel(t *testing.T) {
	if os.Getenv("ARM_TEST_ACCOUNT_SID") == "" || os.Getenv("ARM_TEST_AUTH_TOKEN") == "" || os.Getenv("ARM_TEST_PHONE_NUMBER") == "" {
		t.Skip("Skipping as one of `ARM_TEST_ACCOUNT_SID`, `ARM_TEST_AUTH_TOKEN`, `ARM_TEST_PHONE_NUMBER` was not specified")
	}
}
