package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMActionGroup_basic(t *testing.T) {
	resourceName := "azurerm_action_group.test"
	ri := acctest.RandInt()
	config := testAccAzureRMActionGroup_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActionGroupExists(resourceName),
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
		},
	})
}

func TestAccAzureRMActionGroup_importBasic(t *testing.T) {
	resourceName := "azurerm_action_group.test"

	ri := acctest.RandInt()
	config := testAccAzureRMActionGroup_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMActionGroup_empty(t *testing.T) {
	resourceName := "azurerm_action_group.test"
	ri := acctest.RandInt()
	config := testAccAzureRMActionGroup_empty(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMActionGroup_importEmpty(t *testing.T) {
	resourceName := "azurerm_action_group.test"

	ri := acctest.RandInt()
	config := testAccAzureRMActionGroup_empty(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMActionGroup_disabledEmpty(t *testing.T) {
	resourceName := "azurerm_action_group.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMActionGroup_disabledEmpty(ri, location)
	postConfig := testAccAzureRMActionGroup_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMActionGroup_importDisabledEmpty(t *testing.T) {
	resourceName := "azurerm_action_group.test"

	ri := acctest.RandInt()
	config := testAccAzureRMActionGroup_disabledEmpty(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMActionGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_action_group" "test" {
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

func testAccAzureRMActionGroup_empty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_action_group" "test" {
  name                = "acctestActionGroup-%d"
  location            = "Global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}
`, rInt, location, rInt)
}

func testAccAzureRMActionGroup_disabledEmpty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_action_group" "test" {
  name                = "acctestActionGroup-%d"
  location            = "Global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
  enabled             = false
}
`, rInt, location, rInt)
}

func testCheckAzureRMActionGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).actionGroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_action_group" {
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

func testCheckAzureRMActionGroupExists(name string) resource.TestCheckFunc {
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
