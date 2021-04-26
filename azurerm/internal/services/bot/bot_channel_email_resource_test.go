package bot_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BotChannelEmailResource struct {
}

func TestAccBotChannelEmail_basic(t *testing.T) {
	if ok := skipEmailChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_EMAIL`, AND `ARM_TEST_EMAIL_PASSWORD` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_email", "test")
	r := BotChannelEmailResource{}

	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
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
	if ok := skipSlackChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_SLACK_CLIENT_ID`, `ARM_TEST_SLACK_CLIENT_SECRET`, or `ARM_TEST_SLACK_VERIFICATION_TOKEN` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_email", "test")
	r := BotChannelEmailResource{}

	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("email_password"),
		{
			Config: r.basicUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("email_password"),
	})
}

func (t BotChannelEmailResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func skipEmailChannel() bool {
	if os.Getenv("ARM_TEST_EMAIL") == "" || os.Getenv("ARM_TEST_EMAIL_PASSWORD") == "" {
		return true
	}

	return false
}
