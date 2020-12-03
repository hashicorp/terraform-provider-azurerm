package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccLogicAppTriggerRecurrence_month(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Month", 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Month"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Month", 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Month"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "1"),
				),
			},
			{
				Config:      testAccAzureRMLogicAppTriggerRecurrence_requiresImport(data, "Month", 1),
				ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_recurrence"),
			},
		},
	})
}

func TestAccLogicAppTriggerRecurrence_week(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Week", 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Week"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Day", 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Day"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_minute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Minute", 4),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Minute"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_second(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Second", 30),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Second"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "30"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_hour(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Hour", 4),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Hour"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Month", 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Month"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "1"),
				),
			},
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_basic(data, "Month", 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "Month"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval", "3"),
				),
			},
		},
	})
}

func TestAccLogicAppTriggerRecurrence_startTime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_startTime(data, "2020-01-01T01:02:03Z"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_time", "2020-01-01T01:02:03Z"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_startTimeWithTimeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_startTimeWithTimeZone(data, "2020-01-01T01:02:03Z", "US Eastern Standard Time"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_time", "2020-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr(data.ResourceName, "time_zone", "US Eastern Standard Time"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogicAppTriggerRecurrence_startTimeWithTimeZone(data, "2020-01-01T01:02:03Z", "Egypt Standard Time"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppTriggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_time", "2020-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr(data.ResourceName, "time_zone", "Egypt Standard Time"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMLogicAppTriggerRecurrence_basic(data acceptance.TestData, frequency string, interval int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_logic_app_trigger_recurrence" "test" {
  name         = "frequency-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  frequency    = "%s"
  interval     = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, frequency, interval)
}

func testAccAzureRMLogicAppTriggerRecurrence_startTime(data acceptance.TestData, startTime string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_logic_app_trigger_recurrence" "test" {
  name         = "frequency-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  frequency    = "Month"
  interval     = 1
  start_time   = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, startTime)
}

func testAccAzureRMLogicAppTriggerRecurrence_startTimeWithTimeZone(data acceptance.TestData, startTime string, timeZone string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_logic_app_trigger_recurrence" "test" {
  name         = "frequency-trigger"
  logic_app_id = azurerm_logic_app_workflow.test.id
  frequency    = "Month"
  interval     = 1
  start_time   = "%s"
  time_zone    = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, startTime, timeZone)
}

func testAccAzureRMLogicAppTriggerRecurrence_requiresImport(data acceptance.TestData, frequency string, interval int) string {
	template := testAccAzureRMLogicAppTriggerRecurrence_basic(data, frequency, interval)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_recurrence" "import" {
  name         = azurerm_logic_app_trigger_recurrence.test.name
  logic_app_id = azurerm_logic_app_trigger_recurrence.test.logic_app_id
  frequency    = azurerm_logic_app_trigger_recurrence.test.frequency
  interval     = azurerm_logic_app_trigger_recurrence.test.interval
}
`, template)
}
