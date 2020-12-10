package hpccache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HPCCacheResource struct {
}

func TestAccHPCCache_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_mtu(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.mtu(data, 1000),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.mtu(data, 1500),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.mtu(data, 1000),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_rootSquash(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.rootSquash(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.rootSquash(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.rootSquash(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (HPCCacheResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CacheID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HPCCache.CachesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving HPC Cache (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.CacheProperties != nil), nil
}

func (r HPCCacheResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "import" {
  name                = azurerm_hpc_cache.test.name
  resource_group_name = azurerm_hpc_cache.test.resource_group_name
  location            = azurerm_hpc_cache.test.location
  cache_size_in_gb    = azurerm_hpc_cache.test.cache_size_in_gb
  subnet_id           = azurerm_hpc_cache.test.subnet_id
  sku_name            = azurerm_hpc_cache.test.sku_name
}
`, r.template(data))
}

func (r HPCCacheResource) mtu(data acceptance.TestData, mtu int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  mtu                 = %d
}
`, r.template(data), data.RandomInteger, mtu)
}

func (r HPCCacheResource) rootSquash(data acceptance.TestData, enable bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  root_squash_enabled = %t
}
`, r.template(data), data.RandomInteger, enable)
}

func (HPCCacheResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
