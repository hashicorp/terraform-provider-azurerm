package automation_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationScheduleResource struct {
}

func TestAccAzureRMAutomationSchedule_oneTime_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.oneTime_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.oneTime_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMAutomationSchedule_oneTime_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	// the API returns the time in the timezone we pass in
	// it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("Australia/Perth")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.oneTime_complete(data, startTime),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_oneTime_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	// the API returns the time in the timezone we pass in
	// it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("Australia/Perth")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.oneTime_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.oneTime_complete(data, startTime),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_hourly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.recurring_basic(data, "Hour", 7),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_daily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.recurring_basic(data, "Day", 7),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_weekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.recurring_basic(data, "Week", 7),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_monthly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.recurring_basic(data, "Month", 7),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_weekly_advanced(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.recurring_advanced_week(data, "Monday"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_monthly_advanced_by_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.recurring_advanced_month(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAutomationSchedule_monthly_advanced_by_week_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.recurring_advanced_month_week_day(data, "Monday", 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t AutomationScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	name := id.Path["schedules"]
	resGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := clients.Automation.ScheduleClient.Get(ctx, resGroup, accountName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Schedule '%s' (resource group: '%s') does not exist", name, id.ResourceGroup)
	}

	return utils.Bool(resp.ScheduleProperties != nil), nil
}

func (AutomationScheduleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
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

func (AutomationScheduleResource) oneTime_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
}
`, AutomationScheduleResource{}.template(data), data.RandomInteger)
}

func (AutomationScheduleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "import" {
  name                    = azurerm_automation_schedule.test.name
  resource_group_name     = azurerm_automation_schedule.test.resource_group_name
  automation_account_name = azurerm_automation_schedule.test.automation_account_name
  frequency               = azurerm_automation_schedule.test.frequency
}
`, AutomationScheduleResource{}.oneTime_basic(data))
}

func (AutomationScheduleResource) oneTime_complete(data acceptance.TestData, startTime string) string {
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
`, AutomationScheduleResource{}.template(data), data.RandomInteger, startTime)
}

// nolint unparam
func (AutomationScheduleResource) recurring_basic(data acceptance.TestData, frequency string, interval int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "%s"
  interval                = "%d"
}
`, AutomationScheduleResource{}.template(data), data.RandomInteger, frequency, interval)
}

func (AutomationScheduleResource) recurring_advanced_week(data acceptance.TestData, weekDay string) string {
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
`, AutomationScheduleResource{}.template(data), data.RandomInteger, weekDay)
}

func (AutomationScheduleResource) recurring_advanced_month(data acceptance.TestData, monthDay int) string {
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
`, AutomationScheduleResource{}.template(data), data.RandomInteger, monthDay)
}

func (AutomationScheduleResource) recurring_advanced_month_week_day(data acceptance.TestData, weekDay string, weekDayOccurrence int) string {
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
`, AutomationScheduleResource{}.template(data), data.RandomInteger, weekDay, weekDayOccurrence)
}
