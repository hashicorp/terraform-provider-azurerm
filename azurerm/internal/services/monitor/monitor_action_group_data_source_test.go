package monitor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MonitorActionGroupDataSource struct {
}

func TestAccDataSourceMonitorActionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_action_group", "test")
	r := MonitorActionGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("short_name").HasValue("acctestag"),
				check.That(data.ResourceName).Key("email_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("itsm_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_app_push_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("sms_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("webhook_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("automation_runbook_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("voice_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("logic_app_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_function_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("arm_role_receiver.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceMonitorActionGroup_disabledBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_action_group", "test")
	r := MonitorActionGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.disabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
				check.That(data.ResourceName).Key("short_name").HasValue("acctestag"),
				check.That(data.ResourceName).Key("email_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("itsm_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_app_push_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("sms_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("webhook_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("automation_runbook_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("voice_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("logic_app_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_function_receiver.#").HasValue("0"),
				check.That(data.ResourceName).Key("arm_role_receiver.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceMonitorActionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_action_group", "test")
	r := MonitorActionGroupDataSource{}

	aaName := fmt.Sprintf("acctestAA-%d", data.RandomInteger)
	faName := fmt.Sprintf("acctestFA-%d", data.RandomInteger)
	laName := fmt.Sprintf("acctestLA-%d", data.RandomInteger)
	webhookName := "webhook_alert"
	resGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	aaResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName)
	aaWebhookResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/webhooks/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName, webhookName)
	faResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, faName)
	laResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, laName)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("email_receiver.#").HasValue("2"),
				check.That(data.ResourceName).Key("email_receiver.0.email_address").HasValue("admin@contoso.com"),
				check.That(data.ResourceName).Key("email_receiver.1.email_address").HasValue("devops@contoso.com"),
				check.That(data.ResourceName).Key("email_receiver.1.use_common_alert_schema").HasValue("false"),
				check.That(data.ResourceName).Key("itsm_receiver.#").HasValue("1"),
				check.That(data.ResourceName).Key("itsm_receiver.0.workspace_id").HasValue("6eee3a18-aac3-40e4-b98e-1f309f329816"),
				check.That(data.ResourceName).Key("itsm_receiver.0.connection_id").HasValue("53de6956-42b4-41ba-be3c-b154cdf17b13"),
				check.That(data.ResourceName).Key("itsm_receiver.0.ticket_configuration").HasValue("{}"),
				check.That(data.ResourceName).Key("itsm_receiver.0.region").HasValue("southcentralus"),
				check.That(data.ResourceName).Key("azure_app_push_receiver.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_app_push_receiver.0.email_address").HasValue("admin@contoso.com"),
				check.That(data.ResourceName).Key("sms_receiver.#").HasValue("2"),
				check.That(data.ResourceName).Key("sms_receiver.0.country_code").HasValue("1"),
				check.That(data.ResourceName).Key("sms_receiver.0.phone_number").HasValue("1231231234"),
				check.That(data.ResourceName).Key("sms_receiver.1.country_code").HasValue("1"),
				check.That(data.ResourceName).Key("sms_receiver.1.phone_number").HasValue("5551238888"),
				check.That(data.ResourceName).Key("webhook_receiver.#").HasValue("2"),
				check.That(data.ResourceName).Key("webhook_receiver.0.service_uri").HasValue("http://example.com/alert"),
				check.That(data.ResourceName).Key("webhook_receiver.1.service_uri").HasValue("https://backup.example.com/warning"),
				check.That(data.ResourceName).Key("webhook_receiver.1.use_common_alert_schema").HasValue("false"),
				check.That(data.ResourceName).Key("webhook_receiver.1.aad_auth.#").HasValue("0"),
				check.That(data.ResourceName).Key("automation_runbook_receiver.#").HasValue("1"),
				check.That(data.ResourceName).Key("automation_runbook_receiver.0.automation_account_id").HasValue(aaResourceID),
				check.That(data.ResourceName).Key("automation_runbook_receiver.0.runbook_name").HasValue(webhookName),
				check.That(data.ResourceName).Key("automation_runbook_receiver.0.webhook_resource_id").HasValue(aaWebhookResourceID),
				check.That(data.ResourceName).Key("automation_runbook_receiver.0.service_uri").HasValue("https://s13events.azure-automation.net/webhooks?token=randomtoken"),
				check.That(data.ResourceName).Key("automation_runbook_receiver.0.use_common_alert_schema").HasValue("false"),
				check.That(data.ResourceName).Key("voice_receiver.#").HasValue("2"),
				check.That(data.ResourceName).Key("voice_receiver.0.country_code").HasValue("1"),
				check.That(data.ResourceName).Key("voice_receiver.0.phone_number").HasValue("1231231234"),
				check.That(data.ResourceName).Key("voice_receiver.1.country_code").HasValue("1"),
				check.That(data.ResourceName).Key("voice_receiver.1.phone_number").HasValue("5551238888"),
				check.That(data.ResourceName).Key("logic_app_receiver.#").HasValue("1"),
				check.That(data.ResourceName).Key("logic_app_receiver.0.resource_id").HasValue(laResourceID),
				check.That(data.ResourceName).Key("logic_app_receiver.0.callback_url").HasValue("http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%2Fversions%2F08587100027316071865%2Ftriggers%2FmanualTrigger%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"),
				check.That(data.ResourceName).Key("logic_app_receiver.0.use_common_alert_schema").HasValue("false"),
				check.That(data.ResourceName).Key("azure_function_receiver.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_function_receiver.0.function_app_resource_id").HasValue(faResourceID),
				check.That(data.ResourceName).Key("azure_function_receiver.0.function_name").HasValue("myfunc"),
				check.That(data.ResourceName).Key("azure_function_receiver.0.http_trigger_url").HasValue("https://example.com/trigger"),
				check.That(data.ResourceName).Key("azure_function_receiver.0.use_common_alert_schema").HasValue("false"),
				check.That(data.ResourceName).Key("arm_role_receiver.#").HasValue("1"),
				check.That(data.ResourceName).Key("arm_role_receiver.0.role_id").HasValue("43d0d8ad-25c7-4714-9337-8ba259a9fe05"),
				check.That(data.ResourceName).Key("arm_role_receiver.0.use_common_alert_schema").HasValue("false"),
			),
		},
	})
}

