package devtestlabs_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccDevTestLabSchedule_autoShutdownBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLabScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLabSchedule_autoShutdownBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.status", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.0.time", "0100"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDevTestLabSchedule_autoShutdownBasicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.status", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.time_in_minutes", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.webhook_url", "https://www.bing.com/2/4"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.0.time", "0900"),
				),
			},
		},
	})
}

func TestAccDevTestLabSchedule_autoStartupBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLabScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLabSchedule_autoStartupBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.0.time", "1100"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.0.week_days.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.0.week_days.1", "Tuesday"),
				),
			},
			data.ImportStep("task_type"),
			{
				Config: testAccDevTestLabSchedule_autoStartupBasicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.0.time", "1000"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.0.week_days.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "weekly_recurrence.0.week_days.1", "Thursday"),
				),
			},
		},
	})
}

func TestAccDevTestLabSchedule_concurrent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_schedule", "test")
	secondResourceName := "azurerm_dev_test_schedule.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLabScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLabSchedule_concurrent(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabScheduleExists(data.ResourceName),
					testCheckDevTestLabScheduleExists(secondResourceName),
				),
			},
		},
	})
}

func testCheckDevTestLabScheduleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.LabSchedulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		devTestLabName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, devTestLabName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on devTestLabSchedulesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Dev Test Lab Schedule %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckDevTestLabScheduleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.LabSchedulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_schedule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		devTestLabName := rs.Primary.Attributes["azurerm_dev_test_lab"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, devTestLabName, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Dev Test Lab Schedule still exists:\n%#v", resp.ScheduleProperties)
		}
	}

	return nil
}

func testAccDevTestLabSchedule_autoShutdownBasic(data acceptance.TestData) string {
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

func testAccDevTestLabSchedule_autoShutdownBasicUpdate(data acceptance.TestData) string {
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

func testAccDevTestLabSchedule_autoStartupBasic(data acceptance.TestData) string {
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

func testAccDevTestLabSchedule_autoStartupBasicUpdate(data acceptance.TestData) string {
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

func testAccDevTestLabSchedule_concurrent(data acceptance.TestData) string {
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
