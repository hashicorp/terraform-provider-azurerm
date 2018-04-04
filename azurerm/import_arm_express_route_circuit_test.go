package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func testAccAzureRMExpressRouteCircuit_importMetered(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"

	ri := acctest.RandInt()
	config := testAccAzureRMExpressRouteCircuit_basicMeteredConfig(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMExpressRouteCircuit_importUnlimited(t *testing.T) {
	resourceName := "azurerm_express_route_circuit.test"

	ri := acctest.RandInt()
	config := testAccAzureRMExpressRouteCircuit_basicUnlimitedConfig(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
