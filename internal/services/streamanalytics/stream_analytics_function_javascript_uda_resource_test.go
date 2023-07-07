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

type StreamAnalyticsFunctionJavaScriptUDAResource struct{}

func TestAccStreamAnalyticsFunctionJavaScriptUDA_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_uda", "test")
	r := StreamAnalyticsFunctionJavaScriptUDAResource{}

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

func TestAccStreamAnalyticsFunctionJavaScriptUDA_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_uda", "test")
	r := StreamAnalyticsFunctionJavaScriptUDAResource{}

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

func TestAccStreamAnalyticsFunctionJavaScriptUDA_inputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_uda", "test")
	r := StreamAnalyticsFunctionJavaScriptUDAResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsFunctionJavaScriptUDA_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_uda", "test")
	r := StreamAnalyticsFunctionJavaScriptUDAResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func TestAccStreamAnalyticsFunctionJavaScriptUDA_isConfigurationParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_uda", "test")
	r := StreamAnalyticsFunctionJavaScriptUDAResource{}

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

func (r StreamAnalyticsFunctionJavaScriptUDAResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r StreamAnalyticsFunctionJavaScriptUDAResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_uda" "test" {
  name                    = "acctestinput-%d"
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id

  script = <<SCRIPT
function main() {
    this.init = function () {
        this.state = 0;
    }

    this.accumulate = function (value, timestamp) {
        this.state += value;
    }

    this.computeResult = function () {
        return this.state;
    }
}
SCRIPT


  input {
    type = "bigint"
  }

  output {
    type = "bigint"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r StreamAnalyticsFunctionJavaScriptUDAResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_uda" "import" {
  name                    = azurerm_stream_analytics_function_javascript_uda.test.name
  stream_analytics_job_id = azurerm_stream_analytics_function_javascript_uda.test.stream_analytics_job_id
  script                  = azurerm_stream_analytics_function_javascript_uda.test.script

  input {
    type = azurerm_stream_analytics_function_javascript_uda.test.input.0.type
  }

  output {
    type = azurerm_stream_analytics_function_javascript_uda.test.output.0.type
  }
}
`, r.basic(data))
}

func (r StreamAnalyticsFunctionJavaScriptUDAResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_uda" "test" {
  name                    = "acctestinput-%d"
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id

  script = <<SCRIPT
function main() {
    this.init = function () {
        this.state = 0;
    }

    this.accumulate = function (value, timestamp) {
        this.state += value;
    }

    this.computeResult = function () {
        return this.state;
    }
}
SCRIPT


  input {
    type = "bigint"
  }

  input {
    type = "float"
  }

  output {
    type = "bigint"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r StreamAnalyticsFunctionJavaScriptUDAResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_uda" "test" {
  name                    = "acctestinput-%d"
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id

  script = <<SCRIPT
function main() {
    this.init = function () {
        this.state = 0;
    }

    this.accumulate = function (value) {
        this.state += value;
    }

    this.computeResult = function () {
        return this.state;
    }
}
SCRIPT


  input {
    type = "float"
  }

  output {
    type = "float"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r StreamAnalyticsFunctionJavaScriptUDAResource) isConfigurationParameter(data acceptance.TestData, isConfigurationParameter bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_uda" "test" {
  name                    = "acctestinput-%d"
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id

  script = <<SCRIPT
function main() {
    this.init = function () {
        this.state = 0;
    }

    this.accumulate = function (value, timestamp) {
        this.state += value;
    }

    this.computeResult = function () {
        return this.state;
    }
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
`, r.template(data), data.RandomInteger, isConfigurationParameter)
}

func (r StreamAnalyticsFunctionJavaScriptUDAResource) template(data acceptance.TestData) string {
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
