package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/advisor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAdvisorSuppression_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor_suppression", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisorSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAdvisorSuppression_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor_suppression", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisorSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorSuppressionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAdvisorSuppression_requiresImport),
		},
	})
}

func TestAccAzureRMAdvisorSuppression_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor_suppression", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisorSuppression_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAdvisorSuppression_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor_suppression", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorSuppressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisorSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAdvisorSuppression_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAdvisorSuppression_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAdvisorSuppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorSuppressionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAdvisorSuppressionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Advisor.SuppressionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("advisor Suppression not found: %s", resourceName)
		}

		id, err := parse.AdvisorSuppressionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceUri, id.RecommendationName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Advisor Suppression %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Advisor Suppression Client: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAdvisorSuppressionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Advisor.SuppressionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_advisor_suppression" {
			continue
		}

		id, err := parse.AdvisorSuppressionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceUri, id.RecommendationName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Advisor Suppression Client: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMAdvisorSuppression_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "test" {}

resource "azurerm_advisor_suppression" "test" {
  name              = "acctest-sp-%d"
  recommendation_id = data.azurerm_advisor_recommendations.test.recommendations.0.recommendation_id
}
`, data.RandomInteger)
}

func testAccAzureRMAdvisorSuppression_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAdvisorSuppression_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_advisor_suppression" "import" {
  name              = azurerm_advisor_suppression.test.name
  recommendation_id = azurerm_advisor_suppression.test.recommendation_id
}
`, template)
}

func testAccAzureRMAdvisorSuppression_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "test" {}

resource "azurerm_advisor_suppression" "test" {
  name              = "acctest-sp-%d"
  recommendation_id = data.azurerm_advisor_recommendations.test.recommendations.0.recommendation_id
  duration_in_days  = 1
}
`, data.RandomInteger)
}

func testAccAzureRMAdvisorSuppression_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "test" {}

resource "azurerm_advisor_suppression" "test" {
  name              = "acctest-sp-%d"
  recommendation_id = data.azurerm_advisor_recommendations.test.recommendations.0.recommendation_id
  duration_in_days  = 2
}
`, data.RandomInteger)
}
