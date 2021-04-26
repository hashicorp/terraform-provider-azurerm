package storage_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
)

type StorageBlobInventoryPolicyResource struct{}

func TestAccStorageBlobInventoryPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob_inventory_policy", "test")
	r := StorageBlobInventoryPolicyResource{}
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

func TestAccStorageBlobInventoryPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob_inventory_policy", "test")
	r := StorageBlobInventoryPolicyResource{}
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

func TestAccStorageBlobInventoryPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob_inventory_policy", "test")
	r := StorageBlobInventoryPolicyResource{}
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

func TestAccStorageBlobInventoryPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob_inventory_policy", "test")
	r := StorageBlobInventoryPolicyResource{}
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

func (r StorageBlobInventoryPolicyResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.BlobInventoryPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Storage.BlobInventoryPoliciesClient.Get(ctx, id.ResourceGroup, id.StorageAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	if props := resp.BlobInventoryPolicyProperties; props != nil {
		if policy := props.Policy; policy != nil {
			if policy.Enabled == nil || !*policy.Enabled {
				return utils.Bool(false), nil
			}
		}
	}
	return utils.Bool(true), nil
}

func (r StorageBlobInventoryPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_blob_public_access = true
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageBlobInventoryPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob_inventory_policy" "test" {
  storage_account_id     = azurerm_storage_account.test.id
  storage_container_name = azurerm_storage_container.test.name
  rules {
    name = "rule1"
    filter {
      blob_types   = ["blockBlob"]
    }
  }
}
`, r.template(data))
}

func (r StorageBlobInventoryPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob_inventory_policy" "import" {
  storage_account_id     = azurerm_storage_blob_inventory_policy.test.storage_account_id
  storage_container_name = azurerm_storage_blob_inventory_policy.test.storage_container_name
  rules {
    name = azurerm_storage_blob_inventory_policy.test.rules.0.name
    filter {
      blob_types = azurerm_storage_blob_inventory_policy.test.rules.0.filter.0.blob_types

    }
  }
}
`, config)
}

func (r StorageBlobInventoryPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob_inventory_policy" "test" {
  name                = "acctest-sbip-%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_storage_account.test.name
  policy {
    destination = "containerName"
    enabled     = true
    rules {
      name = "inventoryPolicyRule1"
      definition {
        filters {
          blob_types            = ["blockBlob", "appendBlob", "pageBlob"]
          include_blob_versions = true
          include_snapshots     = true
          prefix_match          = ["inventoryprefix1", "inventoryprefix2"]
        }
      }
      enabled = true
    }
    type = "Inventory"
  }
}
`, template, data.RandomInteger)
}

func (r StorageBlobInventoryPolicyResource) updatePolicy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_blob_inventory_policy" "test" {
  name                = "acctest-sbip-%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_storage_account.test.name
  policy {
    destination = "containerName"
    enabled     = true
    rules {
      name = "inventoryPolicyRule1"
      definition {
        filters {
          blob_types            = ["blockBlob", "appendBlob", "pageBlob"]
          include_blob_versions = true
          include_snapshots     = true
          prefix_match          = ["inventoryprefix1", "inventoryprefix2"]
        }
      }
      enabled = true
    }
    type = "Inventory"
  }
}
`, template, data.RandomInteger)
}
