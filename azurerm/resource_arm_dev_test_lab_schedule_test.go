package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMDevTestLabSchedule_autoShutdownBasic(t *testing.T) {
	resourceName := "azurerm_dev_test_schedule.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMDevTestLabSchedule_autoShutdownBasic(ri, location)
	postConfig := testAccAzureRMDevTestLabSchedule_autoShutdownBasicUpdate(ri, location)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestLabScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_settings.0.status", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "daily_recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "daily_recurrence.0.time", "0100"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_settings.0.status", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "notification_settings.0.time_in_minutes", "30"),
					resource.TestCheckResourceAttr(resourceName, "notification_settings.0.webhook_url", "https://www.bing.com/2/4"),
					resource.TestCheckResourceAttr(resourceName, "daily_recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "daily_recurrence.0.time", "0900"),
				),
			},
		},
	})
}

func TestAccAzureRMDevTestLabSchedule_autoStartupBasic(t *testing.T) {
	resourceName := "azurerm_dev_test_schedule.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMDevTestLabSchedule_autoStartupBasic(ri, location)
	postConfig := testAccAzureRMDevTestLabSchedule_autoStartupBasicUpdate(ri, location)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestLabScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.0.time", "1100"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.0.week_days.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.0.week_days.1", "Tuesday"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"task_type"},
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.0.time", "1000"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.0.week_days.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "weekly_recurrence.0.week_days.1", "Thursday"),
				),
			},
		},
	})
}

func TestAccAzureRMDevTestLabSchedule_concurrent(t *testing.T) {
	firstResourceName := "azurerm_dev_test_schedule.test"
	secondResourceName := "azurerm_dev_test_schedule.test2"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMDevTestLabSchedule_concurrent(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestLabScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabScheduleExists(firstResourceName),
					testCheckAzureRMDevTestLabScheduleExists(secondResourceName),
				),
			},
		},
	})
}

func testCheckAzureRMDevTestLabScheduleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		devTestLabName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.LabSchedulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMDevTestLabScheduleDestroy(s *terraform.State) error {
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

func testAccAzureRMDevTestLabSchedule_autoShutdownBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmsShutdown"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
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
`, rInt, location, rInt)
}

func testAccAzureRMDevTestLabSchedule_autoShutdownBasicUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmsShutdown"
  status              = "Enabled"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
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

`, rInt, location, rInt)
}

func testAccAzureRMDevTestLabSchedule_autoStartupBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmAutoStart"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
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
`, rInt, location, rInt)
}

func testAccAzureRMDevTestLabSchedule_autoStartupBasicUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmAutoStart"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
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

`, rInt, location, rInt)
}

func testAccAzureRMDevTestLabSchedule_concurrent(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctdtl-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

}

resource "azurerm_dev_test_schedule" "test" {
  name                = "LabVmAutoStart"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
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
`, rInt, location, rInt)
}
