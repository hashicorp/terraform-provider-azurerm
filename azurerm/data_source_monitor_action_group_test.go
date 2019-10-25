package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceArmMonitorActionGroup_basic(t *testing.T) {
	dataSourceName := "data.azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccDataSourceArmMonitorActionGroup_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "short_name", "acctestag"),
					resource.TestCheckResourceAttr(dataSourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceArmMonitorActionGroup_disabledBasic(t *testing.T) {
	dataSourceName := "data.azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccDataSourceArmMonitorActionGroup_disabledBasic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "short_name", "acctestag"),
					resource.TestCheckResourceAttr(dataSourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceArmMonitorActionGroup_complete(t *testing.T) {
	dataSourceName := "data.azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccDataSourceArmMonitorActionGroup_complete(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "email_receiver.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(dataSourceName, "email_receiver.1.email_address", "devops@contoso.com"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.1.service_uri", "https://backup.example.com/warning"),
				),
			},
		},
	})
}

func testAccDataSourceArmMonitorActionGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}

data "azurerm_monitor_action_group" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_monitor_action_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceArmMonitorActionGroup_disabledBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
  enabled             = false
}

data "azurerm_monitor_action_group" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_monitor_action_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceArmMonitorActionGroup_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
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

  sms_receiver {
    name         = "remotesupport"
    country_code = "86"
    phone_number = "13888888888"
  }

  webhook_receiver {
    name        = "callmyapiaswell"
    service_uri = "http://example.com/alert"
  }

  webhook_receiver {
    name        = "callmybackupapi"
    service_uri = "https://backup.example.com/warning"
  }
}

data "azurerm_monitor_action_group" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_monitor_action_group.test.name}"
}
`, rInt, location, rInt)
}
