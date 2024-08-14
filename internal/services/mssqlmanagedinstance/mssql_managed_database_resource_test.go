// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type MsSqlManagedDatabase struct{}

func TestAccMsSqlManagedDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")
	r := MsSqlManagedDatabase{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(""),
	})
}

func TestAccMsSqlManagedDatabase_withRetentionPolicies(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "test")
	r := MsSqlManagedDatabase{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRetentionPolicies(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(""),
		{
			Config: r.withRetentionPoliciesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccMsSqlManagedDatabase_pointInTimeRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_database", "pitr")
	r := MsSqlManagedDatabase{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data), // Sets up for pitr
		},
		{
			PreConfig: func() { time.Sleep(11 * time.Minute) },
			// Valid point in time range from 7 days early to now and not before source server creation time
			Config: r.pointInTimeRestore(data, time.Now().Add(-15*time.Minute).UTC().Format(time.RFC3339)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(""),
	})
}

func (r MsSqlManagedDatabase) Exists(ctx context.Context, client *clients.Client, state *acceptance.InstanceState) (*bool, error) {
	id, err := commonids.ParseManagedInstanceDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQLManagedInstance.ManagedDatabasesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving SQL Managed Database %q: %+v", id.ID(), err)
	}

	return pointer.To(true), nil
}

func (r MsSqlManagedDatabase) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_database" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  name                = "acctest-%[2]d"
}
`, MsSqlManagedInstanceResource{}.basic(data), data.RandomInteger)
}

func (r MsSqlManagedDatabase) withRetentionPolicies(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_database" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  name                = "acctest-%[2]d"

  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 1
  }

  short_term_retention_days = 3

}
`, MsSqlManagedInstanceResource{}.basic(data), data.RandomInteger)
}

func (r MsSqlManagedDatabase) withRetentionPoliciesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_database" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  name                = "acctest-%[2]d"

  long_term_retention_policy {
    weekly_retention  = "P10D"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 4
  }

  short_term_retention_days = 4

}
`, MsSqlManagedInstanceResource{}.basic(data), data.RandomInteger)
}

func (r MsSqlManagedDatabase) pointInTimeRestore(data acceptance.TestData, restorePointInTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_database" "pitr" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  name                = "acctest-%[2]d-pitr"

  point_in_time_restore {
    restore_point_in_time = "%[3]s"
    source_database_id    = azurerm_mssql_managed_database.test.id
  }
}
`, MsSqlManagedDatabase{}.basic(data), data.RandomInteger, restorePointInTime)
}
