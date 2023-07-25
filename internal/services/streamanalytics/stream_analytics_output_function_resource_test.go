// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsOutputFunctionResource struct{}

func TestAccStreamAnalyticsOutputFunction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_function", "test")
	r := StreamAnalyticsOutputFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccStreamAnalyticsOutputFunction_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_function", "test")
	r := StreamAnalyticsOutputFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccStreamAnalyticsOutputFunction_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_function", "test")
	r := StreamAnalyticsOutputFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccStreamAnalyticsOutputFunction_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_function", "test")
	r := StreamAnalyticsOutputFunctionResource{}

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

func (r StreamAnalyticsOutputFunctionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := outputs.ParseOutputID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.OutputsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsOutputFunctionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_function" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  function_app              = azurerm_function_app.test.name
  function_name             = "somefunctionname"
  api_key                   = "test"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputFunctionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_function" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  function_app              = azurerm_function_app.test.name
  function_name             = "somefunctionname"
  api_key                   = "test"
  batch_max_in_bytes        = 128
  batch_max_count           = 200
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputFunctionResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_function" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  function_app              = azurerm_function_app.test.name
  function_name             = "adifferentfunctionname"
  api_key                   = "withanewkey!"
  batch_max_in_bytes        = 128
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputFunctionResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_function" "import" {
  name                      = azurerm_stream_analytics_output_function.test.name
  stream_analytics_job_name = azurerm_stream_analytics_output_function.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_output_function.test.resource_group_name
  function_app              = azurerm_stream_analytics_output_function.test.function_app
  function_name             = azurerm_stream_analytics_output_function.test.function_name
  api_key                   = azurerm_stream_analytics_output_function.test.api_key
}
`, template)
}

func (r StreamAnalyticsOutputFunctionResource) template(data acceptance.TestData) string {
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

resource "azurerm_app_service_plan" "test" {
  name                = "acctestplan-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "FunctionApp"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctestfunction-%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  os_type                    = "linux"
  version                    = "~3"
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%[1]d"
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = azurerm_resource_group.test.location
  compatibility_level                      = "1.0"
  data_locale                              = "en-GB"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
