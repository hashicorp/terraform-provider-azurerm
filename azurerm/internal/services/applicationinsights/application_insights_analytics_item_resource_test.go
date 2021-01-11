package applicationinsights_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppInsightsAnalyticsItemResource struct {
}

func TestAccApplicationInsightsAnalyticsItem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_analytics_item", "test")
	r := AppInsightsAnalyticsItemResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("testquery"),
				check.That(data.ResourceName).Key("scope").HasValue("shared"),
				check.That(data.ResourceName).Key("type").HasValue("query"),
				check.That(data.ResourceName).Key("content").HasValue("requests #test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsAnalyticsItem_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_analytics_item", "test")
	r := AppInsightsAnalyticsItemResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("testquery"),
				check.That(data.ResourceName).Key("scope").HasValue("shared"),
				check.That(data.ResourceName).Key("type").HasValue("query"),
				check.That(data.ResourceName).Key("content").HasValue("requests #test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("testquery"),
				check.That(data.ResourceName).Key("scope").HasValue("shared"),
				check.That(data.ResourceName).Key("type").HasValue("query"),
				check.That(data.ResourceName).Key("content").HasValue("requests #updated"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsAnalyticsItem_multiple(t *testing.T) {
	r1 := acceptance.BuildTestData(t, "azurerm_application_insights_analytics_item", "test1")
	r2 := acceptance.BuildTestData(t, "azurerm_application_insights_analytics_item", "test2")
	r3 := acceptance.BuildTestData(t, "azurerm_application_insights_analytics_item", "test3")
	r := AppInsightsAnalyticsItemResource{}

	r1.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiple(r1),
			Check: resource.ComposeTestCheckFunc(
				check.That(r1.ResourceName).ExistsInAzure(r),
				check.That(r2.ResourceName).ExistsInAzure(r),
				check.That(r3.ResourceName).ExistsInAzure(r),
				check.That(r1.ResourceName).Key("name").HasValue("testquery1"),
				check.That(r1.ResourceName).Key("scope").HasValue("shared"),
				check.That(r1.ResourceName).Key("type").HasValue("query"),
				check.That(r1.ResourceName).Key("content").HasValue("requests #test1"),
				check.That(r2.ResourceName).Key("name").HasValue("testquery2"),
				check.That(r2.ResourceName).Key("scope").HasValue("user"),
				check.That(r2.ResourceName).Key("type").HasValue("query"),
				check.That(r2.ResourceName).Key("content").HasValue("requests #test2"),
				check.That(r3.ResourceName).Key("name").HasValue("testfunction1"),
				check.That(r3.ResourceName).Key("scope").HasValue("shared"),
				check.That(r3.ResourceName).Key("type").HasValue("function"),
				check.That(r3.ResourceName).Key("content").HasValue("requests #test3"),
				check.That(r3.ResourceName).Key("function_alias").HasValue("myfunction"),
			),
		},
		r1.ImportStep(),
		r2.ImportStep(),
		r3.ImportStep(),
	})
}

func (t AppInsightsAnalyticsItemResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resGroup, appInsightsName, itemScopePath, itemID, err := applicationinsights.ResourcesArmApplicationInsightsAnalyticsItemParseID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse ID (id: %s): %+v", state.ID, err)
	}

	resp, err := clients.AppInsights.AnalyticsItemsClient.Get(ctx, resGroup, appInsightsName, itemScopePath, itemID, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Application Insights AnalyticsItem %q (resource group: %q, app insight: %s, item scope: %s): %+v", resGroup, appInsightsName, itemScopePath, itemID, err)
	}

	return utils.Bool(resp.StatusCode != http.StatusNotFound), nil
}

func (AppInsightsAnalyticsItemResource) basic(data acceptance.TestData) string {
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

resource "azurerm_application_insights_analytics_item" "test" {
  name                    = "testquery"
  application_insights_id = azurerm_application_insights.test.id
  content                 = "requests #test"
  scope                   = "shared"
  type                    = "query"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsAnalyticsItemResource) basic2(data acceptance.TestData) string {
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

resource "azurerm_application_insights_analytics_item" "test" {
  name                    = "testquery"
  application_insights_id = azurerm_application_insights.test.id
  content                 = "requests #updated"
  scope                   = "shared"
  type                    = "query"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsAnalyticsItemResource) multiple(data acceptance.TestData) string {
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

resource "azurerm_application_insights_analytics_item" "test1" {
  name                    = "testquery1"
  application_insights_id = azurerm_application_insights.test.id
  content                 = "requests #test1"
  scope                   = "shared"
  type                    = "query"
}

resource "azurerm_application_insights_analytics_item" "test2" {
  name                    = "testquery2"
  application_insights_id = azurerm_application_insights.test.id
  content                 = "requests #test2"
  scope                   = "user"
  type                    = "query"
}

resource "azurerm_application_insights_analytics_item" "test3" {
  name                    = "testfunction1"
  application_insights_id = azurerm_application_insights.test.id
  content                 = "requests #test3"
  scope                   = "shared"
  type                    = "function"
  function_alias          = "myfunction"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
