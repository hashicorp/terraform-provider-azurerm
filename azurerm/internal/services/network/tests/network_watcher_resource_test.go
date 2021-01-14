package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
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
		"PacketCaptureOld": {
			"localDisk":                  testAccAzureRMPacketCapture_localDisk,
			"storageAccount":             testAccAzureRMPacketCapture_storageAccount,
			"storageAccountAndLocalDisk": testAccAzureRMPacketCapture_storageAccountAndLocalDisk,
			"withFilters":                testAccAzureRMPacketCapture_withFilters,
			"requiresImport":             testAccAzureRMPacketCapture_requiresImport,
		},
		"ConnectionMonitor": {
			"addressBasic":                   testAccNetworkConnectionMonitor_addressBasic,
			"addressComplete":                testAccNetworkConnectionMonitor_addressComplete,
			"addressUpdate":                  testAccNetworkConnectionMonitor_addressUpdate,
			"vmBasic":                        testAccNetworkConnectionMonitor_vmBasic,
			"vmComplete":                     testAccNetworkConnectionMonitor_vmComplete,
			"vmUpdate":                       testAccNetworkConnectionMonitor_vmUpdate,
			"destinationUpdate":              testAccNetworkConnectionMonitor_destinationUpdate,
			"missingDestinationInvalid":      testAccNetworkConnectionMonitor_missingDestination,
			"bothDestinationsInvalid":        testAccNetworkConnectionMonitor_conflictingDestinations,
			"requiresImport":                 testAccNetworkConnectionMonitor_requiresImport,
			"httpConfiguration":              testAccNetworkConnectionMonitor_httpConfiguration,
			"icmpConfiguration":              testAccNetworkConnectionMonitor_icmpConfiguration,
			"bothAddressAndVirtualMachineId": testAccNetworkConnectionMonitor_withAddressAndVirtualMachineId,
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
			"version":              testAccAzureRMNetworkWatcherFlowLog_version,
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
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkWatcher_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkWatcher_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_watcher"),
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkWatcher_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcher_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcher_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcher_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherExists(data.ResourceName),
					testCheckAzureRMNetworkWatcherDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNetworkWatcherExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WatcherClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.NetworkWatcherID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if id.ResourceGroup == "" {
			return fmt.Errorf("Bad: no resource group found in state for Network Watcher: %q", id.Name)
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Watcher %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on watcherClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkWatcherDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WatcherClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		id, err := parse.NetworkWatcherID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if id.ResourceGroup == "" {
			return fmt.Errorf("Bad: no resource group found in state for Network Watcher: %q", id.Name)
		}

		future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WatcherClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_watcher" {
			continue
		}

		id, err := parse.NetworkWatcherID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Network Watcher still exists:\n%#v", resp)
			}
		}
	}

	return nil
}

func testAccAzureRMNetworkWatcher_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-watcher-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestNW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMNetworkWatcher_requiresImportConfig(data acceptance.TestData) string {
	template := testAccAzureRMNetworkWatcher_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher" "import" {
  name                = azurerm_network_watcher.test.name
  location            = azurerm_network_watcher.test.location
  resource_group_name = azurerm_network_watcher.test.resource_group_name
}
`, template)
}

func testAccAzureRMNetworkWatcher_completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-watcher-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestNW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "Source" = "AccTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
