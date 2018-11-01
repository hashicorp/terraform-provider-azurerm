package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLocalNetworkGateway_importBasic(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_basic(rInt, testLocation()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLocalNetworkGateway_importBGPSettingsComplete(t *testing.T) {
	resourceName := "azurerm_local_network_gateway.test"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLocalNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLocalNetworkGatewayConfig_bgpSettingsComplete(rInt, testLocation()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
