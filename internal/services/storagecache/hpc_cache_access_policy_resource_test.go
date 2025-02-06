// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/caches"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HPCCacheAccessPolicyResource struct{}

func TestAccHPCCacheAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheAccessPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func TestAccHPCCacheAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_access_policy", "test")
	r := HPCCacheAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r HPCCacheAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.StorageCache.Caches

	id, err := parse.CacheAccessPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, caches.NewCacheID(id.SubscriptionId, id.ResourceGroup, id.CacheName))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	m := resp.Model
	if m == nil {
		return pointer.To(false), nil
	}

	props := m.Properties
	if props == nil {
		return pointer.To(false), nil
	}

	settings := props.SecuritySettings
	if settings == nil {
		return pointer.To(false), nil
	}

	policies := settings.AccessPolicies
	if policies == nil {
		return pointer.To(false), nil
	}

	return pointer.To(storagecache.CacheGetAccessPolicyByName(*policies, id.Name) != nil), nil
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
  address_prefixes     = ["10.0.2.0/24"]
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
