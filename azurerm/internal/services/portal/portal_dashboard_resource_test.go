package portal_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccPortalDashboard_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceArmDashboard_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDashboardExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDashboardExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Portal.DashboardsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dashboardName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Dashboard: %s", dashboardName)
		}

		resp, err := client.Get(ctx, resourceGroup, dashboardName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dashboardsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Dashboard %q (resource group: %q) does not exist", dashboardName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDashboardDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Portal.DashboardsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dashboard" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Dashboard still exists:\n%+v", resp)
		}
	}

	return nil
}

func testResourceArmDashboard_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dashboard" "test" {
  name                 = "my-test-dashboard"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  dashboard_properties = <<DASH
{
   "lenses": {
        "0": {
            "order": 0,
            "parts": {
                "0": {
                    "position": {
                        "x": 0,
                        "y": 0,
                        "rowSpan": 2,
                        "colSpan": 3
                    },
                    "metadata": {
                        "inputs": [],
                        "type": "Extension/HubsExtension/PartType/MarkdownPart",
                        "settings": {
                            "content": {
                                "settings": {
                                    "content": "## This is only a test :)",
                                    "subtitle": "",
                                    "title": "Test MD Tile"
                                }
                            }
                        }
                    }
				}
			}
		}
	}
}
DASH
}
`, data.RandomInteger, data.Locations.Primary)
}
