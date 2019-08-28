package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMMonitorActionGroup_basic(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(ri, location),
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

func TestAccAzureRMMonitorActionGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorActionGroup_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_app_service_custom_hostname_binding"),
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_emailReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_emailReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
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

func TestAccAzureRMMonitorActionGroup_smsReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_smsReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
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

func TestAccAzureRMMonitorActionGroup_webhookReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_webhookReceiver(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMMonitorActionGroup_complete(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMonitorActionGroup_complete(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.1.service_uri", "https://backup.example.com/warning"),
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

func TestAccAzureRMMonitorActionGroup_disabledUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMMonitorActionGroup_disabledBasic(ri, location)
	postConfig := testAccAzureRMMonitorActionGroup_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
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
		},
	})
}

func TestAccAzureRMMonitorActionGroup_singleReceiverUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	emailConfig := testAccAzureRMMonitorActionGroup_emailReceiver(ri, location)
	smsConfig := testAccAzureRMMonitorActionGroup_smsReceiver(ri, location)
	webhookConfig := testAccAzureRMMonitorActionGroup_webhookReceiver(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: emailConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
				),
			},
			{
				Config: smsConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
				),
			},
			{
				Config: webhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_multipleReceiversUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	basicConfig := testAccAzureRMMonitorActionGroup_basic(ri, location)
	completeConfig := testAccAzureRMMonitorActionGroup_complete(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
				),
			},
			{
				Config: completeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.1.email_address", "devops@contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.1.service_uri", "https://backup.example.com/warning"),
				),
			},
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook_receiver.#", "0"),
				),
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMMonitorActionGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_group" "import" {
  name                = "${azurerm_monitor_action_group.test.name}"
  resource_group_name = "${azurerm_monitor_action_group.test.resource_group_name}"
  short_name          = "${azurerm_monitor_action_group.test.short_name}"
}
`, template)
}

func testAccAzureRMMonitorActionGroup_emailReceiver(rInt int, location string) string {
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
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_smsReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_webhookReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  webhook_receiver {
    name        = "callmyapiaswell"
    service_uri = "http://example.com/alert"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_complete(rInt int, location string) string {
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
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_disabledBasic(rInt int, location string) string {
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
`, rInt, location, rInt)
}

func testCheckAzureRMMonitorActionGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).monitor.ActionGroupsClient
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

func testCheckAzureRMMonitorActionGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Action Group Instance: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).monitor.ActionGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorActionGroupsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Action Group Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}
