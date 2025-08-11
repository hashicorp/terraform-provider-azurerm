// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsJobStorageAccountResource struct{}

func TestAccStreamAnalyticsJobStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job_storage_account", "test")
	r := StreamAnalyticsJobStorageAccountResource{}

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

func TestAccStreamAnalyticsJobStorageAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job_storage_account", "test")
	r := StreamAnalyticsJobStorageAccountResource{}

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

func TestAccStreamAnalyticsJobStorageAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job_storage_account", "test")
	r := StreamAnalyticsJobStorageAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StreamAnalyticsJobStorageAccountResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := streamingjobs.ParseStreamingJobID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.JobsClient.Get(ctx, *id, streamingjobs.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.JobStorageAccount == nil {
		return pointer.To(false), nil
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StreamAnalyticsJobStorageAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job_storage_account" "test" {
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id
  storage_account_name    = azurerm_storage_account.test.name
  authentication_mode     = "Msi"
}
`, r.template(data))
}

func (r StreamAnalyticsJobStorageAccountResource) connectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job_storage_account" "test" {
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id
  storage_account_name    = azurerm_storage_account.test.name
  storage_account_key     = azurerm_storage_account.test.primary_access_key
  authentication_mode     = "ConnectionString"
}
`, r.template(data))
}

func (r StreamAnalyticsJobStorageAccountResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job_storage_account" "import" {
  stream_analytics_job_id = azurerm_stream_analytics_job_storage_account.test.stream_analytics_job_id
  storage_account_name    = azurerm_stream_analytics_job_storage_account.test.storage_account_name
  authentication_mode     = azurerm_stream_analytics_job_storage_account.test.authentication_mode
}
`, r.basic(data))
}

func (r StreamAnalyticsJobStorageAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_stream_analytics_job" "test" {
  name                = "acctestjob-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  streaming_units     = 3

  tags = {
    environment = "Test"
  }

  identity {
    type = "SystemAssigned"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY


  lifecycle {
    ignore_changes = [job_storage_account]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
