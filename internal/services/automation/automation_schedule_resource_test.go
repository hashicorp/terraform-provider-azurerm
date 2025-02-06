// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationScheduleResource struct{}

func TestAccAutomationSchedule_oneTime_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneTime_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// test for: https://github.com/hashicorp/terraform-provider-azurerm/issues/21854
func TestAccAutomationSchedule_expiryTimeOfEuropeTimeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expiryTimeOfEuropeTimeZone(data, "foo"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.expiryTimeOfEuropeTimeZone(data, "bar"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneTime_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutomationSchedule_oneTime_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	// the API returns the time in the timezone we pass in
	// it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("Australia/Perth")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneTime_complete(data, startTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_oneTime_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	// the API returns the time in the timezone we pass in
	// it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("Australia/Perth")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneTime_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.oneTime_complete(data, startTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_hourly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurring_basic(data, "Hour", 7),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_daily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurring_basic(data, "Day", 7),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_weekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurring_basic(data, "Week", 7),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_monthly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurring_basic(data, "Month", 7),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_weekly_advanced(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurring_advanced_week(data, "Monday"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_monthly_advanced_by_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurring_advanced_month(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationSchedule_monthly_advanced_by_week_day(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_schedule", "test")
	r := AutomationScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.recurring_advanced_month_week_day(data, "Monday", 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t AutomationScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := schedule.ParseScheduleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Automation.Schedule.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
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

func (a AutomationScheduleResource) expiryTimeOfEuropeTimeZone(data acceptance.TestData, desc string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "Week"
  interval                = 1
  timezone                = "Europe/Amsterdam"
  start_time              = "2026-04-15T18:01:15+02:00"
  description             = "%s"
  week_days               = ["Monday"]
}
`, a.template(data), data.RandomInteger, desc)
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
