package devtestlabs_test

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

type DevTestLabScheduleResource struct {
}

func TestAccDevTestLabSchedule_autoShutdownBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_schedule", "test")
	r := DevTestLabScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoShutdownBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Disabled"),
				check.That(data.ResourceName).Key("notification_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification_settings.0.status").HasValue("Disabled"),
				check.That(data.ResourceName).Key("daily_recurrence.#").HasValue("1"),
				check.That(data.ResourceName).Key("daily_recurrence.0.time").HasValue("0100"),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoShutdownBasicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Enabled"),
				check.That(data.ResourceName).Key("notification_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification_settings.0.status").HasValue("Enabled"),
				check.That(data.ResourceName).Key("notification_settings.0.time_in_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("notification_settings.0.webhook_url").HasValue("https://www.bing.com/2/4"),
				check.That(data.ResourceName).Key("daily_recurrence.#").HasValue("1"),
				check.That(data.ResourceName).Key("daily_recurrence.0.time").HasValue("0900"),
			),
		},
	})
}

func TestAccDevTestLabSchedule_autoStartupBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_schedule", "test")
	r := DevTestLabScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoStartupBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Disabled"),
				check.That(data.ResourceName).Key("weekly_recurrence.#").HasValue("1"),
				check.That(data.ResourceName).Key("weekly_recurrence.0.time").HasValue("1100"),
				check.That(data.ResourceName).Key("weekly_recurrence.0.week_days.#").HasValue("2"),
				check.That(data.ResourceName).Key("weekly_recurrence.0.week_days.1").HasValue("Tuesday"),
			),
		},
		data.ImportStep("task_type"),
		{
			Config: r.autoStartupBasicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Enabled"),
				check.That(data.ResourceName).Key("weekly_recurrence.#").HasValue("1"),
				check.That(data.ResourceName).Key("weekly_recurrence.0.time").HasValue("1000"),
				check.That(data.ResourceName).Key("weekly_recurrence.0.week_days.#").HasValue("3"),
				check.That(data.ResourceName).Key("weekly_recurrence.0.week_days.1").HasValue("Thursday"),
			),
		},
	})
}

func TestAccDevTestLabSchedule_concurrent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_schedule", "test")
	r := DevTestLabScheduleResource{}
	secondResourceName := "azurerm_dev_test_schedule.test2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.concurrent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (DevTestLabScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	devTestLabName := id.Path["labs"]
	name := id.Path["schedules"]

	resp, err := clients.DevTestLabs.LabSchedulesClient.Get(ctx, id.ResourceGroup, devTestLabName, name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Dev Test Lab Schedule %q (resource group: %q): %+v", name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ScheduleProperties != nil), nil
}

func (DevTestLabScheduleResource) autoShutdownBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmsShutdown"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lab_name            = azurerm_dev_test_lab.test.name
  daily_recurrence {
    time = "0100"
  }
  time_zone_id = "India Standard Time"
  task_type    = "LabVmsShutdownTask"
  notification_settings {
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DevTestLabScheduleResource) autoShutdownBasicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmsShutdown"
  status              = "Enabled"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lab_name            = azurerm_dev_test_lab.test.name
  daily_recurrence {
    time = "0900"
  }
  time_zone_id = "India Standard Time"
  task_type    = "LabVmsShutdownTask"
  notification_settings {
    time_in_minutes = 30
    webhook_url     = "https://www.bing.com/2/4"
    status          = "Enabled"
  }
  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DevTestLabScheduleResource) autoStartupBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmAutoStart"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lab_name            = azurerm_dev_test_lab.test.name
  weekly_recurrence {
    time      = "1100"
    week_days = ["Monday", "Tuesday"]
  }

  time_zone_id = "India Standard Time"
  task_type    = "LabVmsStartupTask"

  notification_settings {
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DevTestLabScheduleResource) autoStartupBasicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmAutoStart"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lab_name            = azurerm_dev_test_lab.test.name
  weekly_recurrence {
    time      = "1000"
    week_days = ["Wednesday", "Thursday", "Friday"]
  }

  time_zone_id = "India Standard Time"
  task_type    = "LabVmsStartupTask"

  notification_settings {
  }

  status = "Enabled"

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DevTestLabScheduleResource) concurrent(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmAutoStart"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lab_name            = azurerm_dev_test_lab.test.name
  weekly_recurrence {
    time      = "1100"
    week_days = ["Monday", "Tuesday"]
  }

  time_zone_id = "India Standard Time"
  task_type    = "LabVmsStartupTask"

  notification_settings {
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_dev_test_schedule" "test2" {
  name                = "LabVmsShutdown"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lab_name            = azurerm_dev_test_lab.test.name
  daily_recurrence {
    time = "0100"
  }
  time_zone_id = "India Standard Time"
  task_type    = "LabVmsShutdownTask"
  notification_settings {
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
