package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
					resource.TestCheckResourceAttr(dataSourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_function_receiver.#", "0"),
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
					resource.TestCheckResourceAttr(dataSourceName, "itsm_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_app_push_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "automation_runbook_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "voice_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "logic_app_receiver.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_function_receiver.#", "0"),
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
					resource.TestCheckResourceAttr(dataSourceName, "itsm_receiver.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "itsm_receiver.0.workspace_id", "6eee3a18-aac3-40e4-b98e-1f309f329816"),
					resource.TestCheckResourceAttr(dataSourceName, "itsm_receiver.0.connection_id", "53de6956-42b4-41ba-be3c-b154cdf17b13"),
					resource.TestCheckResourceAttr(dataSourceName, "itsm_receiver.0.ticket_configuration", "{}"),
					resource.TestCheckResourceAttr(dataSourceName, "itsm_receiver.0.region", "southcentralus"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_app_push_receiver.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_app_push_receiver.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(dataSourceName, "sms_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.0.service_uri", "http://example.com/alert"),
					resource.TestCheckResourceAttr(dataSourceName, "webhook_receiver.1.service_uri", "https://backup.example.com/warning"),
					resource.TestCheckResourceAttr(dataSourceName, "automation_runbook_receiver.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "automation_runbook_receiver.automation_account_id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001"),
					resource.TestCheckResourceAttr(dataSourceName, "automation_runbook_receiver.runbook_name", "my runbook"),
					resource.TestCheckResourceAttr(dataSourceName, "automation_runbook_receiver.webhook_resource_id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"),
					resource.TestCheckResourceAttrSet(dataSourceName, "automation_runbook_receiver.is_global_runbook"),
					resource.TestCheckResourceAttr(dataSourceName, "automation_runbook_receiver.service_uri", "https://s13events.azure-automation.net/webhooks?token=randomtoken"),
					resource.TestCheckResourceAttr(dataSourceName, "voice_receiver.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "voice_receiver.0.country_code", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "voice_receiver.0.phone_number", "1231231234"),
					resource.TestCheckResourceAttr(dataSourceName, "voice_receiver.1.country_code", "86"),
					resource.TestCheckResourceAttr(dataSourceName, "voice_receiver.1.phone_number", "13888888888"),
					resource.TestCheckResourceAttr(dataSourceName, "logic_app_receiver.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "logic_app_receiver.resource_id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-logicapp/providers/Microsoft.Logic/workflows/logicapp"),
					resource.TestCheckResourceAttr(dataSourceName, "logic_app_receiver.callback_url", "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%2Fversions%2F08587100027316071865%2Ftriggers%2FmanualTrigger%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_function_receiver.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_function_receiver.function_app_resource_id", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_function_receiver.function_name", "myfunc"),
					resource.TestCheckResourceAttr(dataSourceName, "azure_function_receiver.http_trigger_url", "https://example.com/trigger"),
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

	itsm_receiver {
    name          = "createorupdateticket"
		workspace_id = "6eee3a18-aac3-40e4-b98e-1f309f329816"
		connection_id = "53de6956-42b4-41ba-be3c-b154cdf17b13"
		ticket_configuration = "{}"
		region = "southcentralus"
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
    name = "action_name_1"
    automation_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001"
    runbook_name = "my runbook"
    webhook_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg-runbooks/providers/microsoft.automation/automationaccounts/aaa001/webhooks/webhook_alert"
    is_global_runbook = true
    service_uri = "https://s13events.azure-automation.net/webhooks?token=randomtoken"
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
		name = "logicappaction"
		resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-logicapp/providers/Microsoft.Logic/workflows/logicapp"
		callback_url = "http://test-host:100/workflows/fb9c8d79b15f41ce9b12861862f43546/versions/08587100027316071865/triggers/manualTrigger/paths/invoke?api-version=2015-08-01-preview&sp=%%2Fversions%%2F08587100027316071865%%2Ftriggers%%2FmanualTrigger%%2Frun&sv=1.0&sig=IxEQ_ygZf6WNEQCbjV0Vs6p6Y4DyNEJVAa86U5B4xhk"
	}

	azure_function_receiver {
		name = "funcaction"
		function_app_resource_id = ""
		function_name = "myfunc"
		http_trigger_url = "https://example.com/trigger"
	}
}

data "azurerm_monitor_action_group" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_monitor_action_group.test.name}"
}
`, rInt, location, rInt)
}
