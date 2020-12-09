package iottimeseriesinsights_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccIoTTimeSeriesInsightsReferenceDataSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_reference_data_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckIoTTimeSeriesInsightsReferenceDataSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIoTTimeSeriesInsightsReferenceDataSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckIoTTimeSeriesInsightsReferenceDataSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccIoTTimeSeriesInsightsReferenceDataSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_reference_data_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckIoTTimeSeriesInsightsReferenceDataSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIoTTimeSeriesInsightsReferenceDataSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckIoTTimeSeriesInsightsReferenceDataSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccIoTTimeSeriesInsightsReferenceDataSet_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckIoTTimeSeriesInsightsReferenceDataSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccIoTTimeSeriesInsightsReferenceDataSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckIoTTimeSeriesInsightsReferenceDataSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckIoTTimeSeriesInsightsReferenceDataSetExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parse.ReferenceDataSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on TimeSeriesInsightsReferenceDataSetClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Time Series Insights Reference Data Set %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckIoTTimeSeriesInsightsReferenceDataSetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iot_time_series_insights_reference_data_set" {
			continue
		}

		id, err := parse.ReferenceDataSetID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("time Series Insights Reference Data Set still exists: %q", id.Name)
		}
	}

	return nil
}

func testAccIoTTimeSeriesInsightsReferenceDataSet_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}
resource "azurerm_iot_time_series_insights_standard_environment" "test" {
  name                = "accTEst_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S1_1"
  data_retention_time = "P30D"
}

resource "azurerm_iot_time_series_insights_reference_data_set" "test" {
  name                                = "accTEsttsd%d"
  time_series_insights_environment_id = azurerm_iot_time_series_insights_standard_environment.test.id
  location                            = azurerm_resource_group.test.location

  key_property {
    name = "keyProperty1"
    type = "String"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccIoTTimeSeriesInsightsReferenceDataSet_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}
resource "azurerm_iot_time_series_insights_standard_environment" "test" {
  name                = "accTEst_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S1_1"
  data_retention_time = "P30D"
}
resource "azurerm_iot_time_series_insights_reference_data_set" "test" {
  name                                = "accTEsttsd%d"
  time_series_insights_environment_id = azurerm_iot_time_series_insights_standard_environment.test.id
  location                            = azurerm_resource_group.test.location

  key_property {
    name = "keyProperty1"
    type = "String"
  }

  key_property {
    name = "keyProperty2"
    type = "Bool"
  }

  tags = {
    Environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
