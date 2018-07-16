package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func testAccAzureRMExpressRouteCircuitPeering_importAzurePrivatePeering(t *testing.T) {
	rInt := acctest.RandInt()
	location := testLocation()
	resourceName := "azurerm_express_route_circuit_peering.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_privatePeering(rInt, location),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_key"},
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuitPeering_importMicrosoftPeering(t *testing.T) {
	rInt := acctest.RandInt()
	location := testLocation()
	resourceName := "azurerm_express_route_circuit_peering.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_msPeering(rInt, location),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
