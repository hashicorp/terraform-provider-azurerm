package azurerm

import (
	"fmt"
	"testing"

	//"github.com/hashicorp/terraform/helper/acctest"
	//"github.com/hashicorp/terraform/helper/resource"
	//"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

// These tests are actually run as part of the resoure ones due to
// Azure only being happy about provisioning one per subscription at once
// (which our test suite can't easily workaround)

func testAccDataSourceAzureRMMonitorScheduledQueryRules_logToMetricAction(t *testing.T) {
	/*dataSourceName := "data.azurerm_monitor_scheduled_query_rules.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorScheduledQueryRules_logToMetricActionConfig(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locations.#"),
				),
			},
		},
	})
	*/
	return
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingAction(t *testing.T) {
	/*dataSourceName := "data.azurerm_monitor_scheduled_query_rules.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionConfig(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locations.#"),
				),
			},
		},
	})*/
	return
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionCrossResource(t *testing.T) {
/*	dataSourceName := "data.azurerm_monitor_scheduled_query_rules.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionCrossResourceConfig(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locations.#"),
				),
			},
		},
	})*/
	return
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_logToMetricActionConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-01"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
}

resource "azurerm_monitor_scheduled_query_rules" "test" {
  name                = "acctestehns-%s"
	location            = "${azurerm_resource_group.test.location}"
	description         = "test alerting action"
	enabled             = true
	type                = "LogToMetricAction"

	source {
		query        = "let data=datatable(id:int, value:string) [1, 'test1', 2, 'testtwo']; data | extend strlen = strlen(value)"
		dataSourceId = "${azurerm_log_analytics_workspace.test.id}"
		queryType    = "ResultCount"
	}

	schedule {
		frequencyInMinutes  = 60
    timeWindowInMinutes = 60
	}

	action {
		severity     = 3
    aznsAction {
      actionGroup = [
        "${azurerm_monitor_action_group.test.id}"
      ]
      emailSubject": "Custom alert email subject"
		}

    trigger {
      thresholdOperator = GreaterThan"
			threshold         = 5000
			metricTrigger {
				thresholdOperator = "GreaterThan"
				threshold         = 5
				metricTriggerType = "Consecutive"
				metricColumn      = "Computer"
			}
    }
	}
}

data "azurerm_monitor_scheduled_query_rules" "test" {
  name = "${azurerm_monitor_alerting_action.test.name}"
}
`, rInt, location, rString, rInt, location)
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-01"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
}

resource "azurerm_monitor_scheduled_query_rules" "test" {
  name                = "acctestehns-%s"
	location            = "${azurerm_resource_group.test.location}"
	description         = "test alerting action"
	enabled             = true
	type                = "AlertingAction"

	source {
		query        = "let data=datatable(id:int, value:string) [1, 'test1', 2, 'testtwo']; data | extend strlen = strlen(value)"
		dataSourceId = "${azurerm_log_analytics_workspace.test.id}"
		queryType    = "ResultCount"
	}

	schedule {
		frequencyInMinutes  = 60
    timeWindowInMinutes = 60
	}

	action {
		severity     = 3
    aznsAction {
      actionGroup = [
        "${azurerm_monitor_action_group.test.id}"
      ]
      emailSubject": "Custom alert email subject"
		}

    trigger {
      thresholdOperator = "GreaterThan"
      threshold         = 5000
    }
	}
}

data "azurerm_monitor_scheduled_query_rules" "test" {
  name = "${azurerm_monitor_alerting_action.test.name}"
}
`, rInt, location, rString, rInt, location)
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionCrossResourceConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-01"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace" "test2" {
  name                = "acctest-02"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
}

resource "azurerm_monitor_scheduled_query_rules" "test" {
  name                = "acctestehns-%s"
	location            = "${azurerm_resource_group.test.location}"
	description         = "test alerting action"
	enabled             = true
	type                = "AlertingAction"

	source {
		query        = "let data=datatable(id:int, value:string) [1, 'test1', 2, 'testtwo']; data | extend strlen = strlen(value)"
		dataSourceId = "${azurerm_log_analytics_workspace.test.id}"
		queryType    = "ResultCount"
		"authorizedResources": [
			"${azurerm_log_analytics_workspace.test2.id}"
      ],
	}

	schedule {
		frequencyInMinutes  = 60
    timeWindowInMinutes = 60
	}

	action {
		severity     = 3
    aznsAction {
      actionGroup = [
        "${azurerm_monitor_action_group.test.id}"
      ]
      emailSubject": "Custom alert email subject"
		}

    trigger {
      thresholdOperator = GreaterThan"
      threshold         = 5000
    }
	}
}

data "azurerm_monitor_scheduled_query_rules" "test" {
  name = "${azurerm_monitor_alerting_action.test.name}"
}
`, rInt, location, rString, rInt, location)
}
