// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/ipampools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerIpamPoolResource struct{}

func testAccNetworkManagerIpamPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool", "test")
	r := ManagerIpamPoolResource{}

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

func testAccNetworkManagerIpamPool_basicIPv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool", "test")
	r := ManagerIpamPoolResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicIpv6(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerIpamPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool", "test")
	r := ManagerIpamPoolResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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

func testAccNetworkManagerIpamPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool", "test")
	r := ManagerIpamPoolResource{}

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

func testAccNetworkManagerIpamPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool", "test")
	r := ManagerIpamPoolResource{}

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

func (r ManagerIpamPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := ipampools.ParseIPamPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.IPamPools.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerIpamPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool" "test" {
  name               = "acctest-ipampool-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool1"
  address_prefixes   = ["10.0.0.0/27"]
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolResource) basicIpv6(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool" "test" {
  name               = "acctest-ipampool-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool1"
  address_prefixes   = ["2001:db8::/46"]
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_ipam_pool" "import" {
  name               = azurerm_network_manager_ipam_pool.test.name
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = azurerm_network_manager_ipam_pool.test.display_name
  address_prefixes   = azurerm_network_manager_ipam_pool.test.address_prefixes
}
`, r.basic(data))
}

func (r ManagerIpamPoolResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool" "test" {
  name               = "acctest-ipampool-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool2"
  address_prefixes   = ["10.0.0.0/27"]
  description        = "This is another test IPAM pool"

  tags = {
    foo = "bar"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool" "parent" {
  name               = "acctest-ipampool-p-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool0"
  address_prefixes   = ["10.0.0.0/26", "10.0.1.0/27"]
}

resource "azurerm_network_manager_ipam_pool" "test" {
  name               = "acctest-ipampool-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool1"
  parent_pool_name   = azurerm_network_manager_ipam_pool.parent.name
  address_prefixes   = ["10.0.0.0/27", "10.0.1.0/28"]
  description        = "This is a test IPAM pool"

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-ipampool-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-ipam-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity"]
}
`, data.RandomInteger, data.Locations.Primary)
}
