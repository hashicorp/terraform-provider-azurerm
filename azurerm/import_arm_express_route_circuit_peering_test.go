package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMExpressRouteCircuitPeering_importAzurePrivatePeering(t *testing.T) {
	rInt := acctest.RandInt()
	location := testLocation()
	resourceName := "azurerm_express_route_circuit_peering.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_azurePrivatePeering(rInt, location),
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

func TestAccAzureRMExpressRouteCircuitPeering_importAzurePublicPeering(t *testing.T) {
	rInt := acctest.RandInt()
	location := testLocation()
	resourceName := "azurerm_express_route_circuit_peering.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_azurePublicPeering(rInt, location),
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

func TestAccAzureRMExpressRouteCircuitPeering_importMicrosoftPeering(t *testing.T) {
	rInt := acctest.RandInt()
	location := testLocation()
	resourceName := "azurerm_express_route_circuit_peering.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMExpressRouteCircuitPeering_microsoftPeering(rInt, location),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
