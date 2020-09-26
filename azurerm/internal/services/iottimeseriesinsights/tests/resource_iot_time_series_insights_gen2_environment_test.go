package tests

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

func TestAccAzureRMIoTTimeSeriesInsightsGen2Environment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_gen2_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsGen2Environment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIoTTimeSeriesInsightsGen2Environment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_gen2_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsGen2Environment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsGen2Environment_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsGen2Environment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIoTTimeSeriesInsightsGen2Environment_multiple_property_ids(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_gen2_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsGen2Environment_multiple_property_ids(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
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
			return fmt.Errorf("Bad: Get on TimeSeriesInsightsGen2EnvironmentClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Time Series Insights Gen2 Environment %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMIoTTimeSeriesInsightsGen2EnvironmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iot_time_series_insights_gen2_environment" {
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
			return fmt.Errorf("time Series Insights Gen2 Environment still exists: %q", id.Name)
		}
	}

	return nil
}

func testAccAzureRMIoTTimeSeriesInsightsGen2Environment_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}
resource "azurerm_storage_account" "storage" {
  name                     = "acctestsatsi%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "test" {
  name                = "acctest_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "L1"
  data_retention_time = "P30D"
  property {
    ids = ["id"]
  }
  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMIoTTimeSeriesInsightsGen2Environment_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}
resource "azurerm_storage_account" "storage" {
  name                     = "acctestsatsi%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "test" {
  name                = "acctest_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "L1"
  data_retention_time = "P30D"
  property {
    ids = ["newid"]
  }
  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMIoTTimeSeriesInsightsGen2Environment_multiple_property_ids(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}
resource "azurerm_storage_account" "storage" {
  name                     = "acctestsatsi%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "test" {
  name                = "acctest_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "L1"
  data_retention_time = "P30D"
  property {
    ids = ["id", "secondId"]
  }
  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
