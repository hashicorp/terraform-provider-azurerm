package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMBotChannelMsTeams_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_ms_teams", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelMsTeamsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelMsTeams_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelMsTeamsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMBotChannelMsTeams_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_ms_teams", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelMsTeamsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelMsTeams_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelMsTeamsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBotChannelMsTeams_basicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelMsTeamsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBotChannelMsTeams_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelMsTeamsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMBotChannelMsTeamsExists(name string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for Bot Channel MsTeams")
		}

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Bot.ChannelClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bot_channel_ms_teams" {
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

func testAccAzureRMBotChannelMsTeams_basicConfig(data acceptance.TestData) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_ms_teams" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccAzureRMBotChannelMsTeams_basicUpdate(data acceptance.TestData) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_ms_teams" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  calling_web_hook    = "https://example.com/"
  enable_calling      = true
}
`, template)
}