func (MonitorActionGroupDataSource) basic(data acceptance.TestData) string {
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

func (MonitorActionGroupDataSource) disabledBasic(data acceptance.TestData) string {
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

func (MonitorActionGroupDataSource) complete(data acceptance.TestData) string {
	aaName := fmt.Sprintf("acctestAA-%d", data.RandomInteger)
	webhookName := "webhook_alert"
	resGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	aaWebhookResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/webhooks/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, aaName, webhookName)
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
    name         = "remotesupportmsg"
    country_code = "1"
    phone_number = "5551238888"
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
    automation_account_id   = azurerm_automation_account.test.id
    runbook_name            = "webhook_alert"
    webhook_resource_id     = "%s"
    is_global_runbook       = true
    service_uri             = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
    use_common_alert_schema = false
  }

  voice_receiver {
    name         = "oncall"
    country_code = "1"
    phone_number = "1231231234"
  }

  voice_receiver {
    name         = "remotesupport"
    country_code = "1"
    phone_number = "5551238888"
  }

  logic_app_receiver {
    name                    = "logicappaction"
    resource_id             = azurerm_logic_app_workflow.test.id
    callback_url            = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
    use_common_alert_schema = false
  }

  azure_function_receiver {
    name                     = "funcaction"
    function_app_resource_id = azurerm_function_app.test.id
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT
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

data "azurerm_monitor_action_group" "test" {
  name                = azurerm_monitor_action_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, aaWebhookResourceID, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}
