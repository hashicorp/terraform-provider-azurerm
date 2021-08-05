package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MonitorAutoScaleSettingResource struct {
}

func TestAccMonitorAutoScaleSetting_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.name").HasValue("metricRules"),
				check.That(data.ResourceName).Key("profile.0.rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.rule.0.metric_trigger.0.time_aggregation").HasValue("Last"),
				check.That(data.ResourceName).Key("notification.#").HasValue("0"),
				acceptance.TestCheckNoResourceAttr(data.ResourceName, "tags.$type"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorAutoScaleSetting_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_autoscale_setting"),
		},
	})
}

func TestAccMonitorAutoScaleSetting_multipleProfiles(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleProfiles(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("profile.#").HasValue("2"),
				check.That(data.ResourceName).Key("profile.0.name").HasValue("primary"),
				check.That(data.ResourceName).Key("profile.1.name").HasValue("secondary"),
			),
		},
	})
}

func TestAccMonitorAutoScaleSetting_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capacity(data, 1, 3, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.minimum").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.maximum").HasValue("3"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.default").HasValue("2"),
			),
		},
		{
			Config: r.capacity(data, 2, 4, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.minimum").HasValue("2"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.maximum").HasValue("4"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.default").HasValue("3"),
			),
		},
		{
			Config: r.capacity(data, 2, 45, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.minimum").HasValue("2"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.maximum").HasValue("45"),
				check.That(data.ResourceName).Key("profile.0.capacity.0.default").HasValue("3"),
			),
		},
	})
}

func TestAccMonitorAutoScaleSetting_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.name").HasValue("metricRules"),
				check.That(data.ResourceName).Key("profile.0.rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.rule.0.scale_action.0.direction").HasValue("Increase"),
				check.That(data.ResourceName).Key("notification.#").HasValue("0"),
			),
		},
		{
			Config: r.multipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.name").HasValue("metricRules"),
				check.That(data.ResourceName).Key("profile.0.rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("profile.0.rule.0.scale_action.0.direction").HasValue("Increase"),
				check.That(data.ResourceName).Key("profile.0.rule.1.scale_action.0.direction").HasValue("Decrease"),
				check.That(data.ResourceName).Key("notification.#").HasValue("0"),
			),
		},
	})
}

func TestAccMonitorAutoScaleSetting_customEmails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.email(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("notification.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification.0.email.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification.0.email.0.custom_emails.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification.0.email.0.custom_emails.0").HasValue(fmt.Sprintf("acctest1-%d@example.com", data.RandomInteger)),
			),
		},
		{
			Config: r.emailUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("notification.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification.0.email.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification.0.email.0.custom_emails.#").HasValue("2"),
				check.That(data.ResourceName).Key("notification.0.email.0.custom_emails.0").HasValue(fmt.Sprintf("acctest1-%d@example.com", data.RandomInteger)),
				check.That(data.ResourceName).Key("notification.0.email.0.custom_emails.1").HasValue(fmt.Sprintf("acctest2-%d@example.com", data.RandomInteger)),
			),
		},
	})
}

func TestAccMonitorAutoScaleSetting_recurrence(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurrence(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.name").HasValue("recurrence"),
				check.That(data.ResourceName).Key("profile.0.recurrence.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorAutoScaleSetting_recurrenceUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurrence(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("notification.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.#").HasValue("3"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.0").HasValue("Monday"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.1").HasValue("Wednesday"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.2").HasValue("Friday"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.hours.0").HasValue("18"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.minutes.0").HasValue("0"),
			),
		},
		{
			Config: r.recurrenceUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("profile.0.recurrence.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.#").HasValue("3"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.0").HasValue("Monday"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.1").HasValue("Tuesday"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.days.2").HasValue("Wednesday"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.hours.0").HasValue("20"),
				check.That(data.ResourceName).Key("profile.0.recurrence.0.minutes.0").HasValue("15"),
			),
		},
	})
}

func TestAccMonitorAutoScaleSetting_fixedDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fixedDate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("profile.0.name").HasValue("fixedDate"),
				check.That(data.ResourceName).Key("profile.0.fixed_date.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_multipleRulesDimensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_autoscale_setting", "test")
	r := MonitorAutoScaleSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleRulesDimensions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleRulesDimensionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t MonitorAutoScaleSettingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["autoscalesettings"]

	resp, err := clients.Monitor.AutoscaleSettingsClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading autoscale settings (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (MonitorAutoScaleSettingResource) basic(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) requiresImport(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.basic(data)
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

func (MonitorAutoScaleSettingResource) multipleProfiles(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) capacity(data acceptance.TestData, min int, max int, defaultVal int) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) multipleRules(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) email(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) emailUpdated(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) recurrence(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) recurrenceUpdated(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) fixedDate(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) multipleRulesDimensions(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) multipleRulesDimensionsUpdate(data acceptance.TestData) string {
	template := MonitorAutoScaleSettingResource{}.template(data)
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

func (MonitorAutoScaleSettingResource) template(data acceptance.TestData) string {
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

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
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
