package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StreamAnalyticsJobResource struct{}

func TestAccStreamAnalyticsJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

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

func TestAccStreamAnalyticsJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsJob_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r StreamAnalyticsJobResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := client.StreamAnalytics.JobsClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), err
		}
		return nil, fmt.Errorf("retrieving Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_stream_analytics_job" "test" {
  name                = "acctestjob-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  streaming_units     = 3

  tags = {
    environment = "Test"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StreamAnalyticsJobResource) complete(data acceptance.TestData) string {
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
  data_locale                              = "en-GB"
  compatibility_level                      = "1.0"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3

  tags = {
    environment = "Test"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StreamAnalyticsJobResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job" "import" {
  name                                     = azurerm_stream_analytics_job.test.name
  resource_group_name                      = azurerm_stream_analytics_job.test.resource_group_name
  location                                 = azurerm_stream_analytics_job.test.location
  compatibility_level                      = azurerm_stream_analytics_job.test.compatibility_level
  data_locale                              = azurerm_stream_analytics_job.test.data_locale
  events_late_arrival_max_delay_in_seconds = azurerm_stream_analytics_job.test.events_late_arrival_max_delay_in_seconds
  events_out_of_order_max_delay_in_seconds = azurerm_stream_analytics_job.test.events_out_of_order_max_delay_in_seconds
  events_out_of_order_policy               = azurerm_stream_analytics_job.test.events_out_of_order_policy
  output_error_policy                      = azurerm_stream_analytics_job.test.output_error_policy
  streaming_units                          = azurerm_stream_analytics_job.test.streaming_units
  transformation_query                     = azurerm_stream_analytics_job.test.transformation_query
  tags                                     = azurerm_stream_analytics_job.test.tags
}
`, template)
}

func (r StreamAnalyticsJobResource) updated(data acceptance.TestData) string {
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
  data_locale                              = "en-GB"
  compatibility_level                      = "1.1"
  events_late_arrival_max_delay_in_seconds = 10
  events_out_of_order_max_delay_in_seconds = 20
  events_out_of_order_policy               = "Drop"
  output_error_policy                      = "Stop"
  streaming_units                          = 6

  transformation_query = <<QUERY
    SELECT *
    INTO [SomeOtherOutputAlias]
    FROM [SomeOtherInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StreamAnalyticsJobResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_stream_analytics_job" "test" {
  name                = "acctestjob-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  streaming_units     = 3

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
