package azurerm

import (
	"fmt"
	"testing"
)

func TestAccAzureRMAutoscaleSettings_basic(t *testing.T) {

	// ri := acctest.RandInt()
	// config := fmt.Sprintf(testAccAzureRMVAvailabilitySet_basic, ri, ri)

	// resource.Test(t, resource.TestCase{
	// 	PreCheck:     func() { testAccPreCheck(t) },
	// 	Providers:    testAccProviders,
	// 	CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
	// 	Steps: []resource.TestStep{
	// 		resource.TestStep{
	// 			Config: config,
	// 			Check: resource.ComposeTestCheckFunc(
	// 				testCheckAzureRMAvailabilitySetExists("azurerm_availability_set.test"),
	// 				resource.TestCheckResourceAttr(
	// 					"azurerm_availability_set.test", "platform_update_domain_count", "5"),
	// 				resource.TestCheckResourceAttr(
	// 					"azurerm_availability_set.test", "platform_fault_domain_count", "3"),
	// 			),
	// 		},
	// 	},
	// })
}

func testAccAzureRMAutoscaleSettings_basic(rInt int) string {
	return fmt.Sprintf(`
variable "vmss_uri" {
  type    = "string"
  default = "/subscriptions/c3adb315-00d7-4bd0-96be-c96210ced312/resourceGroups/autoscaling-test/providers/Microsoft.Compute/virtualMachineScaleSets/elast"
}

variable "time_zone" {
  type    = "string"
  default = "Pacific Standard Time"
}

data azurerm_resource_group "test" {
  name = "autoscaling-test"
}

resource azurerm_autoscale_settings "test" {
  name                = "tf_ss"
  enabled             = true
  resource_group_name = "${data.azurerm_resource_group.test.name}"
  location            = "${data.azurerm_resource_group.test.location}"
  target_resource_uri = "${var.vmss_uri}"
  profile {
    name = "profile1"
    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }
    rule {
      "metric_trigger" {
        "metric_name"         = "Percentage CPU"
        "metric_resource_uri" = "${var.vmss_uri}"
        "time_grain"          = "PT1M"
        "statistic"           = "Average"
        "time_window"         = "PT5M"
        "time_aggregation"    = "Average"
        "operator"            = "GreaterThan"
        "threshold"           = 75
      }
      "scale_action" {
        "direction" = "Increase"
        "type"      = "ChangeCount"
        "value"     = "1"
        "cooldown"  = "PT1M"
      }
    }
    rule {
      "metric_trigger" {
        "metric_name"         = "Percentage CPU"
        "metric_resource_uri" = "${var.vmss_uri}"
        "time_grain"          = "PT1M"
        "statistic"           = "Average"
        "time_window"         = "PT5M"
        "time_aggregation"    = "Average"
        "operator"            = "LessThan"
        "threshold"           = 25
      }
      "scale_action" {
        "direction" = "Decrease"
        "type"      = "ChangeCount"
        "value"     = "1"
        "cooldown"  = "PT1M"
      }
    }

/*
    "fixed_date" {
      "time_zone" = "${var.time_zone}"
      "start"     = "2017-06-17T00:00:00Z"
      "end"       = "2017-06-17T23:59:59Z"
    }
*/

    "recurrence" {
      "frequency" = "Week"
      "schedule" {
        "time_zone" = "${var.time_zone}"
        "days"    = [
          "Monday",
          "Wednesday",
          "Friday"
        ]
        "hours"   = [ 18 ]
        "minutes" = [ 0 ]
      }
    }
  }

  notification {
    "operation" = "Scale"
    "email" {
      "send_to_subscription_administrator"    = true
      "send_to_subscription_co_administrator" = false
      "custom_emails" = [ "foobar@asdf.com" ]
    }
    "webhook" {
      "service_uri" = "https://www.contoso.com/webhook"
      "properties"  = {
        foo = "bar"
      }
    }
  }
}`, rInt)
}
