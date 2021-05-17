package storage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StorageObjectReplicationPolicyResource struct{}

func TestAccStorageObjectReplicationPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication_policy", "test")
	r := StorageObjectReplicationPolicyResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_policy_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageObjectReplicationPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication_policy", "test")
	r := StorageObjectReplicationPolicyResource{}
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

func TestAccStorageObjectReplicationPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication_policy", "test")
	r := StorageObjectReplicationPolicyResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_policy_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageObjectReplicationPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication_policy", "test")
	r := StorageObjectReplicationPolicyResource{}
	loc, _ := time.LoadLocation("Australia/Perth")
	copyTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z")
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_policy_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_policy_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, copyTime),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_policy_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateMoreRules(data, copyTime),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_policy_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_policy_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageObjectReplicationPolicyResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ObjectReplicationPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Storage.ObjectReplicationPolicyClient.Get(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r StorageObjectReplicationPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "src" {
  name     = "acctest-storage-src-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "src" {
  name                     = "stracctsrc%[3]s"
  resource_group_name      = azurerm_resource_group.src.name
  location                 = azurerm_resource_group.src.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  blob_properties {
    versioning_enabled  = true
    change_feed_enabled = true
  }
}

resource "azurerm_storage_container" "src" {
  name                  = "strcsrc%[3]s"
  storage_account_name  = azurerm_storage_account.src.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "src_second" {
  name                  = "strcsrcsecond%[3]s"
  storage_account_name  = azurerm_storage_account.src.name
  container_access_type = "private"
}

resource "azurerm_resource_group" "dst" {
  name     = "acctest-storage-dst-%[1]d"
  location = "%[4]s"
}

resource "azurerm_storage_account" "dst" {
  name                     = "stracctdst%[3]s"
  resource_group_name      = azurerm_resource_group.dst.name
  location                 = azurerm_resource_group.dst.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  blob_properties {
    versioning_enabled  = true
    change_feed_enabled = true
  }
}

resource "azurerm_storage_container" "dst" {
  name                  = "strcdst%[3]s"
  storage_account_name  = azurerm_storage_account.dst.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "dst_second" {
  name                  = "strcdstsecond%[3]s"
  storage_account_name  = azurerm_storage_account.dst.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Secondary)
}

func (r StorageObjectReplicationPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication_policy" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
  }
}
`, r.template(data))
}

func (r StorageObjectReplicationPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication_policy" "import" {
  source_storage_account_id      = azurerm_storage_object_replication_policy.test.source_storage_account_id
  destination_storage_account_id = azurerm_storage_object_replication_policy.test.destination_storage_account_id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
  }
}
`, r.basic(data))
}

func (r StorageObjectReplicationPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication_policy" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
    copy_over_from_time        = "Everything"
    filter_prefix_matches      = ["blobA", "blobB"]
  }
}

`, r.template(data))
}

func (r StorageObjectReplicationPolicyResource) update(data acceptance.TestData, copyTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication_policy" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
    copy_over_from_time        = "%s"
    filter_prefix_matches      = ["blobA", "blobB", "blobC"]
  }
}

`, r.template(data), copyTime)
}

func (r StorageObjectReplicationPolicyResource) updateMoreRules(data acceptance.TestData, copyTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication_policy" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
    copy_over_from_time        = "%s"
    filter_prefix_matches      = ["blobA", "blobB", "blobC"]
  }
  rules {
    source_container_name      = azurerm_storage_container.src_second.name
    destination_container_name = azurerm_storage_container.dst_second.name
    copy_over_from_time        = "Everything"
  }
}

`, r.template(data), copyTime)
}
