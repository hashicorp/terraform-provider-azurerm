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

func TestAccAzureRMMonitorAutoScaleSetting_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.rule.0.metric_trigger.0.time_aggregation", "Last"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "0"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "tags.$type"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorAutoScaleSetting_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_monitor_autoscale_setting"),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_multipleProfiles(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_multipleProfiles(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.name", "primary"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.1.name", "secondary"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_capacity(data, 1, 3, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.minimum", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.maximum", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.default", "2"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_capacity(data, 2, 4, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.minimum", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.maximum", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.default", "3"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_capacity(data, 2, 45, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.minimum", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.maximum", "45"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.capacity.0.default", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.rule.0.scale_action.0.direction", "Increase"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "0"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_multipleRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.rule.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.rule.0.scale_action.0.direction", "Increase"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.rule.1.scale_action.0.direction", "Decrease"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_customEmails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_email(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.0.email.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.0.email.0.custom_emails.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.0.email.0.custom_emails.0", fmt.Sprintf("acctest1-%d@example.com", data.RandomInteger)),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_emailUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.0.email.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.0.email.0.custom_emails.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.0.email.0.custom_emails.0", fmt.Sprintf("acctest1-%d@example.com", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.0.email.0.custom_emails.1", fmt.Sprintf("acctest2-%d@example.com", data.RandomInteger)),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_recurrence(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_recurrence(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.name", "recurrence"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_recurrenceUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_recurrence(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.1", "Wednesday"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.2", "Friday"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.hours.0", "18"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.minutes.0", "0"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_recurrenceUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.1", "Tuesday"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.days.2", "Wednesday"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.hours.0", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.recurrence.0.minutes.0", "15"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_fixedDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_fixedDate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.name", "fixedDate"),
					resource.TestCheckResourceAttr(data.ResourceName, "profile.0.fixed_date.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_multipleRulesDimensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_multipleRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_multipleRulesDimensions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_multipleRulesDimensionsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_multipleRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMonitorAutoScaleSettingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.AutoscaleSettingsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		autoscaleSettingName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Monitor AutoScale Setting: %s", autoscaleSettingName)
		}

		resp, err := conn.Get(ctx, resourceGroup, autoscaleSettingName)
		if err != nil {
			return fmt.Errorf("Bad: Get on Monitor AutoScale Setting: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: AutoScale Setting %q (Resource Group: %q) does not exist", autoscaleSettingName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMonitorAutoScaleSettingDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_autoscale_setting" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("AutoScale Setting still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMonitorAutoScaleSetting_basic(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Last"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "import" {
  name                = azurerm_monitor_autoscale_setting.test.name
  resource_group_name = azurerm_monitor_autoscale_setting.test.resource_group_name
  location            = azurerm_monitor_autoscale_setting.test.location
  target_resource_id  = azurerm_monitor_autoscale_setting.test.target_resource_id

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template)
}

func testAccAzureRMMonitorAutoScaleSetting_multipleProfiles(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id

  profile {
    name = "primary"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }

  profile {
    name = "secondary"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    recurrence {
      timezone = "Pacific Standard Time"

      days = [
        "Monday",
        "Wednesday",
        "Friday",
      ]

      hours   = [18]
      minutes = [0]
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_capacity(data acceptance.TestData, min int, max int, defaultVal int) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id
  enabled             = false

  profile {
    name = "metricRules"

    capacity {
      default = %d
      minimum = %d
      maximum = %d
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, data.RandomInteger, defaultVal, min, max)
}

func testAccAzureRMMonitorAutoScaleSetting_multipleRules(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id
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
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
        metric_namespace   = "microsoft.compute/virtualmachinescalesets"
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 25
        metric_namespace   = "microsoft.compute/virtualmachinescalesets"
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_email(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
      custom_emails                         = ["acctest1-%d@example.com"]
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_emailUpdated(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
      custom_emails                         = ["acctest1-%d@example.com", "acctest2-%d@example.com"]
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_recurrence(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id

  profile {
    name = "recurrence"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    recurrence {
      timezone = "Pacific Standard Time"

      days = [
        "Monday",
        "Wednesday",
        "Friday",
      ]

      hours   = [18]
      minutes = [0]
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_recurrenceUpdated(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id

  profile {
    name = "recurrence"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    recurrence {
      timezone = "Pacific Standard Time"

      days = [
        "Monday",
        "Tuesday",
        "Wednesday",
      ]

      hours   = [20]
      minutes = [15]
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_fixedDate(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id

  profile {
    name = "fixedDate"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    fixed_date {
      timezone = "Pacific Standard Time"
      start    = "2020-06-18T00:00:00Z"
      end      = "2020-06-18T23:59:59Z"
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_multipleRulesDimensions(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id
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
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
        metric_namespace   = "microsoft.compute/virtualmachinescalesets"
        dimensions {
          name     = "AppName"
          operator = "Equals"
          values   = ["App1"]
        }
      }


      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 25
        metric_namespace   = "microsoft.compute/virtualmachinescalesets"
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_multipleRulesDimensionsUpdate(data acceptance.TestData) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  target_resource_id  = azurerm_virtual_machine_scale_set.test.id
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
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
        metric_namespace   = "microsoft.compute/virtualmachinescalesets"
        dimensions {
          name     = "AppName2"
          operator = "NotEquals"
          values   = ["App2"]
        }

        dimensions {
          name     = "Deployment"
          operator = "Equals"
          values   = ["default"]
        }
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = azurerm_virtual_machine_scale_set.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 25
        metric_namespace   = "microsoft.compute/virtualmachinescalesets"
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorAutoScaleSetting_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile-%d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = azurerm_subnet.test.id
      primary   = true
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "StandardSSD_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
