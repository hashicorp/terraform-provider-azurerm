package monitor_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMonitorActionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorActionGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_monitor_action_group"),
			},
		},
	})
}

func TestAccAzureRMMonitorActionGroup_emailReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_emailReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_itsmReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_itsmReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_azureAppPushReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_azureAppPushReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_smsReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_smsReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_webhookReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_webhookReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_automationRunbookReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_automationRunbookReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_voiceReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_voiceReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_logicAppReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_logicAppReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_azureFunctionReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_azureFunctionReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_armRoleReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_armRoleReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_disabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_disabledBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_disabledBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_singleReceiverUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_emailReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_itsmReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_azureAppPushReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_smsReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_webhookReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_automationRunbookReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_voiceReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_logicAppReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_azureFunctionReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_armRoleReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_multipleReceiversUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionGroup_secureWebhookReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionGroup_webhookReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionGroup_secureWebhookReceiver(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMMonitorActionGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActionGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_group" "import" {
  name                = azurerm_monitor_action_group.test.name
  resource_group_name = azurerm_monitor_action_group.test.resource_group_name
  short_name          = azurerm_monitor_action_group.test.short_name
}
`, template)
}

func testAccAzureRMMonitorActionGroup_emailReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  email_receiver {
    name                    = "sendtoadmin"
    email_address           = "admin@contoso.com"
    use_common_alert_schema = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_itsmReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  itsm_receiver {
    name                 = "createorupdateticket"
    workspace_id         = "6eee3a18-aac3-40e4-b98e-1f309f329816"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{}"
    region               = "eastus"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_azureAppPushReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  azure_app_push_receiver {
    name          = "pushtoadmin"
    email_address = "admin@contoso.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_smsReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_webhookReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  webhook_receiver {
    name                    = "callmyapiaswell"
    service_uri             = "http://example.com/alert"
    use_common_alert_schema = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_automationRunbookReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"

  automation_runbook_receiver {
    name                    = "action_name_1"
    automation_account_id   = "${azurerm_automation_account.test.id}"
    runbook_name            = "${azurerm_automation_runbook.test.name}"
    webhook_resource_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
    is_global_runbook       = true
    service_uri             = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
    use_common_alert_schema = false
  }
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  log_verbose             = "true"
  log_progress            = "true"
  description             = "This is an test runbook"
  runbook_type            = "PowerShellWorkflow"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_voiceReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  voice_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_logicAppReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  logic_app_receiver {
    name                    = "logicappaction"
    resource_id             = azurerm_logic_app_workflow.test.id
    callback_url            = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
    use_common_alert_schema = true
  }
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestLA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_azureFunctionReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  azure_function_receiver {
    name                     = "funcaction"
    function_app_resource_id = azurerm_function_app.test.id
    function_name            = "myfunc"
    http_trigger_url         = "https://example.com/trigger"
    use_common_alert_schema  = true
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestSP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_armRoleReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  arm_role_receiver {
    name                    = "Monitoring Reader"
    role_id                 = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
    use_common_alert_schema = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
    name                    = "sendtodevops"
    email_address           = "devops@contoso.com"
    use_common_alert_schema = true
  }

  itsm_receiver {
    name                 = "createorupdateticket"
    workspace_id         = "6eee3a18-aac3-40e4-b98e-1f309f329816"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{}"
    region               = "eastus"
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
    country_code = "61"
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
    name                    = "action_name_1"
    automation_account_id   = "${azurerm_automation_account.test.id}"
    runbook_name            = "${azurerm_automation_runbook.test.name}"
    webhook_resource_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
    is_global_runbook       = true
    service_uri             = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
    use_common_alert_schema = false
  }

  voice_receiver {
    name         = "oncallvoice"
    country_code = "1"
    phone_number = "1231231234"
  }

  logic_app_receiver {
    name                    = "logicappaction"
    resource_id             = "${azurerm_logic_app_workflow.test.id}"
    callback_url            = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
    use_common_alert_schema = false
  }

  azure_function_receiver {
    name                     = "funcaction"
    function_app_resource_id = "${azurerm_function_app.test.id}"
    function_name            = "myfunc"
    http_trigger_url         = "https://example.com/trigger"
    use_common_alert_schema  = false
  }

  arm_role_receiver {
    name                    = "Monitoring Reader"
    role_id                 = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
    use_common_alert_schema = false
  }
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  log_verbose             = "true"
  log_progress            = "true"
  description             = "This is an test runbook"
  runbook_type            = "PowerShellWorkflow"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
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
  name                       = "acctestFA-%d"
  location                   = "${azurerm_resource_group.test.location}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
  app_service_plan_id        = "${azurerm_app_service_plan.test.id}"
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_disabledBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
  enabled             = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
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
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMMonitorActionGroup_secureWebhookReceiverTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azuread_application" "test" {
  name            = "acctestspa-%[1]d"
  identifier_uris = ["https://uri"]
  app_role {
    allowed_member_types = ["Application"]
    description          = "This is a role for Action Groups to join"
    display_name         = "ActionGroupsSecureWebhook"
    is_enabled           = true
    value                = "ActionGroupsSecureWebhook"
  }
}

resource "azuread_application" "test_alt" {
  name            = "acctestspa-alt-%[1]d"
  identifier_uris = ["https://uri2"]
  app_role {
    allowed_member_types = ["Application"]
    description          = "This is a role for Action Groups to join"
    display_name         = "ActionGroupsSecureWebhook"
    is_enabled           = true
    value                = "ActionGroupsSecureWebhook"
  }
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_service_principal" "test_alt" {
  application_id = azuread_application.test_alt.application_id
}

data "azuread_service_principal" "action_group" {
  application_id = "461e8683-5575-4561-ac7f-899cc907d62a"
}

# Here we miss one AAD resource to do New-AzureADServiceAppRoleAssignment -Id $myApp.AppRoles[0].Id -ResourceId $myServicePrincipal.ObjectId -ObjectId $actionGroupsSP.ObjectId -PrincipalId $actionGroupsSP.ObjectId
# Details in https://docs.microsoft.com/en-us/azure/azure-monitor/platform/action-groups#secure-webhook

`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMonitorActionGroup_secureWebhookReceiver(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActionGroup_secureWebhookReceiverTemplate(data)
	return fmt.Sprintf(`
%s 

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  webhook_receiver {
    name                    = "callmyapiaswell"
    service_uri             = "http://example.com/alert"
    use_common_alert_schema = true
  }

  webhook_receiver {
    name                    = "secureWebhook"
    service_uri             = "http://secureWebhook.com/alert"
    use_common_alert_schema = false
    use_aad_auth            = true
    aad_auth_object_id      = azuread_application.test.object_id
    aad_auth_identifier_uri = azuread_application.test.identifier_uris[0]
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_secureWebhookReceiverUpdate(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActionGroup_secureWebhookReceiverTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  webhook_receiver {
    name                    = "secureWebhook"
    service_uri             = "http://secureWebhook2.com/alert"
    use_common_alert_schema = true
    use_aad_auth            = true
    aad_auth_object_id      = azuread_application.test_alt.object_id
    aad_auth_identifier_uri = azuread_application.test_alt.identifier_uris[0]
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorActionGroup_secureWebhookReceiverRestore(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActionGroup_secureWebhookReceiverTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}
`, template, data.RandomInteger)
}
