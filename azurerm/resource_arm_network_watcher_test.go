package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkWatcher(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning one per region at once
	// (which our test suite can't easily workaround)
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":          testAccAzureRMNetworkWatcher_basic,
			"complete":       testAccAzureRMNetworkWatcher_complete,
			"update":         testAccAzureRMNetworkWatcher_update,
			"disappears":     testAccAzureRMNetworkWatcher_disappears,
			"importBasic":    testAccAzureRMNetworkWatcher_importBasic,
			"importComplete": testAccAzureRMNetworkWatcher_importComplete,
		},
		"PacketCapture": {
			"import":                     testAccAzureRMPacketCapture_importBasic,
			"localDisk":                  testAccAzureRMPacketCapture_localDisk,
			"storageAccount":             testAccAzureRMPacketCapture_storageAccount,
			"storageAccountAndLocalDisk": testAccAzureRMPacketCapture_storageAccountAndLocalDisk,
			"withFilters":                testAccAzureRMPacketCapture_withFilters,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMNetworkWatcher_basic(t *testing.T) {
	resourceGroup := "azurerm_network_watcher.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceGroup),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_complete(t *testing.T) {
	resourceGroup := "azurerm_network_watcher.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_completeConfig(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceGroup),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_update(t *testing.T) {
	resourceGroup := "azurerm_network_watcher.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceGroup),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcher_completeConfig(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceGroup),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_disappears(t *testing.T) {
	resourceGroup := "azurerm_network_watcher.test"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceGroup),
					testCheckAzureRMNetworkWatcherDisappears(resourceGroup),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNetworkWatcherExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Network Watcher: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).watcherClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Watcher %q (resource group: %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on watcherClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkWatcherDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Network Watcher: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).watcherClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Bad: Delete on watcherClient: %+v", err)
			}
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Bad: Delete on watcherClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkWatcherDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {

		if rs.Type != "azurerm_network_watcher" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).watcherClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Network Watcher still exists:\n%#v", resp)
			}
		}
	}

	return nil
}

func testAccAzureRMNetworkWatcher_basicConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestnw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMNetworkWatcher_completeConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestnw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tags {
	"Source" = "AccTests"
  }
}
`, rInt, location, rInt)
}
