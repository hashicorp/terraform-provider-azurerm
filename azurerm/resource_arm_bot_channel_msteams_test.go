package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMBotChannelMsTeams_basic(t *testing.T) {
	// if ok := skipMsTeamsChannel(); ok {
	// 	t.Skip("Skipping as one of `ARM_TEST_MsTeams_CLIENT_ID`, `ARM_TEST_MsTeams_CLIENT_SECRET`, or `ARM_TEST_MsTeams_VERIFICATION_TOKEN` was not specified")
	// }
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMBotChannelMsTeams_basicConfig(ri, testLocation())
	resourceName := "azurerm_bot_channel_msteams.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBotChannelMsTeamsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelMsTeamsExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enabled",
					"calling_web_hook",
					"enable_calling",
				},
			},
		},
	})
}

func testAccAzureRMBotChannelMsTeams_update(t *testing.T) {
	// if ok := skipMsTeamsChannel(); ok {
	// 	t.Skip("Skipping as one of `ARM_TEST_MsTeams_CLIENT_ID`, `ARM_TEST_MsTeams_CLIENT_SECRET`, or `ARM_TEST_MsTeams_VERIFICATION_TOKEN` was not specified")
	// }
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMBotChannelMsTeams_basicConfig(ri, testLocation())
	config2 := testAccAzureRMBotChannelMsTeams_basicUpdate(ri, testLocation())
	resourceName := "azurerm_bot_channel_msteams.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBotChannelMsTeamsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelMsTeamsExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enabled",
					"calling_web_hook",
					"enable_calling",
				},
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelMsTeamsExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enabled",
					"calling_web_hook",
					"enable_calling",
				},
			},
		},
	})
}

func testCheckAzureRMBotChannelMsTeamsExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Bot Channel MsTeams")
		}

		client := testAccProvider.Meta().(*ArmClient).Bot.ChannelClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameMsTeamsChannel))
		if err != nil {
			return fmt.Errorf("Bad: Get on botChannelClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Bot Channel MsTeams %q (resource group: %q / bot: %q) does not exist", name, resourceGroup, botName)
		}

		return nil
	}
}

func testCheckAzureRMBotChannelMsTeamsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Bot.ChannelClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bot_channel_msteams" {
			continue
		}

		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameMsTeamsChannel))

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bot Channel MsTeams still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMBotChannelMsTeams_basicConfig(rInt int, location string) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_msteams" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  calling_web_hook    = "http://example.com"
  enable_calling       = true
}
`, template)
}

func testAccAzureRMBotChannelMsTeams_basicUpdate(rInt int, location string) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_msteams" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  enabled 			  = false
  calling_web_hook    = "http://example2.com"
  enable_calling      = false
}
`, template)
}

// func skipMsTeamsChannel() bool {
// 	if os.Getenv("ARM_TEST_MsTeams_CLIENT_ID") == "" || os.Getenv("ARM_TEST_MsTeams_CLIENT_SECRET") == "" || os.Getenv("ARM_TEST_MsTeams_VERIFICATION_TOKEN") == "" {
// 		return true
// 	}

// 	return false
// }
