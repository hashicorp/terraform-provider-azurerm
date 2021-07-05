package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkWatcherResource struct {
}

func TestAccNetworkWatcher(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning one per region at once
	// (which our test suite can't easily workaround)

	// NOTE: Normally these tests can be separated to its own test cases, rather than this big composite one, since
	// we are not calling the `t.Parallel()` for each sub-test. However, currently nightly test are using the jen20/teamcity-go-test
	// which will invoke a `go test` for each test function, which effectively making them to be in parallel, even if they are intended
	// to be run in sequential.
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":          testAccNetworkWatcher_basic,
			"requiresImport": testAccNetworkWatcher_requiresImport,
			"complete":       testAccNetworkWatcher_complete,
			"update":         testAccNetworkWatcher_update,
			"disappears":     testAccNetworkWatcher_disappears,
		},
		"DataSource": {
			"basic": testAccDataSourceNetworkWatcher_basic,
		},
		"PacketCaptureOld": {
			"localDisk":                  testAccPacketCapture_localDisk,
			"storageAccount":             testAccPacketCapture_storageAccount,
			"storageAccountAndLocalDisk": testAccPacketCapture_storageAccountAndLocalDisk,
			"withFilters":                testAccPacketCapture_withFilters,
			"requiresImport":             testAccPacketCapture_requiresImport,
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
			"endpointType":                   testAccNetworkConnectionMonitor_endpointDeprecated,
			"updateEndpoint":                 testAccNetworkConnectionMonitor_updateEndpointIPAddressAndCoverageLevel,
		},
		"PacketCapture": {
			"localDisk":                  testAccNetworkPacketCapture_localDisk,
			"storageAccount":             testAccNetworkPacketCapture_storageAccount,
			"storageAccountAndLocalDisk": testAccNetworkPacketCapture_storageAccountAndLocalDisk,
			"withFilters":                testAccNetworkPacketCapture_withFilters,
			"requiresImport":             testAccNetworkPacketCapture_requiresImport,
		},
		"FlowLog": {
			"basic":                testAccNetworkWatcherFlowLog_basic,
			"disabled":             testAccNetworkWatcherFlowLog_disabled,
			"reenabled":            testAccNetworkWatcherFlowLog_reenabled,
			"retentionPolicy":      testAccNetworkWatcherFlowLog_retentionPolicy,
			"updateStorageAccount": testAccNetworkWatcherFlowLog_updateStorageAccount,
			"trafficAnalytics":     testAccNetworkWatcherFlowLog_trafficAnalytics,
			"version":              testAccNetworkWatcherFlowLog_version,
			"location":             testAccNetworkWatcherFlowLog_location,
			"tags":                 testAccNetworkWatcherFlowLog_tags,
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

func testAccNetworkWatcher_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")
	r := NetworkWatcherResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcher_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")
	r := NetworkWatcherResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_watcher"),
		},
	})
}

func testAccNetworkWatcher_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")
	r := NetworkWatcherResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcher_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")
	r := NetworkWatcherResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func testAccNetworkWatcher_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher", "test")
	r := NetworkWatcherResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicConfig,
			TestResource: r,
		}),
	})
}

func (t NetworkWatcherResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkWatcherID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.WatcherClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Network Watcher (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (NetworkWatcherResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NetworkWatcherID(state.ID)
	if err != nil {
		return nil, err
	}

	future, err := client.Network.WatcherClient.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("deleting Network Watcher %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Network.WatcherClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for Deletion on NetworkWatcherClient: %+v", err)
	}

	return utils.Bool(true), nil
}

func (NetworkWatcherResource) basicConfig(data acceptance.TestData) string {
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

func (r NetworkWatcherResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher" "import" {
  name                = azurerm_network_watcher.test.name
  location            = azurerm_network_watcher.test.location
  resource_group_name = azurerm_network_watcher.test.resource_group_name
}
`, r.basicConfig(data))
}

func (NetworkWatcherResource) completeConfig(data acceptance.TestData) string {
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
