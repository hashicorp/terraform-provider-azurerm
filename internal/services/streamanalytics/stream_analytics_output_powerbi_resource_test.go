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

type StreamAnalyticsOutputPowerBIResource struct{}

func (r StreamAnalyticsOutputPowerBIResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func TestAccStreamAnalyticsOutputPowerBI_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_powerbi", "test")
	r := StreamAnalyticsOutputPowerBIResource{}

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

func TestAccStreamAnalyticsOutputPowerBI_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_powerbi", "test")
	r := StreamAnalyticsOutputPowerBIResource{}

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

func TestAccStreamAnalyticsOutputPowerBI_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_powerbi", "test")
	r := StreamAnalyticsOutputPowerBIResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r StreamAnalyticsOutputPowerBIResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_powerbi" "test" {
  name                    = "acctestoutput-%d"
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id
  dataset                 = "foo"
  table                   = "bar"
  group_id                = "85b3dbca-5974-4067-9669-67a141095a76"
  group_name              = "some-test-group-name"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputPowerBIResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_powerbi" "test" {
  name                    = "acctestoutput-%d"
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id
  dataset                 = "updated-dataset"
  table                   = "updated-table"
  group_id                = "e18ff5df-fb66-4f6d-8f27-88c4dcbfc002"
  group_name              = "some-updated-group-id"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputPowerBIResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_powerbi" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_id   = azurerm_stream_analytics_job.test.id
  dataset                   = "complete-dataset"
  table                     = "complete-table"
  group_id                  = "e18ff5df-fb66-4f6d-8f27-88c4dcbfc002"
  group_name                = "some-test-group-name"
  token_user_principal_name = "bobsmith@contoso.com"
  token_user_display_name   = "Bob Smith"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputPowerBIResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
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
`, data.RandomInteger, data.Locations.Primary)
}
