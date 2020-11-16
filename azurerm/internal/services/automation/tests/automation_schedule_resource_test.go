package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationSchedule_oneTime_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_basic(data),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_basic(data.ResourceName),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_basic(data),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_basic(data.ResourceName),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationSchedule_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationSchedule_oneTime_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	// the API returns the time in the timezone we pass in
	// it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("Australia/Perth")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_complete(data, startTime),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_complete(data.ResourceName, startTime),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_oneTime_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	// the API returns the time in the timezone we pass in
	// it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("Australia/Perth")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_basic(data),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_basic(data.ResourceName),
			},
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_complete(data, startTime),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_complete(data.ResourceName, startTime),
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_hourly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(data, "Hour", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(data.ResourceName, "Hour", 7),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_daily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(data, "Day", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(data.ResourceName, "Day", 7),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_weekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(data, "Week", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(data.ResourceName, "Week", 7),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(data, "Month", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(data.ResourceName, "Month", 7),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_weekly_advanced(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_advanced_week(data, "Monday"),
				Check:  checkAccAzureRMAutomationSchedule_recurring_advanced_week(data.ResourceName),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthly_advanced_by_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_advanced_month(data, 2),
				Check:  checkAccAzureRMAutomationSchedule_recurring_advanced_month(data.ResourceName),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthly_advanced_by_week_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(data, "Monday", 2),
				Check:  checkAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(data.ResourceName, "Monday", 2),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAutomationScheduleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.ScheduleClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_schedule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Schedule: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Schedule still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationScheduleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.ScheduleClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Schedule: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Schedule '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationScheduleClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAutomationSchedule_prerequisites(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationSchedule_oneTime_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
}
`, testAccAzureRMAutomationSchedule_prerequisites(data), data.RandomInteger)
}

func testAccAzureRMAutomationSchedule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationSchedule_oneTime_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "import" {
  name                    = azurerm_automation_schedule.test.name
  resource_group_name     = azurerm_automation_schedule.test.resource_group_name
  automation_account_name = azurerm_automation_schedule.test.automation_account_name
  frequency               = azurerm_automation_schedule.test.frequency
}
`, template)
}

func checkAccAzureRMAutomationSchedule_oneTime_basic(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
	)
}

func testAccAzureRMAutomationSchedule_oneTime_complete(data acceptance.TestData, startTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
  start_time              = "%s"
  timezone                = "Australia/Perth"
  description             = "This is an automation schedule"
}
`, testAccAzureRMAutomationSchedule_prerequisites(data), data.RandomInteger, startTime)
}

func checkAccAzureRMAutomationSchedule_oneTime_complete(resourceName, startTime string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
		resource.TestCheckResourceAttr(resourceName, "start_time", startTime),
		resource.TestCheckResourceAttr(resourceName, "timezone", "Australia/Perth"),
		resource.TestCheckResourceAttr(resourceName, "description", "This is an automation schedule"),
	)
}

// nolint unparam
func testAccAzureRMAutomationSchedule_recurring_basic(data acceptance.TestData, frequency string, interval int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "%s"
  interval                = "%d"
}
`, testAccAzureRMAutomationSchedule_prerequisites(data), data.RandomInteger, frequency, interval)
}

// nolint unparam
func checkAccAzureRMAutomationSchedule_recurring_basic(resourceName string, frequency string, interval int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", frequency),
		resource.TestCheckResourceAttr(resourceName, "interval", strconv.Itoa(interval)),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
	)
}

func testAccAzureRMAutomationSchedule_recurring_advanced_week(data acceptance.TestData, weekDay string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "Week"
  interval                = "1"
  week_days               = ["%s"]
}
`, testAccAzureRMAutomationSchedule_prerequisites(data), data.RandomInteger, weekDay)
}

func checkAccAzureRMAutomationSchedule_recurring_advanced_week(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists("azurerm_automation_schedule.test"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "Week"),
		resource.TestCheckResourceAttr(resourceName, "interval", "1"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
		resource.TestCheckResourceAttr(resourceName, "week_days.#", "1"),
	)
}

func testAccAzureRMAutomationSchedule_recurring_advanced_month(data acceptance.TestData, monthDay int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "Month"
  interval                = "1"
  month_days              = [%d]
}
`, testAccAzureRMAutomationSchedule_prerequisites(data), data.RandomInteger, monthDay)
}

func checkAccAzureRMAutomationSchedule_recurring_advanced_month(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists("azurerm_automation_schedule.test"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "Month"),
		resource.TestCheckResourceAttr(resourceName, "interval", "1"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
		resource.TestCheckResourceAttr(resourceName, "month_days.#", "1"),
	)
}

func testAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(data acceptance.TestData, weekDay string, weekDayOccurrence int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "Month"
  interval                = "1"

  monthly_occurrence {
    day        = "%s"
    occurrence = "%d"
  }
}
`, testAccAzureRMAutomationSchedule_prerequisites(data), data.RandomInteger, weekDay, weekDayOccurrence)
}

func checkAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(resourceName string, monthWeekDay string, monthWeekOccurrence int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists("azurerm_automation_schedule.test"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "Month"),
		resource.TestCheckResourceAttr(resourceName, "interval", "1"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
		resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.0.day", monthWeekDay),
		resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.0.occurrence", strconv.Itoa(monthWeekOccurrence)),
	)
}
