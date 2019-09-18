package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMBotChannelSlack_basic(t *testing.T) {
	if ok := skipSlackChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_SLACK_CLIENT_ID`, `ARM_TEST_SLACK_CLIENT_SECRET`, or `ARM_TEST_SLACK_VERIFICATION_TOKEN` was not specified")
	}
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMBotChannelSlack_basicConfig(ri, testLocation())
	resourceName := "azurerm_bot_channel_slack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBotChannelSlackDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelSlackExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"verification_token",
					"landing_page_url",
				},
			},
		},
	})
}

func testAccAzureRMBotChannelSlack_update(t *testing.T) {
	if ok := skipSlackChannel(); ok {
		t.Skip("Skipping as one of `ARM_TEST_SLACK_CLIENT_ID`, `ARM_TEST_SLACK_CLIENT_SECRET`, or `ARM_TEST_SLACK_VERIFICATION_TOKEN` was not specified")
	}
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMBotChannelSlack_basicConfig(ri, testLocation())
	config2 := testAccAzureRMBotChannelSlack_basicUpdate(ri, testLocation())
	resourceName := "azurerm_bot_channel_slack.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBotChannelSlackDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelSlackExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"verification_token",
					"landing_page_url",
				},
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelSlackExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"verification_token",
					"landing_page_url",
				},
			},
		},
	})
}

func testCheckAzureRMBotChannelSlackExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Bot Channel Slack")
		}

		client := testAccProvider.Meta().(*ArmClient).bot.ChannelClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameSlackChannel))
		if err != nil {
			return fmt.Errorf("Bad: Get on botChannelClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Bot Channel Slack %q (resource group: %q / bot: %q) does not exist", name, resourceGroup, botName)
		}

		return nil
	}
}

func testCheckAzureRMBotChannelSlackDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).bot.ChannelClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bot_channel_slack" {
			continue
		}

		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameSlackChannel))

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bot Channel Slack still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMBotChannelSlack_basicConfig(rInt int, location string) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_slack" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  client_id           = "%s"
  client_secret       = "%s"
  verification_token  = "%s"
}
`, template, os.Getenv("ARM_TEST_SLACK_CLIENT_ID"), os.Getenv("ARM_TEST_SLACK_CLIENT_SECRET"), os.Getenv("ARM_TEST_SLACK_VERIFICATION_TOKEN"))
}

func testAccAzureRMBotChannelSlack_basicUpdate(rInt int, location string) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_slack" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  client_id           = "%s"
  client_secret       = "%s"
  verification_token  = "%s"
  landing_page_url    = "http://example.com"
}
`, template, os.Getenv("ARM_TEST_SLACK_CLIENT_ID"), os.Getenv("ARM_TEST_SLACK_CLIENT_SECRET"), os.Getenv("ARM_TEST_SLACK_VERIFICATION_TOKEN"))
}

func skipSlackChannel() bool {
	if os.Getenv("ARM_TEST_SLACK_CLIENT_ID") == "" || os.Getenv("ARM_TEST_SLACK_CLIENT_SECRET") == "" || os.Getenv("ARM_TEST_SLACK_VERIFICATION_TOKEN") == "" {
		return true
	}

	return false
}
