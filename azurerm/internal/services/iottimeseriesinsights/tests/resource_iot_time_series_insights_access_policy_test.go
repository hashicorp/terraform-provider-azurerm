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
)

func TestAccAzureRMIoTTimeSeriesInsightsAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIoTTimeSeriesInsightsAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsAccessPolicy_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIoTTimeSeriesInsightsAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTTimeSeriesInsights.AccessPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parse.TimeSeriesInsightsAccessPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on TimeSeriesInsightsAccessPolicyClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMIoTTimeSeriesInsightsAccessPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTTimeSeriesInsights.AccessPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iot_time_series_insights_access_policy" {
			continue
		}

		id, err := parse.TimeSeriesInsightsAccessPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("time Series Insights Access Policy still exists: %q", id.Name)
		}
	}

	return nil
}

func testAccAzureRMIoTTimeSeriesInsightsAccessPolicy_basic(data acceptance.TestData) string {
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
resource "azurerm_iot_time_series_insights_access_policy" "test" {
  name                                = "accTEst_tsiap%d"
  time_series_insights_environment_id = azurerm_iot_time_series_insights_standard_environment.test.id

  principal_object_id = "aGUID"
  roles               = ["Reader"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMIoTTimeSeriesInsightsAccessPolicy_update(data acceptance.TestData) string {
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
resource "azurerm_iot_time_series_insights_access_policy" "test" {
  name                                = "accTEst_tsiap%d"
  time_series_insights_environment_id = azurerm_iot_time_series_insights_standard_environment.test.id

  principal_object_id = "aGUID"
  roles               = ["Contributor"]
  description         = "Test Access Policy"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
