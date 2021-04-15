package hpccache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HPCCacheAccessPolicyResource struct{}

func TestAccHPCCacheAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheAccessPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

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

func (r HPCCacheAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.HPCCache.CachesClient

	id, err := parse.CacheAccessPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.CacheName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	props := resp.CacheProperties
	if props == nil {
		return utils.Bool(false), nil
	}

	settings := props.SecuritySettings
	if settings == nil {
		return utils.Bool(false), nil
	}

	policies := settings.AccessPolicies
	if policies == nil {
		return utils.Bool(false), nil
	}

	return utils.Bool(hpccache.CacheGetAccessPolicyByName(*policies, id.Name) != nil), nil
}

func (r HPCCacheAccessPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_access_policy" "test" {
  name         = "testAccessPolicy"
  hpc_cache_id = azurerm_hpc_cache.test.id
  access_rule {
    scope  = "default"
    access = "rw"
  }
}
`, template)
}

func (r HPCCacheAccessPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_access_policy" "test" {
  name         = "testAccessPolicy"
  hpc_cache_id = azurerm_hpc_cache.test.id

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
`, template)
}

func (r HPCCacheAccessPolicyResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_access_policy" "import" {
  name         = azurerm_hpc_cache_access_policy.test.name
  hpc_cache_id = azurerm_hpc_cache_access_policy.test.hpc_cache_id
  access_rule {
    scope  = "default"
    access = "rw"
  }
}
`, template)
}

func (r HPCCacheAccessPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hpcc-%d"
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

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
