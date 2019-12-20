package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMMonitorActionGroup_basic(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ExpectError: acceptance.RequiresImportError("azurerm_app_service_custom_hostname_binding"),
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_emailReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_emailReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_itsmReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_itsmReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_azureAppPushReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_azureAppPushReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_smsReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_webhookReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_automationRunbookReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_automationRunbookReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_voiceReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_voiceReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_logicAppReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_logicAppReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_azureFunctionReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_azureFunctionReceiver(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_armRoleReceiver(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_armRoleReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_complete(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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
	location := acceptance.Location()
	preConfig := testAccAzureRMMonitorActionGroup_disabledBasic(ri, location)
	postConfig := testAccAzureRMMonitorActionGroup_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccAzureRMMonitorActionGroup_singleReceiverUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_emailReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_itsmReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_azureAppPushReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_smsReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_webhookReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_automationRunbookReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_voiceReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_logicAppReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_azureFunctionReceiver(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_armRoleReceiver(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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

func TestAccAzureRMMonitorActionGroup_multipleReceiversUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_action_group.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_complete(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMonitorActionGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(resourceName),
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
		use_common_alert_schema = false
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_itsmReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  itsm_receiver {
		name          = "createorupdateticket"
		workspace_id = "6eee3a18-aac3-40e4-b98e-1f309f329816"
		connection_id = "53de6956-42b4-41ba-be3c-b154cdf17b13"
		ticket_configuration = "{}"
		region = "eastus"
	}
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_azureAppPushReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
	short_name          = "acctestag"

	azure_app_push_receiver {
		name          = "pushtoadmin"
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
    name                    = "callmyapiaswell"
    service_uri             = "http://example.com/alert"
    use_common_alert_schema = true
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_automationRunbookReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  automation_runbook_receiver {
		name = "action_name_1"
		automation_account_id = "${azurerm_automation_account.test.id}"
  	runbook_name = "${azurerm_automation_runbook.test.name}"
		webhook_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
		is_global_runbook = true
		service_uri = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
		use_common_alert_schema = false
	}
}

resource "azurerm_automation_account" "test" {
	name                = "acctestAA-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"

	sku {
		name = "Basic"
	}
}

resource "azurerm_automation_runbook" "test" {
	name                = "Get-AzureVMTutorial"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	account_name        = "${azurerm_automation_account.test.name}"
	log_verbose         = "true"
	log_progress        = "true"
	description         = "This is an test runbook"
	runbook_type        = "PowerShellWorkflow"

	publish_content_link {
		uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
	}
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMonitorActionGroup_voiceReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

	voice_receiver {
		name         = "oncallmsg"
		country_code = "1"
		phone_number = "1231231234"
	}
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_logicAppReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
	short_name          = "acctestag"

	logic_app_receiver {
		name = "logicappaction"
		resource_id = "${azurerm_logic_app_workflow.test.id}"
		callback_url = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
		use_common_alert_schema = true
	}
}

resource "azurerm_logic_app_workflow" "test" {
	name                = "acctestLA-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_logic_app_trigger_http_request" "test" {
	name         = "some-http-trigger"
	logic_app_id = "${azurerm_logic_app_workflow.test.id}"

schema = <<SCHEMA
{
	"type": "object",
	"properties": {
		"hello": {
			"type": "string"
		}
	}
}
SCHEMA
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMonitorActionGroup_azureFunctionReceiver(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

	azure_function_receiver {
		name = "funcaction"
		function_app_resource_id = "${azurerm_function_app.test.id}"
		function_name = "myfunc"
		http_trigger_url = "https://example.com/trigger"
		use_common_alert_schema = true
	}
}

resource "azurerm_storage_account" "test" {
	name                     = "acctestsa%s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
	name                = "acctestSP-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"

	sku {
		tier = "Standard"
		size = "S1"
	}
}

resource "azurerm_function_app" "test" {
	name                      = "acctestFA-%d"
	location                  = "${azurerm_resource_group.test.location}"
	resource_group_name       = "${azurerm_resource_group.test.name}"
	app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
	storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, rInt, rString, rInt, rInt)
}

func testAccAzureRMMonitorActionGroup_armRoleReceiver(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

	arm_role_receiver {
		name = "Monitoring Reader"
		role_id = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
		use_common_alert_schema = false
	}
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorActionGroup_complete(rInt int, rString, location string) string {
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
			name          		= "sendtodevops"
			email_address 		= "devops@contoso.com"
  			use_common_alert_schema = true
		}

		itsm_receiver {
			name          = "createorupdateticket"
			workspace_id = "6eee3a18-aac3-40e4-b98e-1f309f329816"
			connection_id = "53de6956-42b4-41ba-be3c-b154cdf17b13"
			ticket_configuration = "{}"
			region = "eastus"
		}

		azure_app_push_receiver {
			name          = "pushtoadmin"
			email_address = "admin@contoso.com"
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
			name                    = "callmybackupapi"
			service_uri             = "https://backup.example.com/warning"
			use_common_alert_schema = true
		}

		automation_runbook_receiver {
			name = "action_name_1"
			automation_account_id = "${azurerm_automation_account.test.id}"
    	runbook_name = "${azurerm_automation_runbook.test.name}"
			webhook_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
			is_global_runbook = true
			service_uri = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
			use_common_alert_schema = false
		}

		voice_receiver {
			name         = "oncallmsg"
			country_code = "1"
			phone_number = "1231231234"
		}

		logic_app_receiver {
			name = "logicappaction"
			resource_id = "${azurerm_logic_app_workflow.test.id}"
			callback_url = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
			use_common_alert_schema = false
		}

		azure_function_receiver {
			name = "funcaction"
			function_app_resource_id = "${azurerm_function_app.test.id}"
			function_name = "myfunc"
			http_trigger_url = "https://example.com/trigger"
			use_common_alert_schema = false
		}

		arm_role_receiver {
			name = "Monitoring Reader"
			role_id = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
			use_common_alert_schema = false
		}
	}

	resource "azurerm_automation_account" "test" {
		name                = "acctestAA-%d"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"

		sku {
		  name = "Basic"
		}
	}

	resource "azurerm_automation_runbook" "test" {
		name                = "Get-AzureVMTutorial"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"
		account_name        = "${azurerm_automation_account.test.name}"
		log_verbose         = "true"
		log_progress        = "true"
		description         = "This is an test runbook"
		runbook_type        = "PowerShellWorkflow"

		publish_content_link {
		  uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
		}
	}

	resource "azurerm_logic_app_workflow" "test" {
		name                = "acctestLA-%d"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"
	}

	resource "azurerm_logic_app_trigger_http_request" "test" {
		name         = "some-http-trigger"
		logic_app_id = "${azurerm_logic_app_workflow.test.id}"

	schema = <<SCHEMA
	{
		"type": "object",
		"properties": {
			"hello": {
				"type": "string"
			}
		}
	}
	SCHEMA
	}

	resource "azurerm_storage_account" "test" {
		name                     = "acctestsa%s"
		resource_group_name      = "${azurerm_resource_group.test.name}"
		location                 = "${azurerm_resource_group.test.location}"
		account_tier             = "Standard"
		account_replication_type = "LRS"
	}

	resource "azurerm_app_service_plan" "test" {
		name                = "acctestSP-%d"
		location            = "${azurerm_resource_group.test.location}"
		resource_group_name = "${azurerm_resource_group.test.name}"

		sku {
			tier = "Standard"
			size = "S1"
		}
	}

	resource "azurerm_function_app" "test" {
		name                      = "acctestFA-%d"
		location                  = "${azurerm_resource_group.test.location}"
		resource_group_name       = "${azurerm_resource_group.test.name}"
		app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
		storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
	}
`, rInt, location, rInt, rInt, rInt, rString, rInt, rInt)
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
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
