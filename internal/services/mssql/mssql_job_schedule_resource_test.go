// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlJobScheduleResourceTest struct{}

func TestAccMsSqlJobScheduleResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_schedule", "test")
	r := MsSqlJobScheduleResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Once"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlJobScheduleResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_schedule", "test")
	r := MsSqlJobScheduleResourceTest{}

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

func TestAccMsSqlJobScheduleResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_schedule", "test")
	r := MsSqlJobScheduleResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Once"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("interval").HasValue("PT15M"),
				check.That(data.ResourceName).Key("type").HasValue("Recurring"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlJobScheduleResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_schedule", "test")
	r := MsSqlJobScheduleResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("interval").HasValue("PT15M"),
				check.That(data.ResourceName).Key("type").HasValue("Recurring"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlJobScheduleResource_recurringWithoutInterval(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_schedule", "test")
	r := MsSqlJobScheduleResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.recurringWithoutInterval(data),
			ExpectError: regexp.MustCompile("`interval` must be set when `type` is `Recurring`"),
		},
	})
}

func (MsSqlJobScheduleResourceTest) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := jobs.ParseJobID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.JobsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (MsSqlJobScheduleResourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_schedule" "test" {
  enabled = true
  job_id  = azurerm_mssql_job.test.id
  type    = "Once"
}
`, MsSqlJobResourceTest{}.basic(data))
}

func (r MsSqlJobScheduleResourceTest) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_schedule" "import" {
  job_id = azurerm_mssql_job_schedule.test.job_id
  type   = azurerm_mssql_job_schedule.test.type
}
`, r.basic(data))
}

func (r MsSqlJobScheduleResourceTest) complete(data acceptance.TestData) string {
	now := time.Now()
	startTime := now.AddDate(0, 0, 5)
	endTime := now.AddDate(0, 0, 10)

	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_schedule" "test" {
  enabled    = true
  end_time   = "%[2]s"
  interval   = "PT15M"
  job_id     = azurerm_mssql_job.test.id
  start_time = "%[3]s"
  type       = "Recurring"
}
`, MsSqlJobResourceTest{}.basic(data), endTime.Format(time.RFC3339), startTime.Format(time.RFC3339))
}

func (MsSqlJobScheduleResourceTest) recurringWithoutInterval(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_schedule" "test" {
  enabled = true
  job_id  = azurerm_mssql_job.test.id
  type    = "Recurring"
}
`, MsSqlJobResourceTest{}.basic(data))
}
