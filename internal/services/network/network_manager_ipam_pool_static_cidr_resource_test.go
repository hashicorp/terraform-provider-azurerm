// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/staticcidrs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerIpamPoolStaticCidrResource struct{}

func testAccNetworkManagerIpamPoolStaticCidr_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool_static_cidr", "test")
	r := ManagerIpamPoolStaticCidrResource{}

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

func testAccNetworkManagerIpamPoolStaticCidr_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool_static_cidr", "test")
	r := ManagerIpamPoolStaticCidrResource{}

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

func testAccNetworkManagerIpamPoolStaticCidr_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool_static_cidr", "test")
	r := ManagerIpamPoolStaticCidrResource{}

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

func testAccNetworkManagerIpamPoolStaticCidr_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool_static_cidr", "test")
	r := ManagerIpamPoolStaticCidrResource{}

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

func testAccNetworkManagerIpamPoolStaticCidr_ipAddressNumber(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool_static_cidr", "test")
	r := ManagerIpamPoolStaticCidrResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipAddressNumber(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerIpamPoolStaticCidr_ipAddressNumberUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_ipam_pool_static_cidr", "test")
	r := ManagerIpamPoolStaticCidrResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipAddressNumber(data),
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

func (r ManagerIpamPoolStaticCidrResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := staticcidrs.ParseStaticCidrID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.StaticCidrs.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerIpamPoolStaticCidrResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool_static_cidr" "test" {
  name             = "acctest-ipsc-%[2]d"
  ipam_pool_id     = azurerm_network_manager_ipam_pool.test.id
  address_prefixes = ["10.0.0.0/27"]
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolStaticCidrResource) ipAddressNumber(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool_static_cidr" "test" {
  name                               = "acctest-ipsc-%[2]d"
  ipam_pool_id                       = azurerm_network_manager_ipam_pool.test.id
  number_of_ip_addresses_to_allocate = "16"
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolStaticCidrResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_ipam_pool_static_cidr" "import" {
  name             = azurerm_network_manager_ipam_pool_static_cidr.test.name
  ipam_pool_id     = azurerm_network_manager_ipam_pool_static_cidr.test.ipam_pool_id
  address_prefixes = azurerm_network_manager_ipam_pool_static_cidr.test.address_prefixes
}
`, r.basic(data))
}

func (r ManagerIpamPoolStaticCidrResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool_static_cidr" "test" {
  name             = "acctest-ipsc-%[2]d"
  ipam_pool_id     = azurerm_network_manager_ipam_pool.test.id
  address_prefixes = ["10.0.0.0/26", "10.0.0.128/27"]
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolStaticCidrResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-ipampool-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-ipsc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
}

resource "azurerm_network_manager_ipam_pool" "test" {
  name               = "acctest-ipampool-%[1]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool1"
  address_prefixes   = ["10.0.0.0/24"]
}
`, data.RandomInteger, data.Locations.Primary)
}
