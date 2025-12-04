// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerResource struct{}

func TestAccNetworkManager(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning one (connectivity or securityAdmin) network manager per subscription at once
	// (which our test suite can't easily work around)

	testCases := map[string]map[string]func(t *testing.T){
		"Manager": {
			"basic":          testAccNetworkManager_basic,
			"complete":       testAccNetworkManager_complete,
			"update":         testAccNetworkManager_update,
			"requiresImport": testAccNetworkManager_requiresImport,
			"dataSource":     testAccNetworkManagerDataSource_complete,
		},
		"NetworkGroup": {
			"basic":          testAccNetworkManagerNetworkGroup_basic,
			"complete":       testAccNetworkManagerNetworkGroup_complete,
			"update":         testAccNetworkManagerNetworkGroup_update,
			"requiresImport": testAccNetworkManagerNetworkGroup_requiresImport,
			"dataSource":     testAccNetworkManagerNetworkGroupDataSource_complete,
		},
		"SubscriptionConnection": {
			"basic":          testAccNetworkSubscriptionNetworkManagerConnection_basic,
			"complete":       testAccNetworkSubscriptionNetworkManagerConnection_complete,
			"update":         testAccNetworkSubscriptionNetworkManagerConnection_update,
			"requiresImport": testAccNetworkSubscriptionNetworkManagerConnection_requiresImport,
		},
		"ManagementGroupConnection": {
			"basic":          testAccNetworkManagerManagementGroupConnection_basic,
			"complete":       testAccNetworkManagerManagementGroupConnection_complete,
			"update":         testAccNetworkManagerManagementGroupConnection_update,
			"requiresImport": testAccNetworkManagerManagementGroupConnection_requiresImport,
		},
		"ScopeConnection": {
			"basic":          testAccNetworkManagerScopeConnection_basic,
			"complete":       testAccNetworkManagerScopeConnection_complete,
			"update":         testAccNetworkManagerScopeConnection_update,
			"requiresImport": testAccNetworkManagerScopeConnection_requiresImport,
		},
		"StaticMember": {
			"basic":          testAccNetworkManagerStaticMember_basic,
			"requiresImport": testAccNetworkManagerStaticMember_requiresImport,
		},
		"ConnectivityConfiguration": {
			"basic":             testAccNetworkManagerConnectivityConfiguration_basic,
			"basicTopologyMesh": testAccNetworkManagerConnectivityConfiguration_basicTopologyMesh,
			"complete":          testAccNetworkManagerConnectivityConfiguration_complete,
			"update":            testAccNetworkManagerConnectivityConfiguration_update,
			"requiresImport":    testAccNetworkManagerConnectivityConfiguration_requiresImport,
			"dataSource":        testAccNetworkManagerConnectivityConfigurationDataSource_basic,
		},
		"SecurityAdminConfiguration": {
			"basic":          testAccNetworkManagerSecurityAdminConfiguration_basic,
			"complete":       testAccNetworkManagerSecurityAdminConfiguration_complete,
			"update":         testAccNetworkManagerSecurityAdminConfiguration_update,
			"requiresImport": testAccNetworkManagerSecurityAdminConfiguration_requiresImport,
		},
		"AdminRuleCollection": {
			"basic":          testAccNetworkManagerAdminRuleCollection_basic,
			"complete":       testAccNetworkManagerAdminRuleCollection_complete,
			"update":         testAccNetworkManagerAdminRuleCollection_update,
			"requiresImport": testAccNetworkManagerAdminRuleCollection_requiresImport,
		},
		"AdminRule": {
			"basic":          testAccNetworkManagerAdminRule_basic,
			"complete":       testAccNetworkManagerAdminRule_complete,
			"update":         testAccNetworkManagerAdminRule_update,
			"requiresImport": testAccNetworkManagerAdminRule_requiresImport,
		},
		"Deployment": {
			"basic":          testAccNetworkManagerDeployment_basic,
			"basicAdmin":     testAccNetworkManagerDeployment_basicAdmin,
			"complete":       testAccNetworkManagerDeployment_complete,
			"update":         testAccNetworkManagerDeployment_update,
			"withTriggers":   testAccNetworkManagerDeployment_withTriggers,
			"requiresImport": testAccNetworkManagerDeployment_requiresImport,
		},
		"IPAMPool": {
			"basic":          testAccNetworkManagerIpamPool_basic,
			"basicIPv6":      testAccNetworkManagerIpamPool_basicIPv6,
			"complete":       testAccNetworkManagerIpamPool_complete,
			"update":         testAccNetworkManagerIpamPool_update,
			"requiresImport": testAccNetworkManagerIpamPool_requiresImport,
			"dataSource":     testAccNetworkManagerIpamPoolDataSource_complete,
		},
		"IPAMPoolStaticCIDR": {
			"basic":                 testAccNetworkManagerIpamPoolStaticCidr_basic,
			"complete":              testAccNetworkManagerIpamPoolStaticCidr_complete,
			"update":                testAccNetworkManagerIpamPoolStaticCidr_update,
			"requiresImport":        testAccNetworkManagerIpamPoolStaticCidr_requiresImport,
			"ipAddressNumber":       testAccNetworkManagerIpamPoolStaticCidr_ipAddressNumber,
			"updateIpAddressNumber": testAccNetworkManagerIpamPoolStaticCidr_ipAddressNumberUpdated,
		},
		"VerifierWorkspace": {
			"basic":          testAccNetworkManagerVerifierWorkspace_basic,
			"complete":       testAccNetworkManagerVerifierWorkspace_complete,
			"update":         testAccNetworkManagerVerifierWorkspace_update,
			"requiresImport": testAccNetworkManagerVerifierWorkspace_requiresImport,
		},
		"VerifierWorkspaceReachabilityAnalysisIntent": {
			"basic":          testAccNetworkManagerVerifierWorkspaceReachabilityAnalysisIntent_basic,
			"complete":       testAccNetworkManagerVerifierWorkspaceReachabilityAnalysisIntent_complete,
			"requiresImport": testAccNetworkManagerVerifierWorkspaceReachabilityAnalysisIntent_requiresImport,
		},
		"RoutingConfiguration": {
			"basic":          testAccNetworkManagerRoutingConfiguration_basic,
			"complete":       testAccNetworkManagerRoutingConfiguration_complete,
			"update":         testAccNetworkManagerRoutingConfiguration_update,
			"requiresImport": testAccNetworkManagerRoutingConfiguration_requiresImport,
		},
		"RoutingRuleCollection": {
			"basic":          testAccNetworkManagerRoutingRuleCollection_basic,
			"complete":       testAccNetworkManagerRoutingRuleCollection_complete,
			"update":         testAccNetworkManagerRoutingRuleCollection_update,
			"requiresImport": testAccNetworkManagerRoutingRuleCollection_requiresImport,
		},
		"SubnetIPAMPool": {
			"ipAddressPool":              testAccSubnet_ipAddressPool,
			"ipAddressPoolVNet":          testAccSubnet_ipAddressPoolVNet,
			"ipAddressPoolIPv6":          testAccSubnet_ipAddressPoolIPv6,
			"ipAddressPoolBlockUpdated":  testAccSubnet_ipAddressPoolBlockUpdated,
			"ipAddressPoolNumberUpdated": testAccSubnet_ipAddressPoolNumberUpdated,
		},
		"VNETIPANPool": {
			"ipAddressPool":             testAccVirtualNetwork_ipAddressPool,
			"ipAddressPoolIPv6":         testAccVirtualNetwork_ipAddressPoolIPv6,
			"ipAddressPoolMultiple":     testAccVirtualNetwork_ipAddressPoolMultiple,
			"ipAddressPoolUpdateBasic":  testAccVirtualNetwork_ipAddressPoolUpdateBasic,
			"ipAddressPoolUpdateNumber": testAccVirtualNetwork_ipAddressPoolUpdateNumber,
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

func testAccNetworkManager_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager", "test")
	r := ManagerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManager_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager", "test")
	r := ManagerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManager_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager", "test")
	r := ManagerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManager_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager", "test")
	r := ManagerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ManagerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkmanagers.ParseNetworkManagerID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Network.NetworkManagers.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_network_manager" "test" {
  name                = "acctest-networkmanager-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_network_manager" "import" {
  name                = azurerm_network_manager.test.name
  location            = azurerm_network_manager.test.location
  resource_group_name = azurerm_network_manager.test.resource_group_name
  scope {
    subscription_ids = azurerm_network_manager.test.scope.0.subscription_ids
  }
}
`, r.basic(data))
}

func (r ManagerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_network_manager" "test" {
  name                = "acctest-networkmanager-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin", "Routing"]
  description    = "test network manager"
  tags = {
    foo = "bar"
  }
}
`, r.template(data), data.RandomInteger)
}

func (ManagerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-%d"
  location = "%s"
}
data "azurerm_subscription" "current" {
}
`, data.RandomInteger, data.Locations.Primary)
}
