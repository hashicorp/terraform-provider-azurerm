// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				check.That(data.ResourceName).Key("compatibility_level").HasValue("1.2"),
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

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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

func TestAccStreamAnalyticsJob_jobTypeCloud(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jobTypeCloud(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsJob_jobTypeEdge(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jobTypeEdge(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsJob_jobStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	r := StreamAnalyticsJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jobStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("job_storage_account.0.account_key"),
	})
}

func (r StreamAnalyticsJobResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := streamingjobs.ParseStreamingJobID(state.ID)
	if err != nil {
		return nil, err
	}

	var opts streamingjobs.GetOperationOptions
	resp, err := client.StreamAnalytics.JobsClient.Get(ctx, *id, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, err
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
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
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_stream_analytics_cluster" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  streaming_capacity  = 36
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%[1]d"
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = azurerm_resource_group.test.location
  data_locale                              = "en-GB"
  compatibility_level                      = "1.2"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3
  stream_analytics_cluster_id              = azurerm_stream_analytics_cluster.test.id
  tags = {
    environment = "Test"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary)
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

func (r StreamAnalyticsJobResource) jobTypeCloud(data acceptance.TestData) string {
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
  type                = "Cloud"

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StreamAnalyticsJobResource) jobTypeEdge(data acceptance.TestData) string {
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
  type                = "Edge"

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StreamAnalyticsJobResource) jobStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%[3]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_stream_analytics_job" "test" {
  name                   = "acctestjob-%[4]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  streaming_units        = 3
  content_storage_policy = "JobStorageAccount"
  job_storage_account {
    account_name = azurerm_storage_account.test.name
    account_key  = azurerm_storage_account.test.primary_access_key
  }

  tags = {
    environment = "Test"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.RandomString, data.Locations.Primary, data.RandomInteger)
}
