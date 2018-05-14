package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func testAccAzureRMNetworkWatcher_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "azurerm_network_watcher.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, testLocation()),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_importComplete(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "azurerm_network_watcher.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_completeConfig(rInt, testLocation()),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
