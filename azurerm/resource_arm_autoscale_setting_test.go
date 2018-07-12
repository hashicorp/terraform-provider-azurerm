package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMAutoscaleSetting_basic(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoscaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutoscaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMAutoscaleSetting_recurrence(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_recurrence(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoscaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutoscaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "recurrence"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMAutoscaleSetting_fixedDate(t *testing.T) {
	resourceName := "azurerm_autoscale_setting.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAutoscaleSetting_fixedDate(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutoscaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutoscaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "profile.1.fixed_date.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
				),
			},
		},
	})
}

func testAccAzureRMAutoscaleSetting_basic(rInt int, location string) string {
	vmssBasic := testAccAzureRMVirtualMachineScaleSet_basic(rInt, location)
	return fmt.Sprintf(`%s
resource "azurerm_autoscale_setting" "test" {
  name                = "acctestAutoscaleSetting-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
  enabled             = true

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id   = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain           = "PT1M"
        statistic            = "Average"
        time_window          = "PT5M"
        time_aggregation     = "Average"
        operator             = "GreaterThan"
        threshold            = 75
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "LessThan"
        threshold           = 25
      }
      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
  }
}
`, vmssBasic, rInt)
}

func testAccAzureRMAutoscaleSetting_recurrence(rInt int, location string) string {
	vmssBasic := testAccAzureRMVirtualMachineScaleSet_basic(rInt, location)
	return fmt.Sprintf(`%s
resource "azurerm_autoscale_setting" "test" {
  name                = "acctestAutoscaleSetting-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
  enabled             = true

  profile {
    name = "recurrence"
    recurrence {
      frequency = "Week"
      schedule {
        time_zone = "Pacific Standard Time"
        days      = [
          "Monday",
          "Wednesday",
          "Friday"
        ]
        hours     = [ 18 ]
        minutes   = [ 0 ]
      }
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
    }
  }
}`, vmssBasic, rInt)
}

func testAccAzureRMAutoscaleSetting_fixedDate(rInt int, location string) string {
	vmssBasic := testAccAzureRMVirtualMachineScaleSet_basic(rInt, location)
	return fmt.Sprintf(`%s
resource "azurerm_autoscale_setting" "test" {
  name                = "acctestAutoscaleSetting-%d"
  enabled             = true
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "GreaterThan"
        threshold           = 75
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name         = "Percentage CPU"
        metric_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain          = "PT1M"
        statistic           = "Average"
        time_window         = "PT5M"
        time_aggregation    = "Average"
        operator            = "LessThan"
        threshold           = 25
      }
      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = "1"
        cooldown  = "PT1M"
      }
    }
  }

  profile {
    name = "fixedDate"
    fixed_date {
      time_zone = "Pacific Standard Time"
      start     = "2020-06-18T00:00:00Z"
      end       = "2020-06-18T23:59:59Z"
    }
  }
}`, vmssBasic, rInt)
}

func testCheckAzureRMAutoscaleSettingExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		autoscaleSettingName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Autoscale Setting: %s", autoscaleSettingName)
		}

		conn := testAccProvider.Meta().(*ArmClient).autoscaleSettingsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, autoscaleSettingName)
		if err != nil {
			return fmt.Errorf("Bad: Get on Autoscale Setting: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Autoscale Setting Instance %q (resource group: %q) does not exist", autoscaleSettingName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMAutoscaleSettingDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).autoscaleSettingsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_autoscale_setting" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Autoscale Setting still exists:\n%#v", resp)
		}
	}

	return nil
}
