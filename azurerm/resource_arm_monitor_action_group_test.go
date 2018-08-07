package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMMonitorActionGroup_basic(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMonitorActionGroup_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.1.email_address", "devops@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_empty(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMonitorActionGroup_empty(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_disabledEmpty(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMMonitorActionGroup_disabledEmpty(ri, location)
	postConfig := testAccAzureRMMonitorActionGroup_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMMonitorActionGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  location            = "Global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  email_receiver {
    name          = "sendtoadmin"
    email_address = "admin@contoso.com"
  }

  email_receiver {
    name          = "sendtodevops"
    email_address = "devops@contoso.com"
  }

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }

  webhook_receiver {
    name        = "callmyapiaswell"
    service_uri = "http://example.com/alert"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_empty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  location            = "Global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_disabledEmpty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  location            = "Global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
  enabled             = false
}
`, rInt, location, rInt)
}

func testCheckAzureRMMonitorActionGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).actionGroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_action_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Action Group still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMMonitorActionGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Action Group Instance: %s", resourceName)
		}

		conn := testAccProvider.Meta().(*ArmClient).actionGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, resourceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on actionGroupsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Action Group Instance %q (resource group: %q) does not exist", resourceName, resourceGroup)
		}

		return nil
	}
}
