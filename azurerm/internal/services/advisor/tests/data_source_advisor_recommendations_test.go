package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAdvisorRecommendations_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_advisor_recommendations", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckArmAdvisorRecommendations_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.category"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.impact"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.recommendation_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.recommendation_type_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.resource_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.resource_type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.0.updated_time"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAdvisorRecommendations_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_advisor_recommendations", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckArmAdvisorRecommendations_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.#"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAdvisorRecommendations_categoriesFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_advisor_recommendations", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckArmAdvisorRecommendations_categoriesFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "recommendations.#"),
					resource.TestCheckResourceAttr(data.ResourceName, "recommendations.0.category", "Cost"),
				),
			},
		},
	})
}

const testAccCheckArmAdvisorRecommendations_basic = `
provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "test" { }
`

// Advisor genereate recommendations needs long time to take effects, sometimes up to one day or more,
// Please refer to the issue https://github.com/Azure/azure-rest-api-specs/issues/9284
// So here we get an empty list of recommendations
func testAccCheckArmAdvisorRecommendations_complete(data acceptance.TestData) string {
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

const testAccCheckArmAdvisorRecommendations_categoriesFilter = `
provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "test" {
  filter_by_category           = ["cost"]
}
`
