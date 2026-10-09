// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type sqlManagedInstanceStartStopScheduleResource struct{}

func TestAccMsSqlManagedInstanceStartStopSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_start_stop_schedule", "test")
	r := sqlManagedInstanceStartStopScheduleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlManagedInstanceStartStopSchedule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_start_stop_schedule", "test")
	r := sqlManagedInstanceStartStopScheduleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlManagedInstanceStartStopSchedule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_start_stop_schedule", "test")
	r := sqlManagedInstanceStartStopScheduleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.importStep(data),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.importStep(data),
	})
}

func TestAccMsSqlManagedInstanceStartStopSchedule_multipleSchedules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_start_stop_schedule", "test")
	r := sqlManagedInstanceStartStopScheduleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleSchedules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("next_execution_time").Exists(),
				check.That(data.ResourceName).Key("next_run_action").Exists(),
			),
		},
		r.importStep(data),
		{
			Config: r.multipleSchedulesUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("next_execution_time").Exists(),
				check.That(data.ResourceName).Key("next_run_action").Exists(),
			),
		},
		r.importStep(data),
		{
			Config: r.multipleSchedules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("next_execution_time").Exists(),
				check.That(data.ResourceName).Key("next_run_action").Exists(),
			),
		},
		r.importStep(data),
	})
}

func TestAccMsSqlManagedInstanceStartStopSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_start_stop_schedule", "test")
	r := sqlManagedInstanceStartStopScheduleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r sqlManagedInstanceStartStopScheduleResource) importStep(data acceptance.TestData) acceptance.TestStep {
	step := data.ImportStep()
	if features.SixPointOh() {
		return step
	}

	// Import has no existing list order in 5.x. Compare all schedule values and
	// their counts separately, while retaining normal verification for other fields.
	step.ImportStateVerifyIgnore = []string{"schedule"}
	schedules := func(state *pluginsdk.InstanceState) (map[[4]string]int, error) {
		count, err := strconv.Atoi(state.Attributes["schedule.#"])
		if err != nil {
			return nil, fmt.Errorf("reading schedule count: %+v", err)
		}

		items := make(map[[4]string]int, count)
		for i := range count {
			prefix := fmt.Sprintf("schedule.%d.", i)
			items[[4]string{
				state.Attributes[prefix+"start_day"],
				state.Attributes[prefix+"start_time"],
				state.Attributes[prefix+"stop_day"],
				state.Attributes[prefix+"stop_time"],
			}]++
		}
		return items, nil
	}

	var expected map[[4]string]int
	step.ImportStateIdFunc = func(state *pluginsdk.State) (string, error) {
		resource, ok := state.RootModule().Resources[data.ResourceName]
		if !ok || resource.Primary == nil {
			return "", fmt.Errorf("%s was not found in state", data.ResourceName)
		}
		var err error
		expected, err = schedules(resource.Primary)
		return resource.Primary.ID, err
	}
	step.ImportStateCheck = func(states []*pluginsdk.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected one imported resource, got %d", len(states))
		}
		actual, err := schedules(states[0])
		if err != nil {
			return err
		}
		if !maps.Equal(actual, expected) {
			return fmt.Errorf("expected imported schedules %#v, got %#v", expected, actual)
		}
		return nil
	}
	return step
}

func (r sqlManagedInstanceStartStopScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedInstanceStartStopScheduleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MSSQLManagedInstance.ManagedInstanceStartStopSchedulesClient

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	resp, err := client.Get(ctx, managedInstanceId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r sqlManagedInstanceStartStopScheduleResource) template(data acceptance.TestData) string {
	return MsSqlManagedInstanceResource{}.basic(data)
}

func (r sqlManagedInstanceStartStopScheduleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_managed_instance_start_stop_schedule" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  schedule {
    start_day  = "Wednesday"
    start_time = "11:00"
    stop_day   = "Wednesday"
    stop_time  = "23:00"
  }
}
`, r.template(data))
}

func (r sqlManagedInstanceStartStopScheduleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_managed_instance_start_stop_schedule" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  description         = "test description"
  timezone_id         = "Central European Standard Time"
  schedule {
    start_day  = "Wednesday"
    start_time = "11:00"
    stop_day   = "Wednesday"
    stop_time  = "23:00"
  }
}
`, r.template(data))
}

func (r sqlManagedInstanceStartStopScheduleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_managed_instance_start_stop_schedule" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  description         = "updated test description"
  timezone_id         = "Central European Standard Time"
  schedule {
    start_day  = "Wednesday"
    start_time = "10:00"
    stop_day   = "Wednesday"
    stop_time  = "22:00"
  }
  schedule {
    start_day  = "Thursday"
    start_time = "11:00"
    stop_day   = "Thursday"
    stop_time  = "23:00"
  }
}
`, r.template(data))
}

func (r sqlManagedInstanceStartStopScheduleResource) multipleSchedules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_managed_instance_start_stop_schedule" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  description         = "test description"
  timezone_id         = "Central European Standard Time"
  schedule {
    start_day  = "Monday"
    start_time = "08:00"
    stop_day   = "Monday"
    stop_time  = "20:00"
  }

  schedule {
    start_day  = "Friday"
    start_time = "09:00"
    stop_day   = "Friday"
    stop_time  = "21:00"
  }

  schedule {
    start_day  = "Wednesday"
    start_time = "11:00"
    stop_day   = "Wednesday"
    stop_time  = "23:00"
  }
}
`, r.template(data))
}

func (r sqlManagedInstanceStartStopScheduleResource) multipleSchedulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_managed_instance_start_stop_schedule" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  description         = "updated test description"
  timezone_id         = "Central European Standard Time"
  schedule {
    start_day  = "Monday"
    start_time = "08:00"
    stop_day   = "Monday"
    stop_time  = "20:00"
  }
  schedule {
    start_day  = "Friday"
    start_time = "09:00"
    stop_day   = "Friday"
    stop_time  = "21:00"
  }
  schedule {
    start_day  = "Wednesday"
    start_time = "10:00"
    stop_day   = "Wednesday"
    stop_time  = "22:00"
  }
  schedule {
    start_day  = "Thursday"
    start_time = "11:00"
    stop_day   = "Thursday"
    stop_time  = "23:00"
  }
}
`, r.template(data))
}

func (r sqlManagedInstanceStartStopScheduleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_managed_instance_start_stop_schedule" "import" {
  managed_instance_id = azurerm_mssql_managed_instance_start_stop_schedule.test.managed_instance_id
  schedule {
    start_day  = "Wednesday"
    start_time = "11:00"
    stop_day   = "Wednesday"
    stop_time  = "23:00"
  }
}
`, r.basic(data))
}
