package monitor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmMonitorActionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmMonitorActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "short_name", "acctestag"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_function_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "arm_role_receiver.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceArmMonitorActionGroup_disabledBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmMonitorActionGroup_disabledBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "short_name", "acctestag"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_function_receiver.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "arm_role_receiver.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceArmMonitorActionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_action_group", "test")

	aaName := fmt.Sprintf("acctestAA-%d", data.RandomInteger)
	faName := fmt.Sprintf("acctestFA-%d", data.RandomInteger)
	laName := fmt.Sprintf("acctestLA-%d", data.RandomInteger)
	webhookName := "webhook_alert"
	resGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	aaResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName)
	aaWebhookResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/AutomationAccounts/%s/webhooks/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName, webhookName)
	faResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, faName)
	laResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, laName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmMonitorActionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_receiver.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_receiver.1.email_address", "devops@contoso.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_receiver.1.use_common_alert_schema", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "itsm_receiver.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "itsm_receiver.0.workspace_id", "6eee3a18-aac3-40e4-b98e-1f309f329816"),
					resource.TestCheckResourceAttr(data.ResourceName, "itsm_receiver.0.connection_id", "53de6956-42b4-41ba-be3c-b154cdf17b13"),
					resource.TestCheckResourceAttr(data.ResourceName, "itsm_receiver.0.ticket_configuration", "{}"),
					resource.TestCheckResourceAttr(data.ResourceName, "itsm_receiver.0.region", "southcentralus"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_app_push_receiver.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_app_push_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "sms_receiver.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(data.ResourceName, "sms_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(data.ResourceName, "sms_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(data.ResourceName, "webhook_receiver.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(data.ResourceName, "webhook_receiver.1.service_uri", "https://backup.example.com/warning"),
					resource.TestCheckResourceAttr(data.ResourceName, "webhook_receiver.1.use_common_alert_schema", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.0.automation_account_id", aaResourceID),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.0.runbook_name", webhookName),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.0.webhook_resource_id", aaWebhookResourceID),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.0.service_uri", "https://s13events.azure-automation.net/webhooks?token=randomtoken"),
					resource.TestCheckResourceAttr(data.ResourceName, "automation_runbook_receiver.0.use_common_alert_schema", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "voice_receiver.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "voice_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "voice_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(data.ResourceName, "logic_app_receiver.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "logic_app_receiver.0.resource_id", laResourceID),
					resource.TestCheckResourceAttr(data.ResourceName, "logic_app_receiver.0.callback_url", "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%2Fversions%2F08587100027316071865%2Ftriggers%2FmanualTrigger%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"),
					resource.TestCheckResourceAttr(data.ResourceName, "logic_app_receiver.0.use_common_alert_schema", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_function_receiver.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_function_receiver.0.function_app_resource_id", faResourceID),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_function_receiver.0.function_name", "myfunc"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_function_receiver.0.http_trigger_url", "https://example.com/trigger"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_function_receiver.0.use_common_alert_schema", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "arm_role_receiver.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "arm_role_receiver.0.role_id", "43d0d8ad-25c7-4714-9337-8ba259a9fe05"),
					resource.TestCheckResourceAttr(data.ResourceName, "arm_role_receiver.0.use_common_alert_schema", "false"),
				),
			},
		},
	})
}

func testAccDataSourceArmMonitorActionGroup_basic(data acceptance.TestData) string {
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

data "azurerm_monitor_action_group" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_monitor_action_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceArmMonitorActionGroup_disabledBasic(data acceptance.TestData) string {
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

data "azurerm_monitor_action_group" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_monitor_action_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceArmMonitorActionGroup_complete(data acceptance.TestData) string {
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
    name                    = "sendtoadmin"
    email_address           = "admin@contoso.com"
    use_common_alert_schema = false
  }

  email_receiver {
    name          = "sendtodevops"
    email_address = "devops@contoso.com"
  }

  itsm_receiver {
    name                 = "createorupdateticket"
    workspace_id         = "6eee3a18-aac3-40e4-b98e-1f309f329816"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{}"
    region               = "southcentralus"
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
    name        = "callmybackupapi"
    service_uri = "https://backup.example.com/warning"
  }

  automation_runbook_receiver {
    name                    = "action_name_1"
    automation_account_id   = "${azurerm_automation_account.test.id}"
    runbook_name            = "my runbook"
    webhook_resource_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
    is_global_runbook       = true
    service_uri             = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
    use_common_alert_schema = false
  }

  voice_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "1231231234"
  }

  voice_receiver {
    name         = "remotesupport"
    country_code = "86"
    phone_number = "13888888888"
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
  name                       = "acctestFA-%d"
  location                   = "${azurerm_resource_group.test.location}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
  app_service_plan_id        = "${azurerm_app_service_plan.test.id}"
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

data "azurerm_monitor_action_group" "test" {
  name                = "${azurerm_monitor_action_group.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}
