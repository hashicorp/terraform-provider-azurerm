package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type LogicAppTriggerRecurrenceResource struct {
}

func TestAccLogicAppTriggerRecurrence_month(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Month", 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Month"),
				check.That(data.ResourceName).Key("interval").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Month", 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Month"),
				check.That(data.ResourceName).Key("interval").HasValue("1"),
			),
		},
		{
			Config:      r.requiresImport(data, "Month", 1),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_recurrence"),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_week(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Week", 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Week"),
				check.That(data.ResourceName).Key("interval").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Day", 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Day"),
				check.That(data.ResourceName).Key("interval").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_minute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Minute", 4),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Minute"),
				check.That(data.ResourceName).Key("interval").HasValue("4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_second(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Second", 30),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Second"),
				check.That(data.ResourceName).Key("interval").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_hour(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Hour", 4),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Hour"),
				check.That(data.ResourceName).Key("interval").HasValue("4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Month", 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Month"),
				check.That(data.ResourceName).Key("interval").HasValue("1"),
			),
		},
		{
			Config: r.basic(data, "Month", 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frequency").HasValue("Month"),
				check.That(data.ResourceName).Key("interval").HasValue("3"),
			),
		},
	})
}

func TestAccLogicAppTriggerRecurrence_startTime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.startTime(data, "2020-01-01T01:02:03Z"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_time").HasValue("2020-01-01T01:02:03Z"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_startTimeWithTimeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.startTimeWithTimeZone(data, "2020-01-01T01:02:03Z", "US Eastern Standard Time"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_time").HasValue("2020-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("time_zone").HasValue("US Eastern Standard Time"),
			),
		},
		data.ImportStep(),
		{
			Config: r.startTimeWithTimeZone(data, "2020-01-01T01:02:03Z", "Egypt Standard Time"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_time").HasValue("2020-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("time_zone").HasValue("Egypt Standard Time"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerRecurrence_schedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_recurrence", "test")
	r := LogicAppTriggerRecurrenceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.schedule(data, "Week", 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scheduleUpdated(data, "Week", 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Week", 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogicAppTriggerRecurrenceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	return triggerExists(ctx, clients, state)
}

func (LogicAppTriggerRecurrenceResource) basic(data acceptance.TestData, frequency string, interval int) string {
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

func (LogicAppTriggerRecurrenceResource) startTime(data acceptance.TestData, startTime string) string {
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

func (LogicAppTriggerRecurrenceResource) startTimeWithTimeZone(data acceptance.TestData, startTime string, timeZone string) string {
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

func (r LogicAppTriggerRecurrenceResource) requiresImport(data acceptance.TestData, frequency string, interval int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_recurrence" "import" {
  name         = azurerm_logic_app_trigger_recurrence.test.name
  logic_app_id = azurerm_logic_app_trigger_recurrence.test.logic_app_id
  frequency    = azurerm_logic_app_trigger_recurrence.test.frequency
  interval     = azurerm_logic_app_trigger_recurrence.test.interval
}
`, r.basic(data, frequency, interval))
}

func (LogicAppTriggerRecurrenceResource) schedule(data acceptance.TestData, frequency string, interval int) string {
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

  schedule {
    at_these_minutes = [10, 21]
    on_these_days    = ["Monday", "Friday"]
    at_these_hours   = [12, 15]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, frequency, interval)
}

func (LogicAppTriggerRecurrenceResource) scheduleUpdated(data acceptance.TestData, frequency string, interval int) string {
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

  schedule {
    at_these_hours = [10]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, frequency, interval)
}
