// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/storagetargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HPCCacheBlobTargetResource struct{}

func TestAccHPCCacheBlobTarget_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")
	r := HPCCacheBlobTargetResource{}

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

func TestAccHPCCacheBlobTarget_accessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")
	r := HPCCacheBlobTargetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.accessPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.accessPolicyUpdate(data),
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

func TestAccHPCCacheBlobTarget_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")
	r := HPCCacheBlobTargetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.namespace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCacheBlobTarget_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_blob_target", "test")
	r := HPCCacheBlobTargetResource{}

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

func (HPCCacheBlobTargetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storagetargets.ParseStorageTargetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.StorageCache.StorageTargets.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving HPC Cache Blob Target (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HPCCacheBlobTargetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  namespace_path       = "/blob_storage1"
  access_policy_name   = "default"
}
`, r.cacheTemplate(data), data.RandomString)
}

func (r HPCCacheBlobTargetResource) namespace(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  namespace_path       = "/blob_storage2"
}
`, r.cacheTemplate(data), data.RandomString)
}

func (r HPCCacheBlobTargetResource) accessPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_access_policy" "test" {
  name         = "p1"
  hpc_cache_id = azurerm_hpc_cache.test.id
  access_rule {
    scope  = "default"
    access = "rw"
  }

  # This is not needed in Terraform v0.13, whilst needed in v0.14.
  # Once https://github.com/hashicorp/terraform/issues/28193 is fixed, we can remove this lifecycle block.
  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  namespace_path       = "/blob_storage1"
  access_policy_name   = azurerm_hpc_cache_access_policy.test.name
}
`, r.cacheTemplate(data), data.RandomString)
}

func (r HPCCacheBlobTargetResource) accessPolicyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_access_policy" "test" {
  name         = "p2"
  hpc_cache_id = azurerm_hpc_cache.test.id
  access_rule {
    scope  = "default"
    access = "rw"
  }
  # This is necessary to make the Terraform apply order works correctly.
  # Without CBD: azurerm_hpc_cache_access_policy-p1 (destroy) -> azurerm_hpc_cache_blob_target (update) -> azurerm_hpc_cache_access_policy-p2 (create)
  # 			 (the 1st step wil fail as the access policy is under used by the blob target)
  # With CBD   : azurerm_hpc_cache_access_policy-p2 (create) -> azurerm_hpc_cache_blob_target (update) -> azurerm_hpc_cache_access_policy-p1 (delete)
  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_hpc_cache_blob_target" "test" {
  name                 = "acctest-HPCCTGT-%s"
  resource_group_name  = azurerm_resource_group.test.name
  cache_name           = azurerm_hpc_cache.test.name
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  namespace_path       = "/blob_storage1"
  access_policy_name   = azurerm_hpc_cache_access_policy.test.name
}
`, r.cacheTemplate(data), data.RandomString)
}

func (r HPCCacheBlobTargetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_blob_target" "import" {
  name                 = azurerm_hpc_cache_blob_target.test.name
  resource_group_name  = azurerm_hpc_cache_blob_target.test.resource_group_name
  cache_name           = azurerm_hpc_cache_blob_target.test.cache_name
  storage_container_id = azurerm_hpc_cache_blob_target.test.storage_container_id
  namespace_path       = azurerm_hpc_cache_blob_target.test.namespace_path
}
`, r.basic(data))
}

func (r HPCCacheBlobTargetResource) cacheTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"

  # hpc_cache_blob_target depends on below role_assignments, however these role_assignments need up to 5 minutes to take effect.
  # Since hpc_cache_blob_target depends on the hpc_cache and hpc_cache takes far more than 5 minutes to create, put the dependency here so role_assignments are ready before creating hpc_cache_blob_target.
  depends_on = [
    azurerm_role_assignment.test_storage_account_contrib,
    azurerm_role_assignment.test_storage_blob_data_contrib,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (HPCCacheBlobTargetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

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
  address_prefixes     = ["10.0.2.0/24"]
}

data "azuread_service_principal" "test" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststorgacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctest-strgctn-%s"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_role_assignment" "test_storage_account_contrib" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}
