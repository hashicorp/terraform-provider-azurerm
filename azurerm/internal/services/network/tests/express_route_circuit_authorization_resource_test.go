package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMExpressRouteCircuitAuthorization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_authorization", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitAuthorizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitAuthorization_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitAuthorizationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "authorization_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMExpressRouteCircuitAuthorization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_authorization", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitAuthorizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitAuthorization_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitAuthorizationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "authorization_key"),
				),
			},
			{
				Config:      testAccAzureRMExpressRouteCircuitAuthorization_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_express_route_circuit_authorization"),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuitAuthorization_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_authorization", "test1")
	secondResourceName := "azurerm_express_route_circuit_authorization.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitAuthorizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitAuthorization_multipleConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitAuthorizationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "authorization_key"),
					testCheckAzureRMExpressRouteCircuitAuthorizationExists(secondResourceName),
					resource.TestCheckResourceAttrSet(secondResourceName, "authorization_key"),
				),
			},
		},
	})
}

func testCheckAzureRMExpressRouteCircuitAuthorizationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRouteAuthsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		authorizationName := rs.Primary.Attributes["name"]
		expressRouteCircuitName := rs.Primary.Attributes["express_route_circuit_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Express Route Circuit Authorization: %s", expressRouteCircuitName)
		}

		resp, err := client.Get(ctx, resourceGroup, expressRouteCircuitName, authorizationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Express Route Circuit Authorization %q (Circuit %q / Resource Group: %q) does not exist", authorizationName, expressRouteCircuitName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on expressRouteAuthsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMExpressRouteCircuitAuthorizationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRouteAuthsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_express_route_circuit" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		expressRouteCircuitName := rs.Primary.Attributes["express_route_circuit_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, expressRouteCircuitName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Express Route Circuit Authorization still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMExpressRouteCircuitAuthorization_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_authorization" "test" {
  name                       = "acctestauth%d"
  express_route_circuit_name = azurerm_express_route_circuit.test.name
  resource_group_name        = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMExpressRouteCircuitAuthorization_requiresImportConfig(data acceptance.TestData) string {
	template := testAccAzureRMExpressRouteCircuitAuthorization_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_authorization" "import" {
  name                       = azurerm_express_route_circuit_authorization.test.name
  express_route_circuit_name = azurerm_express_route_circuit_authorization.test.express_route_circuit_name
  resource_group_name        = azurerm_express_route_circuit_authorization.test.resource_group_name
}
`, template)
}

func testAccAzureRMExpressRouteCircuitAuthorization_multipleConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_authorization" "test1" {
  name                       = "acctestauth1%d"
  express_route_circuit_name = azurerm_express_route_circuit.test.name
  resource_group_name        = azurerm_resource_group.test.name
}

resource "azurerm_express_route_circuit_authorization" "test2" {
  name                       = "acctestauth2%d"
  express_route_circuit_name = azurerm_express_route_circuit.test.name
  resource_group_name        = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
