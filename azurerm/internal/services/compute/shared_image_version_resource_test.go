package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SharedImageVersionResource struct {
}

func TestAccSharedImageVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setup(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
			),
		},
		{
			Config: r.imageVersionUpdated(data),
			Check: resource.ComposeTestCheckFunc(
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
	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setup(data),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersionStorageAccountType(data, "Standard_LRS"),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setup(data),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersionStorageAccountType(data, "Standard_ZRS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_specializedImageVersionBySnapshot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imageVersionSpecializedBySnapshot(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imageVersionSpecializedByVM(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImageVersion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_version", "test")
	r := SharedImageVersionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setup(data),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.imageVersion(data),
			Check: resource.ComposeTestCheckFunc(
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

func (r SharedImageVersionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SharedImageVersionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.GalleryImageVersionsClient.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName, compute.ReplicationStatusTypesReplicationStatus)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Shared Image Gallery %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (SharedImageVersionResource) revokeSnapshot(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	snapShotName := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	future, err := client.Compute.SnapshotsClient.RevokeAccess(ctx, resourceGroup, snapShotName)
	if err != nil {
		return fmt.Errorf("bad: cannot revoke SAS on the snapshot: %+v", err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Compute.SnapshotsClient.Client); err != nil {
		return fmt.Errorf("bad: waiting the revoke of SAS on the snapshot: %+v", err)
	}

	return nil
}

// nolint: unparam
func (SharedImageVersionResource) setup(data acceptance.TestData) string {
	return ImageResource{}.setupUnmanagedDisks(data, "LRS")
}

func (SharedImageVersionResource) provision(data acceptance.TestData) string {
	template := ImageResource{}.standaloneImageProvision(data, "LRS", "")
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
  name                = "0.0.1"
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
