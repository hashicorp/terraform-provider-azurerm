package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMExpressRouteCircuitAuthorization_basic(t *testing.T) {
	resourceName := "azurerm_express_route_circuit_authorization.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitAuthorizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitAuthorization_basicConfig(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitAuthorizationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_key"),
				),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuitAuthorization_multiple(t *testing.T) {
	firstResourceName := "azurerm_express_route_circuit_authorization.test1"
	secondResourceName := "azurerm_express_route_circuit_authorization.test2"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitAuthorizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitAuthorization_multipleConfig(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitAuthorizationExists(firstResourceName),
					resource.TestCheckResourceAttrSet(firstResourceName, "authorization_key"),
					testCheckAzureRMExpressRouteCircuitAuthorizationExists(secondResourceName),
					resource.TestCheckResourceAttrSet(secondResourceName, "authorization_key"),
				),
			},
		},
	})
}

func testCheckAzureRMExpressRouteCircuitAuthorizationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		authorizationName := rs.Primary.Attributes["name"]
		expressRouteCircuitName := rs.Primary.Attributes["express_route_circuit_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Express Route Circuit Authorization: %s", expressRouteCircuitName)
		}

		client := testAccProvider.Meta().(*ArmClient).expressRouteAuthsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, expressRouteCircuitName, authorizationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Express Route Circuit Authorization %q (Circuit %q / Resource Group: %q) does not exist", name, expressRouteCircuitName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on expressRouteAuthsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMExpressRouteCircuitAuthorizationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).expressRouteAuthsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMExpressRouteCircuitAuthorization_basicConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_authorization" "test" {
  name                       = "acctestauth%d"
  express_route_circuit_name = "${azurerm_express_route_circuit.test.name}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMExpressRouteCircuitAuthorization_multipleConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_authorization" "test1" {
  name                       = "acctestauth1%d"
  express_route_circuit_name = "${azurerm_express_route_circuit.test.name}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
}

resource "azurerm_express_route_circuit_authorization" "test2" {
  name                       = "acctestauth2%d"
  express_route_circuit_name = "${azurerm_express_route_circuit.test.name}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}
