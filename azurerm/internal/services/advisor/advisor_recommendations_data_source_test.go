package advisor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AdvisorRecommendationsDataSourceTests struct{}

func TestAccAdvisorRecommendationsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_advisor_recommendations", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AdvisorRecommendationsDataSourceTests{}.basicConfig(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("recommendations.#").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.category").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.description").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.impact").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.recommendation_name").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.recommendation_type_id").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.resource_name").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.resource_type").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.updated_time").Exists(),
			),
		},
	})
}

func TestAccAdvisorRecommendationsDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_advisor_recommendations", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AdvisorRecommendationsDataSourceTests{}.completeConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("recommendations.#").Exists(),
			),
		},
	})
}

func TestAccAdvisorRecommendationsDataSource_categoriesFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_advisor_recommendations", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AdvisorRecommendationsDataSourceTests{}.categoriesFilterConfig(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("recommendations.#").Exists(),
				check.That(data.ResourceName).Key("recommendations.0.category").HasValue("Cost"),
			),
		},
	})
}

func (AdvisorRecommendationsDataSourceTests) basicConfig() string {
	return `provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "test" {}`
}

// Advisor generated recommendations needs long time to take effects, sometimes up to one day or more,
// Please refer to the issue https://github.com/Azure/azure-rest-api-specs/issues/9284
// So here we get an empty list of recommendations
func (AdvisorRecommendationsDataSourceTests) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-advisor-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                      = "accteststr%s"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  enable_https_traffic_only = false
  account_tier              = "Standard"
  account_replication_type  = "LRS"
}

data "azurerm_advisor_recommendations" "test" {
  filter_by_category        = ["security"]
  filter_by_resource_groups = [azurerm_resource_group.test.name]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (AdvisorRecommendationsDataSourceTests) categoriesFilterConfig() string {
	return `provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "test" {
  filter_by_category = ["cost"]
}
`
}
