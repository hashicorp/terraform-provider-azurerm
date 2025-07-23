// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SharedImageVersionResource struct{}

func TestAccSharedImageVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
			),
		},
		{
			Config: r.imageVersionUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("2"),
				check.That(data.ResourceName).Key("name").HasValue("1234567890.1234567890.1234567890"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_storageAccountTypeLrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setup(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersionStorageAccountType(data, "Standard_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_storageAccountTypeZrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setup(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersionStorageAccountType(data, "Standard_ZRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_blobURI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageVersionBlobURI(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_specializedImageVersionBySnapshot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageVersionSpecializedBySnapshot(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// the share image version will generate a shared access signature (SAS) on the referenced snapshot and keep it active until the replication is complete
				// in the meantime, the service will return success of creation before the replication complete.
				// therefore in this test, we have to revoke the access on the snapshot in order to do the cleaning work
				data.CheckWithClientForResource(r.revokeSnapshot, "azurerm_snapshot.test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_specializedImageVersionByVM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageVersionSpecializedByVM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_diskEncryptionSetID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.diskEncryptionSetID(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_endOfLifeDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	endOfLifeDate := time.Now().Add(time.Hour * 10).Format(time.RFC3339)
	endOfLifeDateUpdated := time.Now().Add(time.Hour * 20).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.endOfLifeDate(data, endOfLifeDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.endOfLifeDate(data, endOfLifeDateUpdated),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_replicationMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.replicationMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_replicatedRegionDeletion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.replicatedRegionDeletion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setup(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_shared_image_version"),
		},
	})
}

func (r SharedImageVersionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := galleryimageversions.ParseImageVersionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.GalleryImageVersionsClient.Get(ctx, *id, galleryimageversions.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (SharedImageVersionResource) revokeSnapshot(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	subscriptionId := client.Account.SubscriptionId
	snapShotName := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()
	}

	snapshotId := snapshots.NewSnapshotID(subscriptionId, resourceGroup, snapShotName)
	if err := client.Compute.SnapshotsClient.RevokeAccessThenPoll(ctx, snapshotId); err != nil {
		return fmt.Errorf("revoking SAS on %s: %+v", snapshotId, err)
	}

	return nil
}

// nolint: unparam
func (SharedImageVersionResource) setup(data acceptance.TestData) string {
	return ImageResource{}.setupUnmanagedDisks(data)
}

func (SharedImageVersionResource) provision(data acceptance.TestData) string {
	template := ImageResource{}.standaloneImageProvision(data, "")
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SharedImageVersionResource) imageVersion(data acceptance.TestData) string {
	template := r.provision(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "1234567890.1234567890.1234567890"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    "foo" = "bar"
  }
}
`, template)
}

func (r SharedImageVersionResource) imageVersionBlobURI(data acceptance.TestData) string {
	template := r.setup(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  blob_uri           = azurerm_virtual_machine.testsource.storage_os_disk[0].vhd_uri
  storage_account_id = azurerm_storage_account.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}
`, template, data.RandomInteger)
}

func (r SharedImageVersionResource) provisionSpecialized(data acceptance.TestData) string {
	template := ImageResource{}.setupManagedDisks(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  specialized         = true

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SharedImageVersionResource) imageVersionSpecializedBySnapshot(data acceptance.TestData) string {
	template := r.provisionSpecialized(data)
	return fmt.Sprintf(`
%s

resource "azurerm_snapshot" "test" {
  name                = "acctestsnapshot%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_resource_id  = azurerm_virtual_machine.testsource.storage_os_disk.0.managed_disk_id
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  os_disk_snapshot_id = azurerm_snapshot.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    "foo" = "bar"
  }
}
`, template, data.RandomInteger)
}

func (r SharedImageVersionResource) imageVersionSpecializedByVM(data acceptance.TestData) string {
	template := r.provisionSpecialized(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  managed_image_id = azurerm_virtual_machine.testsource.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    "foo" = "bar"
  }
}
`, template)
}

func (r SharedImageVersionResource) imageVersionStorageAccountType(data acceptance.TestData, storageAccountType string) string {
	template := r.provision(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
    storage_account_type   = "%s"
  }
}
`, template, storageAccountType)
}

func (r SharedImageVersionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "import" {
  name                = azurerm_shared_image_version.test.name
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name
  location            = azurerm_shared_image_version.test.location
  managed_image_id    = azurerm_shared_image_version.test.managed_image_id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}
`, r.imageVersion(data))
}

func (r SharedImageVersionResource) imageVersionUpdated(data acceptance.TestData) string {
	template := r.provision(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "1234567890.1234567890.1234567890"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  target_region {
    name                   = "%s"
    regional_replica_count = 2
  }
}
`, template, data.Locations.Secondary)
}

func (r SharedImageVersionResource) diskEncryptionSetID(data acceptance.TestData) string {
	template := r.provision(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
  }

  depends_on = [
    azurerm_key_vault_access_policy.disk-encryption,
  ]
}
`, template, data.RandomString, data.RandomInteger)
}

func (r SharedImageVersionResource) endOfLifeDate(data acceptance.TestData, endOfLifeDate string) string {
	template := r.provision(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id
  end_of_life_date    = "%s"

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}
`, template, endOfLifeDate)
}

func (r SharedImageVersionResource) replicationMode(data acceptance.TestData) string {
	template := r.provision(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id
  replication_mode    = "Shallow"

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}
`, template)
}

func (r SharedImageVersionResource) replicatedRegionDeletion(data acceptance.TestData) string {
	template := r.provision(data)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                                     = "0.0.1"
  gallery_name                             = azurerm_shared_image_gallery.test.name
  image_name                               = azurerm_shared_image.test.name
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = azurerm_resource_group.test.location
  managed_image_id                         = azurerm_image.test.id
  deletion_of_replicated_locations_enabled = true

  target_region {
    name                        = azurerm_resource_group.test.location
    regional_replica_count      = 1
    exclude_from_latest_enabled = true
  }
}
`, template)
}
