package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMBotConnection_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMBotConnection_basicConfig(ri, testLocation())
	resourceName := "azurerm_bot_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBotConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"service_provider_name",
				},
			},
		},
	})
}

func testAccAzureRMBotConnection_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMBotConnection_completeConfig(ri, testLocation())
	config2 := testAccAzureRMBotConnection_completeUpdateConfig(ri, testLocation())
	resourceName := "azurerm_bot_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBotConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"service_provider_name",
				},
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBotConnectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"service_provider_name",
				},
			},
		},
	})
}

func testCheckAzureRMBotConnectionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Bot Channels Registration: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).bot.ConnectionClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, botName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on botConnectionClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Bot Connection %q (resource group: %q / bot: %q) does not exist", name, resourceGroup, botName)
		}

		return nil
	}
}

func testCheckAzureRMBotConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).bot.ConnectionClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bot" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		botName := rs.Primary.Attributes["bot_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, botName, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bot Connection still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMBotConnection_basicConfig(rInt int, location string) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = "${azurerm_bot_channels_registration.test.name}"
  location              = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  service_provider_name = "box"
  client_id             = "test"
  client_secret         = "secret"
}
`, template, rInt)
}

func testAccAzureRMBotConnection_completeConfig(rInt int, location string) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = "${azurerm_bot_channels_registration.test.name}"
  location              = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  service_provider_name = "Salesforce"
  client_id             = "test"
  client_secret         = "secret"
  scopes                = "testscope"
  
  parameters = {
    loginUri = "www.example.com"
  }
}
`, template, rInt)
}

func testAccAzureRMBotConnection_completeUpdateConfig(rInt int, location string) string {
	template := testAccAzureRMBotChannelsRegistration_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_bot_connection" "test" {
  name                  = "acctestBc%d"
  bot_name              = "${azurerm_bot_channels_registration.test.name}"
  location              = "${azurerm_bot_channels_registration.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  service_provider_name = "Salesforce"
  client_id             = "test2"
  client_secret         = "secret2"
  scopes                = "testscope2"
  
  parameters = {
    loginUri = "www.example2.com"
  }
}
`, template, rInt)
}
