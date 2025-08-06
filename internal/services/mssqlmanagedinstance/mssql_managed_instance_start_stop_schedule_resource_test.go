// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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

func (r sqlManagedInstanceStartStopScheduleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_managed_instance_start_stop_schedule" "import" {
  managed_instance_id = azurerm_mssql_managed_instance_start_stop_schedule.test.managed_instance_id
  schedule {
    start_day  = azurerm_mssql_managed_instance_start_stop_schedule.test.schedule.0.start_day
    start_time = azurerm_mssql_managed_instance_start_stop_schedule.test.schedule.0.start_time
    stop_day   = azurerm_mssql_managed_instance_start_stop_schedule.test.schedule.0.stop_day
    stop_time  = azurerm_mssql_managed_instance_start_stop_schedule.test.schedule.0.stop_time
  }
}
`, r.basic(data))
}
