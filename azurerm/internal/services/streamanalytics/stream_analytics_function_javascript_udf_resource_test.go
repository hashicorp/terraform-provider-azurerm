package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StreamAnalyticsFunctionJavaScriptUDFResource struct{}

func TestAccStreamAnalyticsFunctionJavaScriptUDF_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	r := StreamAnalyticsFunctionJavaScriptUDFResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsFunctionJavaScriptUDF_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	r := StreamAnalyticsFunctionJavaScriptUDFResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStreamAnalyticsFunctionJavaScriptUDF_inputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	r := StreamAnalyticsFunctionJavaScriptUDFResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.inputs(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StreamAnalyticsFunctionJavaScriptUDFResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	jobName := state.Attributes["stream_analytics_job_name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := client.StreamAnalytics.FunctionsClient.Get(ctx, resourceGroup, jobName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Function JavaScript UDF %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
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
