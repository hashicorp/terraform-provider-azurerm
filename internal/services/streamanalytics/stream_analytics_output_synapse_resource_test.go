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

type StreamAnalyticsOutputSynapseResource struct{}

func TestAccStreamAnalyticsOutputSynapse_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_synapse", "test")
	r := StreamAnalyticsOutputSynapseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccStreamAnalyticsOutputSynapse_sqlPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_synapse", "test")
	r := StreamAnalyticsOutputSynapseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sqlPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccStreamAnalyticsOutputSynapse_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_synapse", "test")
	r := StreamAnalyticsOutputSynapseResource{}

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
		data.ImportStep("password"),
	})
}

func TestAccStreamAnalyticsOutputSynapse_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_synapse", "test")
	r := StreamAnalyticsOutputSynapseResource{}

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

func (r StreamAnalyticsOutputSynapseResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r StreamAnalyticsOutputSynapseResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_stream_analytics_output_synapse" "test" {
  name                      = "acctestoutput-%[2]d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  server   = azurerm_synapse_workspace.test.connectivity_endpoints["sqlOnDemand"]
  user     = azurerm_synapse_workspace.test.sql_administrator_login
  password = azurerm_synapse_workspace.test.sql_administrator_login_password
  database = "master"
  table    = "AccTestTable"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputSynapseResource) sqlPool(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestSP%[3]s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
  storage_account_type = "GRS"
}

resource "azurerm_stream_analytics_output_synapse" "test" {
  name                      = "acctestoutput-%[2]d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  server   = azurerm_synapse_workspace.test.connectivity_endpoints["sql"]
  user     = azurerm_synapse_workspace.test.sql_administrator_login
  password = azurerm_synapse_workspace.test.sql_administrator_login_password
  database = azurerm_synapse_sql_pool.test.name
  table    = "AccTestTable"
}
`, template, data.RandomInteger, data.RandomString)
}

func (r StreamAnalyticsOutputSynapseResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_stream_analytics_output_synapse" "test" {
  name                      = "acctestoutput-%[2]d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  server   = azurerm_synapse_workspace.test.connectivity_endpoints["sqlOnDemand"]
  user     = azurerm_synapse_workspace.test.sql_administrator_login
  password = "updatedPassword"
  database = "master"
  table    = "AccTestTable"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputSynapseResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_stream_analytics_output_synapse" "import" {
  name                      = azurerm_stream_analytics_output_synapse.test.name
  stream_analytics_job_name = azurerm_stream_analytics_output_synapse.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_output_synapse.test.resource_group_name

  server   = azurerm_synapse_workspace.test.connectivity_endpoints["sqlOnDemand"]
  user     = azurerm_synapse_workspace.test.sql_administrator_login
  password = azurerm_synapse_workspace.test.sql_administrator_login_password
  database = "master"
  table    = "AccTestTable"
}
`, template)
}

func (r StreamAnalyticsOutputSynapseResource) template(data acceptance.TestData) string {
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
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctestdlfs%[3]s"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[3]s"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%[3]s"
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
