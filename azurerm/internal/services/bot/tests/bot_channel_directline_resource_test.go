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

func testAccAzureRMBotChannelDirectline_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_directline", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelDirectlineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelDirectline_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelDirectlineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMBotChannelDirectline_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_directline", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelDirectlineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelDirectline_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelDirectlineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMBotChannelDirectline_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_directline", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelDirectlineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelDirectline_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelDirectlineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBotChannelDirectline_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelDirectlineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMBotChannelDirectline_basicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelDirectlineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMBotChannelDirectlineExists(name string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for Bot Channel Directline")
		}

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameDirectLineChannel1))
		if err != nil {
			return fmt.Errorf("Bad: Get on botChannelClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Bot Channel Directline %q (resource group: %q / bot: %q) does not exist", name, resourceGroup, botName)
		}

		return nil
	}
}

func testCheckAzureRMBotChannelDirectlineDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Bot.ChannelClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bot_channel_directline" {
			continue
		}

		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameDirectLineChannel1))
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bot Channel Directline still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMBotChannelDirectline_basicConfig(data acceptance.TestData) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(data)
	return fmt.Sprintf(` 
%s

resource "azurerm_bot_channel_directline" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  site {
    name    = "test"
    enabled = true
  }
}
`, template)
}

func testAccAzureRMBotChannelDirectline_completeConfig(data acceptance.TestData) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(data)
	return fmt.Sprintf(` 
%s

resource "azurerm_bot_channel_directline" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  site {
    name                            = "test"
    enabled                         = true
    v1_allowed                      = true
    v3_allowed                      = true
    enhanced_authentication_enabled = true
    trusted_origins                 = ["https://example.com"]
  }
}
`, template)
}

func testAccAzureRMBotChannelDirectline_basicUpdate(data acceptance.TestData) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(data)
	return fmt.Sprintf(` 
%s

resource "azurerm_bot_channel_directline" "test" {
  bot_name            = "${azurerm_bot_channels_registration.test.name}"
  location            = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  site {
    name    = "test"
    enabled = false
  }
}
`, template)
}
