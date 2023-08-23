// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypackqueries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsQueryPackQueryResource struct{ uuid string }

func (r LogAnalyticsQueryPackQueryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := querypackqueries.ParseQueryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.LogAnalytics.QueryPackQueriesClient.QueriesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccLogAnalyticsQueryPackQuery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_query_pack_query", "test")
	r := LogAnalyticsQueryPackQueryResource{}

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

func TestAccLogAnalyticsQueryPackQuery_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_query_pack_query", "test")
	r := LogAnalyticsQueryPackQueryResource{}

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

func TestAccLogAnalyticsQueryPackQuery_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_query_pack_query", "test")
	r := LogAnalyticsQueryPackQueryResource{uuid: uuid.New().String()}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsQueryPackQuery_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_query_pack_query", "test")
	r := LogAnalyticsQueryPackQueryResource{uuid: uuid.New().String()}

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
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogAnalyticsQueryPackQueryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_query_pack" "test" {
  name                = "acctestlaqp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_log_analytics_query_pack_query" "test" {
  query_pack_id = azurerm_log_analytics_query_pack.test.id
  display_name  = "Exceptions - New in the last 24 hours"

  body = <<BODY
    let newExceptionsTimeRange = 1d;
    let timeRangeToCheckBefore = 7d;
    exceptions
    | where timestamp < ago(timeRangeToCheckBefore)
    | summarize count() by problemId
    | join kind= rightanti (
        exceptions
        | where timestamp >= ago(newExceptionsTimeRange)
        | extend stack = tostring(details[0].rawStack)
        | summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId
    ) on problemId
    | order by count_ desc
  BODY
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsQueryPackQueryResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_query_pack" "test" {
  name                = "acctestlaqp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_log_analytics_query_pack_query" "test" {
  name           = "%[3]s"
  query_pack_id  = azurerm_log_analytics_query_pack.test.id
  display_name   = "Exceptions - New in the last 24 hours"
  description    = "my description"
  categories     = ["network"]
  resource_types = ["microsoft.web/sites"]
  solutions      = ["LogManagement"]

  body = <<BODY
    let newExceptionsTimeRange = 1d;
    let timeRangeToCheckBefore = 7d;
    exceptions
    | where timestamp < ago(timeRangeToCheckBefore)
    | summarize count() by problemId
    | join kind= rightanti (
        exceptions
        | where timestamp >= ago(newExceptionsTimeRange)
        | extend stack = tostring(details[0].rawStack)
        | summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId
    ) on problemId
    | order by count_ desc
  BODY

  additional_settings_json = <<JSON
{
  "Environment": "Test"
}
JSON

  tags = {
    my-label       = "label1,label2"
    my-other-label = "label3,label4"
  }
}
`, r.template(data), data.RandomInteger, r.uuid)
}

func (r LogAnalyticsQueryPackQueryResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_query_pack" "test" {
  name                = "acctestlaqp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_log_analytics_query_pack_query" "test" {
  name           = "%[3]s"
  query_pack_id  = azurerm_log_analytics_query_pack.test.id
  display_name   = "Exceptions - New in the last 48 hours"
  description    = "my test description"
  categories     = ["resources"]
  resource_types = ["microsoft.network/virtualnetworks"]
  solutions      = ["NetworkMonitoring"]

  body = <<BODY
    let newExceptionsTimeRange = 2d;
    let timeRangeToCheckBefore = 7d;
    exceptions
    | where timestamp < ago(timeRangeToCheckBefore)
    | summarize count() by problemId
    | join kind= rightanti (
        exceptions
        | where timestamp >= ago(newExceptionsTimeRange)
        | extend stack = tostring(details[0].rawStack)
        | summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId
    ) on problemId
    | order by count_ desc
  BODY

  additional_settings_json = <<JSON
{
  "Environment2": "Test2"
}
JSON

  tags = {
    my-label       = "label5"
    my-other-label = "label7"
  }
}
`, r.template(data), data.RandomInteger, r.uuid)
}

func (r LogAnalyticsQueryPackQueryResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_query_pack_query" "import" {
  name          = azurerm_log_analytics_query_pack_query.test.name
  query_pack_id = azurerm_log_analytics_query_pack_query.test.query_pack_id
  body          = azurerm_log_analytics_query_pack_query.test.body
  display_name  = azurerm_log_analytics_query_pack_query.test.display_name
}
`, r.basic(data))
}

func (r LogAnalyticsQueryPackQueryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-LA-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
