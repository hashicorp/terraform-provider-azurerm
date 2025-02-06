// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorActionGroupResource struct{}

func TestAccMonitorActionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_action_group"),
		},
	})
}

func TestAccMonitorActionGroup_emailReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emailReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_itsmReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.itsmReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_azureAppPushReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureAppPushReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_smsReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.smsReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_webhookReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webhookReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

/*
@favoretti: Disabling this one, since it's written in such a way that it will never succeed in CI

	func TestAccMonitorActionGroup_secureWebhookReceiver(t *testing.T) {
		data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
		r := MonitorActionGroupResource{}

		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.basic(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.webhookReceiver(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.secureWebhookReceiver(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.basic(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
		})
	}
*/
func TestAccMonitorActionGroup_automationRunbookReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automationRunbookReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_voiceReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.voiceReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_logicAppReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logicAppReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_azureFunctionReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureFunctionReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_armRoleReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.armRoleReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_eventHubReceiver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventHubReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_disabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.disabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_location(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.location(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_singleReceiverUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emailReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.itsmReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.azureAppPushReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.smsReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.webhookReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.automationRunbookReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.voiceReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logicAppReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.azureFunctionReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.armRoleReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.eventHubReceiver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionGroup_multipleReceiversUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_group", "test")
	r := MonitorActionGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MonitorActionGroupResource) basic(data acceptance.TestData) string {
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

func (r MonitorActionGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_group" "import" {
  name                = azurerm_monitor_action_group.test.name
  resource_group_name = azurerm_monitor_action_group.test.resource_group_name
  short_name          = azurerm_monitor_action_group.test.short_name
}
`, r.basic(data))
}

func (MonitorActionGroupResource) emailReceiver(data acceptance.TestData) string {
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

func (MonitorActionGroupResource) itsmReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  itsm_receiver {
    name                 = "createorupdateticket"
    workspace_id         = "${data.azurerm_client_config.current.subscription_id}|${azurerm_log_analytics_workspace.test.workspace_id}"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{\"PayloadRevision\":0,\"WorkItemType\":\"Incident\",\"UseTemplate\":false,\"WorkItemData\":\"{}\",\"CreateOneWIPerCI\":false}"
    region               = "eastus"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (MonitorActionGroupResource) azureAppPushReceiver(data acceptance.TestData) string {
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

func (MonitorActionGroupResource) smsReceiver(data acceptance.TestData) string {
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
    phone_number = "2123456789"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MonitorActionGroupResource) webhookReceiver(data acceptance.TestData) string {
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

/*
@favoretti: Disabling this one, since it's written in such a way that it will never succeed in CI

func (MonitorActionGroupResource) secureWebhookReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azuread_application" "test" {
  object_id = "%s"
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

  webhook_receiver {
    name                    = "callmysecureapi"
    service_uri             = "http://secureExample.com/alert"
    use_common_alert_schema = true
    aad_auth {
      object_id      = data.azuread_application.test.object_id
      identifier_uri = data.azuread_application.test.identifier_uris[0]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_APP_OBJECT_ID"), data.RandomInteger)
}

*/

func (MonitorActionGroupResource) automationRunbookReceiver(data acceptance.TestData) string {
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

func (MonitorActionGroupResource) voiceReceiver(data acceptance.TestData) string {
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
    phone_number = "2123456789"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MonitorActionGroupResource) logicAppReceiver(data acceptance.TestData) string {
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

func (MonitorActionGroupResource) azureFunctionReceiver(data acceptance.TestData) string {
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

func (MonitorActionGroupResource) armRoleReceiver(data acceptance.TestData) string {
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

func (MonitorActionGroupResource) eventHubReceiver(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acceptanceTestEventHubNamespace-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "test" {
  name                = "acceptanceTestEventHub"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"

  event_hub_receiver {
    name                    = "eventhub-test-action"
    event_hub_name          = azurerm_eventhub.test.name
    event_hub_namespace     = azurerm_eventhub.test.namespace_name
    use_common_alert_schema = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MonitorActionGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
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
    workspace_id         = "${data.azurerm_client_config.current.subscription_id}|${azurerm_log_analytics_workspace.test.workspace_id}"
    connection_id        = "53de6956-42b4-41ba-be3c-b154cdf17b13"
    ticket_configuration = "{\"PayloadRevision\":0,\"WorkItemType\":\"Incident\",\"UseTemplate\":false,\"WorkItemData\":\"{}\",\"CreateOneWIPerCI\":false}"
    region               = "eastus"
  }

  azure_app_push_receiver {
    name          = "pushtoadmin"
    email_address = "admin@contoso.com"
  }

  sms_receiver {
    name         = "oncallmsg"
    country_code = "1"
    phone_number = "2123456789"
  }

  sms_receiver {
    name         = "remotesupport"
    country_code = "1"
    phone_number = "3123456789"
  }

  webhook_receiver {
    name        = "callmyapiaswell"
    service_uri = "http://example.com/alert"
  }

  webhook_receiver {
    name                    = "callmybackupapi"
    service_uri             = "https://example.com/warning/backup"
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
    phone_number = "2123456789"
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

  event_hub_receiver {
    name                    = "eventhub-test-action"
    event_hub_name          = azurerm_eventhub.test.name
    event_hub_namespace     = azurerm_eventhub.test.namespace_name
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

resource "azurerm_eventhub_namespace" "test" {
  name                = "acceptanceTestEventHubNamespace-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "test" {
  name                = "acceptanceTestEventHub"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}


`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (MonitorActionGroupResource) disabledBasic(data acceptance.TestData) string {
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

func (t MonitorActionGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := actiongroupsapis.ParseActionGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Monitor.ActionGroupsClient.ActionGroupsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading (%s): %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (MonitorActionGroupResource) location(data acceptance.TestData) string {
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
  location            = "swedencentral"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
