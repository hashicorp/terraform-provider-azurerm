// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageObjectReplicationResource struct{}

func TestAccStorageObjectReplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication", "test")
	r := StorageObjectReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageObjectReplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication", "test")
	r := StorageObjectReplicationResource{}
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

func TestAccStorageObjectReplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication", "test")
	r := StorageObjectReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageObjectReplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication", "test")
	r := StorageObjectReplicationResource{}
	loc, _ := time.LoadLocation("Australia/Perth")
	copyTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z")
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, copyTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateMoreRules(data, copyTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageObjectReplication_crossTenantDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication", "test")
	r := StorageObjectReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossTenantDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageObjectReplication_crossSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_object_replication", "test")
	if data.Subscriptions.Secondary == "" {
		t.Skipf("The secondary subscription is not specified")
	}
	r := StorageObjectReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossSubscription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_object_replication_id").Exists(),
				check.That(data.ResourceName).Key("destination_object_replication_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageObjectReplicationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ObjectReplicationID(state.ID)
	if err != nil {
		return nil, err
	}
	dstResp, err := client.Storage.ResourceManager.ObjectReplicationPolicies.Get(ctx, id.Dst)
	if err != nil {
		if response.WasNotFound(dstResp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	srcResp, err := client.Storage.ResourceManager.ObjectReplicationPolicies.Get(ctx, id.Src)
	if err != nil {
		if response.WasNotFound(srcResp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r StorageObjectReplicationResource) template(data acceptance.TestData) string {
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
  name     = "acctest-storage-alt-%[1]d"
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

func (r StorageObjectReplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
  }
}
`, r.template(data))
}

func (r StorageObjectReplicationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication" "import" {
  source_storage_account_id      = azurerm_storage_object_replication.test.source_storage_account_id
  destination_storage_account_id = azurerm_storage_object_replication.test.destination_storage_account_id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
  }
}
`, r.basic(data))
}

func (r StorageObjectReplicationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name        = azurerm_storage_container.src.name
    destination_container_name   = azurerm_storage_container.dst.name
    copy_blobs_created_after     = "Everything"
    filter_out_blobs_with_prefix = ["blobA", "blobB"]
  }
}
`, r.template(data))
}

func (r StorageObjectReplicationResource) update(data acceptance.TestData, copyTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name        = azurerm_storage_container.src.name
    destination_container_name   = azurerm_storage_container.dst.name
    copy_blobs_created_after     = "%s"
    filter_out_blobs_with_prefix = ["blobA", "blobB", "blobC"]
  }
}
`, r.template(data), copyTime)
}

func (r StorageObjectReplicationResource) updateMoreRules(data acceptance.TestData, copyTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_object_replication" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name        = azurerm_storage_container.src.name
    destination_container_name   = azurerm_storage_container.dst.name
    copy_blobs_created_after     = "%[2]s"
    filter_out_blobs_with_prefix = ["blobA", "blobB", "blobC"]
  }
  rules {
    source_container_name      = azurerm_storage_container.src_second.name
    destination_container_name = azurerm_storage_container.dst_second.name
    copy_blobs_created_after   = "Everything"
  }
}
`, r.template(data), copyTime)
}

func (r StorageObjectReplicationResource) crossSubscription(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azurerm-alt" {
  subscription_id = "%[1]s"
  features {}
}

resource "azurerm_resource_group" "src" {
  name     = "acctest-storage-src-%[2]d"
  location = "%[3]s"
}

resource "azurerm_storage_account" "src" {
  name                     = "stracctsrc%[4]s"
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
  name                  = "strcsrc%[4]s"
  storage_account_name  = azurerm_storage_account.src.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "src_second" {
  name                  = "strcsrcsecond%[4]s"
  storage_account_name  = azurerm_storage_account.src.name
  container_access_type = "private"
}

resource "azurerm_resource_group" "dst" {
  provider = azurerm-alt
  name     = "acctest-storage-dst-%[2]d"
  location = "%[5]s"
}

resource "azurerm_storage_account" "dst" {
  provider                 = azurerm-alt
  name                     = "stracctdst%[4]s"
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
  provider              = azurerm-alt
  name                  = "strcdst%[4]s"
  storage_account_name  = azurerm_storage_account.dst.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "dst_second" {
  provider              = azurerm-alt
  name                  = "strcdstsecond%[4]s"
  storage_account_name  = azurerm_storage_account.dst.name
  container_access_type = "private"
}

resource "azurerm_storage_object_replication" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
  }
}
`, data.Subscriptions.Secondary, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Secondary)
}

func (r StorageObjectReplicationResource) crossTenantDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "src" {
  name     = "acctest-storage-src-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "src" {
  name                             = "stracctsrc%[3]s"
  resource_group_name              = azurerm_resource_group.src.name
  location                         = azurerm_resource_group.src.location
  account_tier                     = "Standard"
  account_replication_type         = "LRS"
  cross_tenant_replication_enabled = false
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
  name     = "acctest-storage-alt-%[1]d"
  location = "%[4]s"
}

resource "azurerm_storage_account" "dst" {
  name                             = "stracctdst%[3]s"
  resource_group_name              = azurerm_resource_group.dst.name
  location                         = azurerm_resource_group.dst.location
  account_tier                     = "Standard"
  account_replication_type         = "LRS"
  cross_tenant_replication_enabled = false
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

resource "azurerm_storage_object_replication" "test" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Secondary)
}
