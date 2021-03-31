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

func TestAccHPCCache_rootSquashDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.rootSquashDeprecated(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.rootSquashDeprecated(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		// Following import verification will cause diff given we simply set whatever is in cfg to state for "root_squash_enabled", since there is no
		// "cfg" during import verification, the state of the "root_squash_enabled" is always false.
		// The clarification is that since this is a deprecated property, users shouldn't import an existing resource to a new .tf file whilst using that
		// deprecated property.
		// data.ImportStep(),
		{
			Config: r.rootSquashDeprecated(data, false),
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

func TestAccHPCCache_defaultAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultAccessPolicyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultAccessPolicyComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultAccessPolicyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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
`, r.basic(data))
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

func (r HPCCacheResource) rootSquashDeprecated(data acceptance.TestData, enable bool) string {
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

func (r HPCCacheResource) defaultAccessPolicyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  default_access_policy {
    access_rule {
      scope  = "default"
      access = "rw"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) defaultAccessPolicyComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  default_access_policy {
    access_rule {
      scope  = "default"
      access = "ro"
    }

    access_rule {
      scope                   = "network"
      access                  = "rw"
      filter                  = "10.0.0.0/24"
      suid_enabled            = true
      submount_access_enabled = true
      root_squash_enabled     = true
      anonymous_uid           = 123
      anonymous_gid           = 123
    }

    access_rule {
      scope  = "host"
      access = "no"
      filter = "10.0.0.1"
    }
  }
}
`, r.template(data), data.RandomInteger)
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
