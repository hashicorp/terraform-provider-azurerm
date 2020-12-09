package bot_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccBotChannelEmail_basic(t *testing.T) {
	if ok := skipEmailChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_EMAIL`, AND `ARM_TEST_EMAIL_PASSWORD` was not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_email", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBotChannelEmailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBotChannelEmail_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBotChannelEmailExists(data.ResourceName),
				),
			},
			data.ImportStep("email_password"),
		},
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBotChannelEmailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBotChannelEmail_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBotChannelEmailExists(data.ResourceName),
				),
			},
			data.ImportStep("email_password"),
			{
				Config: testAccBotChannelEmail_basicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBotChannelEmailExists(data.ResourceName),
				),
			},
			data.ImportStep("email_password"),
		},
	})
}

func testCheckBotChannelEmailExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Bot.ChannelClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Bot Channel Email")
		}

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameEmailChannel))
		if err != nil {
			return fmt.Errorf("Bad: Get on botChannelClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Bot Channel Email %q (resource group: %q / bot: %q) does not exist", name, resourceGroup, botName)
		}

		return nil
	}
}

func testCheckBotChannelEmailDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Bot.ChannelClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bot_channel_email" {
			continue
		}

		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameEmailChannel))
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bot Channel Email still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccBotChannelEmail_basicConfig(data acceptance.TestData) string {
	template := testAccBotChannelsRegistration_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_email" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  email_address       = "%s"
  email_password      = "%s"
}
`, template, os.Getenv("ARM_TEST_EMAIL"), os.Getenv("ARM_TEST_EMAIL_PASSWORD"))
}

func testAccBotChannelEmail_basicUpdate(data acceptance.TestData) string {
	template := testAccBotChannelsRegistration_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_email" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  email_address       = "%s"
  email_password      = "%s"
}
`, template, os.Getenv("ARM_TEST_EMAIL"), os.Getenv("ARM_TEST_EMAIL_PASSWORD"))
}

func skipEmailChannel() bool {
	if os.Getenv("ARM_TEST_EMAIL") == "" || os.Getenv("ARM_TEST_EMAIL_PASSWORD") == "" {
		return true
	}

	return false
}
