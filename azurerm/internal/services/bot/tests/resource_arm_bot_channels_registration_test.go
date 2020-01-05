package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBotChannelsRegistration(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being able provision against one app id at a time
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":    testAccAzureRMBotChannelsRegistration_basic,
			"update":   testAccAzureRMBotChannelsRegistration_update,
			"complete": testAccAzureRMBotChannelsRegistration_complete,
		},
		"connection": {
			"basic":    testAccAzureRMBotConnection_basic,
			"complete": testAccAzureRMBotConnection_complete,
		},
		"channel": {
			"slackBasic":    testAccAzureRMBotChannelSlack_basic,
			"slackUpdate":   testAccAzureRMBotChannelSlack_update,
			"msteamsBasic":  testAccAzureRMBotChannelMsTeams_basic,
			"msteamsUpdate": testAccAzureRMBotChannelMsTeams_update,
		},
		"web_app": {
			"basic":    testAccAzureRMBotWebApp_basic,
			"update":   testAccAzureRMBotWebApp_update,
			"complete": testAccAzureRMBotWebApp_complete,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMBotChannelsRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channels_registration", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelsRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelsRegistration_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelsRegistrationExists(data.ResourceName),
				),
			},
			data.ImportStep("developer_app_insights_api_key"),
		},
	})
}

func testAccAzureRMBotChannelsRegistration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channels_registration", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelsRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelsRegistration_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelsRegistrationExists(data.ResourceName),
				),
			},
			data.ImportStep("developer_app_insights_api_key"),
			{
				Config: testAccAzureRMBotChannelsRegistration_updateConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelsRegistrationExists(data.ResourceName),
				),
			},
			data.ImportStep("developer_app_insights_api_key"),
		},
	})
}

func testAccAzureRMBotChannelsRegistration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channels_registration", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBotChannelsRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBotChannelsRegistration_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotChannelsRegistrationExists(data.ResourceName),
				),
			},
			data.ImportStep("developer_app_insights_api_key"),
		},
	})
}

func testCheckAzureRMBotChannelsRegistrationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Bot Channels Registration: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Bot.BotClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on botClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Bot Channels Registration %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMBotChannelsRegistrationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Bot.BotClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bot" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bot Channels Registration still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMBotChannelsRegistration_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "F0"
  microsoft_app_id    = "${data.azurerm_client_config.current.service_principal_application_id}"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMBotChannelsRegistration_updateConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "F0"
  microsoft_app_id    = "${data.azurerm_client_config.current.service_principal_application_id}"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMBotChannelsRegistration_completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test" {
  name                    = "acctestappinsightsapikey-%d"
  application_insights_id = "${azurerm_application_insights.test.id}"
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_bot_channels_registration" "test" {
  name                = "acctestdf%d"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  microsoft_app_id    = "${data.azurerm_client_config.current.service_principal_application_id}"
  sku                 = "F0"

  endpoint                              = "https://example.com"
  developer_app_insights_api_key        = "${azurerm_application_insights_api_key.test.api_key}"
  developer_app_insights_application_id = "${azurerm_application_insights.test.app_id}"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
