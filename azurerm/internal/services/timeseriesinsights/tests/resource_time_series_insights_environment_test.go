package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/timeseriesinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMTimeSeriesInsightsEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_time_series_insights_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTimeSeriesInsightsEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTimeSeriesInsightsEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTimeSeriesInsightsEnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTimeSeriesInsightsEnvironment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_time_series_insights_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTimeSeriesInsightsEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTimeSeriesInsightsEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTimeSeriesInsightsEnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTimeSeriesInsightsEnvironment_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTimeSeriesInsightsEnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTimeSeriesInsightsEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTimeSeriesInsightsEnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMTimeSeriesInsightsEnvironmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).TimeSeriesInsights.EnvironmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parse.TimeSeriesInsightsEnvironmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on TimeSeriesInsightsEnvironmentClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Time Series Insights Environment %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMTimeSeriesInsightsEnvironmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).TimeSeriesInsights.EnvironmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_time_series_insights_environment" {
			continue
		}

		id, err := parse.TimeSeriesInsightsEnvironmentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Time Series Insights Environment still exists: %q", id.Name)
		}
	}

	return nil
}

func testAccAzureRMTimeSeriesInsightsEnvironment_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}
resource "azurerm_time_series_insights_environment" "test" {
  name                = "accTEst_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S1_1"
  data_retention_time = "P30D"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMTimeSeriesInsightsEnvironment_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}
resource "azurerm_time_series_insights_environment" "test" {
  name                = "accTEst_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S2_5"
  data_retention_time = "P30D"

  storage_limited_exceeded_behavior = "PauseIngress"
  partition_key                     = "foo"

  tags = {
    Environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
