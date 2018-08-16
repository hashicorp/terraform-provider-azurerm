package azurerm

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationSchedule_oneTimeBasic(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTimeBasic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_requiresImport(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTimeBasic(ri, location),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAutomationSchedule_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_automation_schedule"),
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_oneTimeComplete(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	//the API returns the time in the timezone we pass in
	//it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("CET")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTimeComplete(ri, testLocation(), startTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
					resource.TestCheckResourceAttr(resourceName, "start_time", startTime),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Central Europe Standard Time"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is an automation schedule"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_oneTimeUpdate(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	//the API returns the time in the timezone we pass in
	//it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("CET")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTimeBasic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
				),
			},
			{
				Config: testAccAzureRMAutomationSchedule_oneTimeComplete(ri, testLocation(), startTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
					resource.TestCheckResourceAttr(resourceName, "start_time", startTime),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Central Europe Standard Time"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is an automation schedule"),
				),
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_recurringHourly(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurringBasic(ri, testLocation(), "Hour", 7),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "Hour"),
					resource.TestCheckResourceAttr(resourceName, "interval", "7"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_recurringDaily(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurringBasic(ri, testLocation(), "Day", 7),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "Day"),
					resource.TestCheckResourceAttr(resourceName, "interval", "7"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_recurringWeekly(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurringBasic(ri, testLocation(), "Week", 7),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "Week"),
					resource.TestCheckResourceAttr(resourceName, "interval", "7"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_recurringMonthly(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurringBasic(ri, testLocation(), "Month", 7),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "Month"),
					resource.TestCheckResourceAttr(resourceName, "interval", "7"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_weeklyAdvanced(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurringAdvancedWeek(ri, testLocation(), "Monday"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "Week"),
					resource.TestCheckResourceAttr(resourceName, "interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "week_days.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthlyAdvancedByDay(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurringAdvancedMonth(ri, testLocation(), 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "Month"),
					resource.TestCheckResourceAttr(resourceName, "interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "month_days.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthlyAdvancedByWeekday(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurringAdvancedMonthWeekDay(ri, testLocation(), "Monday", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMAutomationScheduleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "Month"),
					resource.TestCheckResourceAttr(resourceName, "interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.0.day", "Monday"),
					resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.0.occurrence", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMAutomationScheduleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationScheduleClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_schedule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]

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

func testCheckAzureRMAutomationScheduleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*ArmClient).automationScheduleClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
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

func testAccAzureRMAutomationSchedule_oneTimeBasic(rInt int, location string) string {
	template := testAccAzureRMAutomationSchedule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "OneTime"
}
`, template, rInt)
}

func testAccAzureRMAutomationSchedule_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAutomationSchedule_oneTimeBasic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "import" {
  name                    = "${azurerm_automation_schedule.test.name}"
  resource_group_name     = "${azurerm_automation_schedule.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_schedule.test.automation_account_name}"
  frequency               = "OneTime"
}
`, template)
}

func testAccAzureRMAutomationSchedule_oneTimeComplete(rInt int, location, startTime string) string {
	template := testAccAzureRMAutomationSchedule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "OneTime"
  start_time              = "%s"
  timezone                = "Central Europe Standard Time" 
  description             = "This is an automation schedule"
}
`, template, rInt, startTime)
}

func testAccAzureRMAutomationSchedule_recurringBasic(rInt int, location, frequency string, interval int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "%s"
  interval                = "%d"
}
`, testAccAzureRMAutomationSchedule_template(rInt, location), rInt, frequency, interval)
}

func testAccAzureRMAutomationSchedule_recurringAdvancedWeek(rInt int, location string, weekDay string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "Week"
  interval                = "1"
  week_days               = ["%s"]
}	
`, testAccAzureRMAutomationSchedule_template(rInt, location), rInt, weekDay)
}

func testAccAzureRMAutomationSchedule_recurringAdvancedMonth(rInt int, location string, monthDay int) string {
	template := testAccAzureRMAutomationSchedule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "Month"
  interval                = "1"
  month_days              = [%d]
}	
`, template, rInt, monthDay)
}

func testAccAzureRMAutomationSchedule_recurringAdvancedMonthWeekDay(rInt int, location string, weekDay string, weekDayOccurrence int) string {
	template := testAccAzureRMAutomationSchedule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "Month"
  interval                = "1"

  monthly_occurrence {
	day        = "%s"
	occurrence = "%d"
  }
}	
`, template, rInt, weekDay, weekDayOccurrence)
}

func testAccAzureRMAutomationSchedule_template(rInt int, location string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" { 
  name     = "acctestRG-%d" 
  location = "%s" 
} 
 
resource "azurerm_automation_account" "test" { 
  name                = "acctestAA-%d" 
  location            = "${azurerm_resource_group.test.location}" 
  resource_group_name = "${azurerm_resource_group.test.name}" 
  sku { 
    name = "Basic" 
  } 
}
`, rInt, location, rInt)
}
