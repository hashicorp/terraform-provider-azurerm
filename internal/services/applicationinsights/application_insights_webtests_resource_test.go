// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	webtests "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AppInsightsWebTestsResource struct{}

func TestAccApplicationInsightsWebTests_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_web_test", "test")
	r := AppInsightsWebTestsResource{}

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

func TestAccApplicationInsightsWebTests_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_web_test", "test")
	r := AppInsightsWebTestsResource{}

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

func TestAccApplicationInsightsWebTests_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_web_test", "test")
	r := AppInsightsWebTestsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_locations.#").HasValue("1"),
				check.That(data.ResourceName).Key("frequency").HasValue("300"),
				check.That(data.ResourceName).Key("timeout").HasValue("30"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_locations.#").HasValue("2"),
				check.That(data.ResourceName).Key("frequency").HasValue("900"),
				check.That(data.ResourceName).Key("timeout").HasValue("120"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_locations.#").HasValue("1"),
				check.That(data.ResourceName).Key("frequency").HasValue("300"),
				check.That(data.ResourceName).Key("timeout").HasValue("30"),
			),
		},
	})
}

func TestAccApplicationInsightsWebTests_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_web_test", "test")
	r := AppInsightsWebTestsResource{}

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

func (t AppInsightsWebTestsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webtests.ParseWebTestID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppInsights.WebTestsClient.WebTestsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (AppInsightsWebTestsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_web_test" "test" {
  name                    = "acctestappinsightswebtests-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  kind                    = "ping"
  geo_locations           = ["us-tx-sn1-azr"]

  configuration = <<XML
<WebTest Name="WebTest1" Id="ABD48585-0831-40CB-9069-682EA6BB3583" Enabled="True" CssProjectStructure="" CssIteration="" Timeout="0" WorkItemIds="" xmlns="http://microsoft.com/schemas/VisualStudio/TeamTest/2010" Description="" CredentialUserName="" CredentialPassword="" PreAuthenticate="True" Proxy="default" StopOnError="False" RecordedResultFile="" ResultsLocale="">
  <Items>
    <Request Method="GET" Guid="a5f10126-e4cd-570d-961c-cea43999a200" Version="1.1" Url="http://microsoft.com" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="200" ExpectedResponseUrl="" ReportingName="" IgnoreHttpStatusCode="False" />
  </Items>
</WebTest>
XML

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (AppInsightsWebTestsResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_web_test" "test" {
  name                    = "acctestappinsightswebtests-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  kind                    = "ping"
  frequency               = 900
  timeout                 = 120
  enabled                 = true
  description             = "web_test"
  retry_enabled           = true
  tags = {
    ENV = "web_test"
  }
  geo_locations = ["us-tx-sn1-azr", "us-il-ch1-azr"]

  configuration = <<XML
<WebTest Name="WebTest1" Id="ABD48585-0831-40CB-9069-682EA6BB3583" Enabled="True" CssProjectStructure="" CssIteration="" Timeout="0" WorkItemIds="" xmlns="http://microsoft.com/schemas/VisualStudio/TeamTest/2010" Description="" CredentialUserName="" CredentialPassword="" PreAuthenticate="True" Proxy="default" StopOnError="False" RecordedResultFile="" ResultsLocale="">
  <Items>
    <Request Method="GET" Guid="a5f10126-e4cd-570d-961c-cea43999a200" Version="1.1" Url="http://microsoft.com" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="200" ExpectedResponseUrl="" ReportingName="" IgnoreHttpStatusCode="False" />
  </Items>
</WebTest>
XML

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (AppInsightsWebTestsResource) requiresImport(data acceptance.TestData) string {
	template := AppInsightsWebTestsResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_web_test" "import" {
  name                    = azurerm_application_insights_web_test.test.name
  location                = azurerm_application_insights_web_test.test.location
  resource_group_name     = azurerm_application_insights_web_test.test.resource_group_name
  application_insights_id = azurerm_application_insights_web_test.test.application_insights_id
  kind                    = azurerm_application_insights_web_test.test.kind
  configuration           = azurerm_application_insights_web_test.test.configuration
  geo_locations           = azurerm_application_insights_web_test.test.geo_locations
}
`, template)
}
