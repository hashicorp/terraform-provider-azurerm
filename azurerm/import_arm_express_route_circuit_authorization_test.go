package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func testAccAzureRMExpressRouteCircuitAuthorization_importBasic(t *testing.T) {
	resourceName := "azurerm_express_route_circuit_authorization.test"

	ri := acctest.RandInt()
	config := testAccAzureRMExpressRouteCircuitAuthorization_basicConfig(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitAuthorizationDestroy,
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
