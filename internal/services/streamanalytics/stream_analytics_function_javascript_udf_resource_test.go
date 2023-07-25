// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsFunctionJavaScriptUDFResource struct{}

func TestAccStreamAnalyticsFunctionJavaScriptUDF_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	r := StreamAnalyticsFunctionJavaScriptUDFResource{}

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

func TestAccStreamAnalyticsFunctionJavaScriptUDF_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	r := StreamAnalyticsFunctionJavaScriptUDFResource{}

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

func TestAccStreamAnalyticsFunctionJavaScriptUDF_inputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	r := StreamAnalyticsFunctionJavaScriptUDFResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.inputs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsFunctionJavaScriptUDF_isConfigurationParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	r := StreamAnalyticsFunctionJavaScriptUDFResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.isConfigurationParameter(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.isConfigurationParameter(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StreamAnalyticsFunctionJavaScriptUDFResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := functions.ParseFunctionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.FunctionsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s : %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r StreamAnalyticsFunctionJavaScriptUDFResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  script = <<SCRIPT
function getRandomNumber(in) {
  return in;
}
SCRIPT


  input {
    type = "bigint"
  }

  output {
    type = "bigint"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsFunctionJavaScriptUDFResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "import" {
  name                      = azurerm_stream_analytics_function_javascript_udf.test.name
  stream_analytics_job_name = azurerm_stream_analytics_function_javascript_udf.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_function_javascript_udf.test.resource_group_name
  script                    = azurerm_stream_analytics_function_javascript_udf.test.script

  input {
    type = azurerm_stream_analytics_function_javascript_udf.test.input.0.type
  }

  output {
    type = azurerm_stream_analytics_function_javascript_udf.test.output.0.type
  }
}
`, template)
}

func (r StreamAnalyticsFunctionJavaScriptUDFResource) inputs(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  script = <<SCRIPT
function getRandomNumber(first, second) {
  return first * second;
}
SCRIPT


  input {
    type = "bigint"
  }

  input {
    type = "bigint"
  }

  output {
    type = "bigint"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsFunctionJavaScriptUDFResource) isConfigurationParameter(data acceptance.TestData, isConfigurationParameter bool) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  script = <<SCRIPT
function getRandomNumber(in) {
  return in;
}
SCRIPT


  input {
    type                    = "bigint"
    configuration_parameter = %t
  }

  output {
    type = "bigint"
  }
}
`, template, data.RandomInteger, isConfigurationParameter)
}

func (r StreamAnalyticsFunctionJavaScriptUDFResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%d"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
