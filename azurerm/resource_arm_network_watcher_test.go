package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkWatcher(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning one per region at once
	// (which our test suite can't easily workaround)
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":          testAccAzureRMNetworkWatcher_basic,
			"requiresImport": testAccAzureRMNetworkWatcher_requiresImport,
			"complete":       testAccAzureRMNetworkWatcher_complete,
			"update":         testAccAzureRMNetworkWatcher_update,
			"disappears":     testAccAzureRMNetworkWatcher_disappears,
		},
		"DataSource": {
			"basic": testAccDataSourceAzureRMNetworkWatcher_basic,
		},
		"ConnectionMonitorOld": {
			"addressBasic":              testAccAzureRMConnectionMonitor_addressBasic,
			"addressComplete":           testAccAzureRMConnectionMonitor_addressComplete,
			"addressUpdate":             testAccAzureRMConnectionMonitor_addressUpdate,
			"vmBasic":                   testAccAzureRMConnectionMonitor_vmBasic,
			"vmComplete":                testAccAzureRMConnectionMonitor_vmComplete,
			"vmUpdate":                  testAccAzureRMConnectionMonitor_vmUpdate,
			"destinationUpdate":         testAccAzureRMConnectionMonitor_destinationUpdate,
			"missingDestinationInvalid": testAccAzureRMConnectionMonitor_missingDestination,
			"bothDestinationsInvalid":   testAccAzureRMConnectionMonitor_conflictingDestinations,
			"requiresImport":            testAccAzureRMConnectionMonitor_requiresImport,
		},
		"PacketCaptureOld": {
			"localDisk":                  testAccAzureRMPacketCapture_localDisk,
			"storageAccount":             testAccAzureRMPacketCapture_storageAccount,
			"storageAccountAndLocalDisk": testAccAzureRMPacketCapture_storageAccountAndLocalDisk,
			"withFilters":                testAccAzureRMPacketCapture_withFilters,
			"requiresImport":             testAccAzureRMPacketCapture_requiresImport,
		},
		"ConnectionMonitor": {
			"addressBasic":              testAccAzureRMNetworkConnectionMonitor_addressBasic,
			"addressComplete":           testAccAzureRMNetworkConnectionMonitor_addressComplete,
			"addressUpdate":             testAccAzureRMNetworkConnectionMonitor_addressUpdate,
			"vmBasic":                   testAccAzureRMNetworkConnectionMonitor_vmBasic,
			"vmComplete":                testAccAzureRMNetworkConnectionMonitor_vmComplete,
			"vmUpdate":                  testAccAzureRMNetworkConnectionMonitor_vmUpdate,
			"destinationUpdate":         testAccAzureRMNetworkConnectionMonitor_destinationUpdate,
			"missingDestinationInvalid": testAccAzureRMNetworkConnectionMonitor_missingDestination,
			"bothDestinationsInvalid":   testAccAzureRMNetworkConnectionMonitor_conflictingDestinations,
			"requiresImport":            testAccAzureRMNetworkConnectionMonitor_requiresImport,
		},
		"PacketCapture": {
			"localDisk":                  testAccAzureRMNetworkPacketCapture_localDisk,
			"storageAccount":             testAccAzureRMNetworkPacketCapture_storageAccount,
			"storageAccountAndLocalDisk": testAccAzureRMNetworkPacketCapture_storageAccountAndLocalDisk,
			"withFilters":                testAccAzureRMNetworkPacketCapture_withFilters,
			"requiresImport":             testAccAzureRMNetworkPacketCapture_requiresImport,
		},
		"FlowLog": {
			"basic":                testAccAzureRMNetworkWatcherFlowLog_basic,
			"disabled":             testAccAzureRMNetworkWatcherFlowLog_disabled,
			"reenabled":            testAccAzureRMNetworkWatcherFlowLog_reenabled,
			"retentionPolicy":      testAccAzureRMNetworkWatcherFlowLog_retentionPolicy,
			"updateStorageAccount": testAccAzureRMNetworkWatcherFlowLog_updateStorageAccount,
			"trafficAnalytics":     testAccAzureRMNetworkWatcherFlowLog_trafficAnalytics,
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
	resourceName := "azurerm_network_watcher.test"
	rInt := tf.AccRandTimeInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceName),
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

func testAccAzureRMNetworkWatcher_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_network_watcher.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkWatcher_requiresImportConfig(rInt, location),
				ExpectError: acceptance.RequiresImportError("azurerm_network_watcher"),
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_complete(t *testing.T) {
	resourceName := "azurerm_network_watcher.test"
	rInt := tf.AccRandTimeInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_completeConfig(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceName),
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

func testAccAzureRMNetworkWatcher_update(t *testing.T) {
	resourceName := "azurerm_network_watcher.test"
	rInt := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcher_completeConfig(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceName),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_disappears(t *testing.T) {
	resourceName := "azurerm_network_watcher.test"
	rInt := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(resourceName),
					testCheckAzureRMNetworkWatcherDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNetworkWatcherExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Network Watcher: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WatcherClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMNetworkWatcherDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Network Watcher: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WatcherClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Bad: Delete on watcherClient: %+v", err)
			}
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WatcherClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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
  name     = "acctestRG-watcher-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestNW-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMNetworkWatcher_requiresImportConfig(rInt int, location string) string {
	template := testAccAzureRMNetworkWatcher_basicConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher" "import" {
  name                = "${azurerm_network_watcher.test.name}"
  location            = "${azurerm_network_watcher.test.location}"
  resource_group_name = "${azurerm_network_watcher.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMNetworkWatcher_completeConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-watcher-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestNW-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    "Source" = "AccTests"
  }
}
`, rInt, location, rInt)
}
