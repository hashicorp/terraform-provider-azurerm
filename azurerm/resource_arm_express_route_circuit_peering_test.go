package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeering(t *testing.T) {
	resourceName := "azurerm_express_route_circuit_peering.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeering(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "peering_type", "AzurePrivatePeering"),
					resource.TestCheckResourceAttr(resourceName, "microsoft_peering_config.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_key"}, //is not returned by the API
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_express_route_circuit_peering.test"
	ri := tf.AccRandTimeInt()

	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeering(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMExpressRouteCircuitPeering_requiresImportConfig(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_express_route_circuit_peering"),
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_microsoftPeering(t *testing.T) {
	resourceName := "azurerm_express_route_circuit_peering.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeering(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "peering_type", "MicrosoftPeering"),
					resource.TestCheckResourceAttr(resourceName, "microsoft_peering_config.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeeringWithCircuitUpdate(t *testing.T) {
	resourceName := "azurerm_express_route_circuit_peering.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeering(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "peering_type", "AzurePrivatePeering"),
					resource.TestCheckResourceAttr(resourceName, "microsoft_peering_config.#", "0"),
				),
			},
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeeringWithCircuitUpdate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMExpressRouteCircuitPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "peering_type", "AzurePrivatePeering"),
					resource.TestCheckResourceAttr(resourceName, "microsoft_peering_config.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_key"}, //is not returned by the API
			},
		},
	})
}

func testCheckAzureRMExpressRouteCircuitPeeringExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		peeringType := rs.Primary.Attributes["peering_type"]
		circuitName := rs.Primary.Attributes["express_route_circuit_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Express Route Circuit Peering: %s", peeringType)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRoutePeeringsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, circuitName, peeringType)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Express Route Circuit Peering %q (Circuit %q / Resource Group %q) does not exist", peeringType, circuitName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on expressRoutePeeringsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMExpressRouteCircuitPeeringDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_express_route_circuit_peering" {
			continue
		}

		peeringType := rs.Primary.Attributes["peering_type"]
		circuitName := rs.Primary.Attributes["express_route_circuit_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, circuitName, peeringType)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Express Route Circuit Peering still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMExpressRouteCircuitPeering_privatePeering(rInt int, location string) string {
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

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = "${azurerm_express_route_circuit.test.name}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  shared_key                    = "SSSSsssssshhhhhItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}
`, rInt, location, rInt)
}

func testAccAzureRMExpressRouteCircuitPeering_requiresImportConfig(rInt int, location string) string {
	template := testAccAzureRMExpressRouteCircuitPeering_privatePeering(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_peering" "import" {
  peering_type                  = "${azurerm_express_route_circuit_peering.test.peering_type}"
  express_route_circuit_name    = "${azurerm_express_route_circuit_peering.test.express_route_circuit_name}"
  resource_group_name           = "${azurerm_express_route_circuit_peering.test.resource_group_name}"
  shared_key                    = "${azurerm_express_route_circuit_peering.test.shared_key}"
  peer_asn                      = "${azurerm_express_route_circuit_peering.test.peer_asn}"
  primary_peer_address_prefix   = "${azurerm_express_route_circuit_peering.test.primary_peer_address_prefix}"
  secondary_peer_address_prefix = "${azurerm_express_route_circuit_peering.test.secondary_peer_address_prefix}"
  vlan_id                       = "${azurerm_express_route_circuit_peering.test.vlan_id}"
}
`, template)
}

func testAccAzureRMExpressRouteCircuitPeering_msPeering(rInt int, location string) string {
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
    tier   = "Premium"
    family = "MeteredData"
  }

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = "${azurerm_express_route_circuit.test.name}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.1.0.0/24"]
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMExpressRouteCircuitPeering_privatePeeringWithCircuitUpdate(rInt int, location string) string {
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

  tags = {
    Environment = "prod"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = "${azurerm_express_route_circuit.test.name}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  shared_key                    = "SSSSsssssshhhhhItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}
`, rInt, location, rInt)
}
