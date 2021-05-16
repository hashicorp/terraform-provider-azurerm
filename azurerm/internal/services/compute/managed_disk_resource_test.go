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

type ManagedDiskResource struct {
}

func TestAccManagedDisk_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_managed_disk"),
		},
	})
}

func TestAccManagedDisk_zeroGbFromPlatformImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.zeroGbFromPlatformImage(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectNonEmptyPlan: true, // since the `disk_size_gb` will have changed
		},
	})
}

func TestAccManagedDisk_import(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}
	vm := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then delete it so we can use the vhd to test import
			Config:             vm.basicLinuxMachine(data),
			Destroy:            false,
			ExpectNonEmptyPlan: true,
			Check: resource.ComposeTestCheckFunc(
				// TODO: switch to using `azurerm_linux_virtual_machine` once Binary Testing is enabled
				check.That("azurerm_virtual_machine.test").ExistsInAzure(vm),
				data.CheckWithClientForResource(r.destroyVirtualMachine, "azurerm_virtual_machine.test"),
			),
		},
		{
			Config: r.importConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_copy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.copy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_fromPlatformImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.platformImage(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
				check.That(data.ResourceName).Key("tags.cost-center").HasValue("ops"),
				check.That(data.ResourceName).Key("disk_size_gb").HasValue("1"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue(string(compute.StorageAccountTypesStandardLRS)),
			),
		},
		{
			Config: r.empty_updated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
				check.That(data.ResourceName).Key("disk_size_gb").HasValue("2"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue(string(compute.StorageAccountTypesPremiumLRS)),
			),
		},
	})
}

func TestAccManagedDisk_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.encryption(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.0.secret_url").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.0.source_vault_id").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.0.key_url").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.0.source_vault_id").Exists(),
			),
		},
	})
}

func TestAccManagedDisk_importEmpty_withZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty_withZone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_create_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.create_withUltraSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_update_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.create_withUltraSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_iops_read_write").HasValue("101"),
				check.That(data.ResourceName).Key("disk_mbps_read_write").HasValue("10"),
			),
		},
		{
			Config: r.update_withUltraSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_iops_read_write").HasValue("102"),
				check.That(data.ResourceName).Key("disk_mbps_read_write").HasValue("11"),
			),
		},
	})
}

func TestAccManagedDisk_import_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.create_withUltraSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.import_withUltraSSD),
	})
}

func TestAccManagedDisk_diskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.diskEncryptionSetEncrypted(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_diskEncryptionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.diskEncryptionSetUnencrypted(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.diskEncryptionSetEncrypted(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_attachedDiskUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.managedDiskAttached(data, 10),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.managedDiskAttached(data, 20),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_size_gb").HasValue("20"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_attachedStorageTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.storageTypeUpdateWhilstAttached(data, "Standard_LRS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.storageTypeUpdateWhilstAttached(data, "Premium_LRS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_attachedTierUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tierUpdateWhileAttached(data, "P10"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("P10"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tierUpdateWhileAttached(data, "P20"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("P20"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMManagedDisk_networkPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: testAccAzureRMManagedDisk_networkPolicy_create(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMManagedDisk_networkPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: testAccAzureRMManagedDisk_networkPolicy_create(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "DenyAll"),
			),
		},
		data.ImportStep(),
		{
			Config: testAccAzureRMManagedDisk_networkPolicy_update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "DenyAll"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMManagedDisk_networkPolicy_create_withAllowPrivate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: testAccAzureRMManagedDisk_networkPolicy_create_withAllowPrivate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMManagedDisk_networkPolicy_update_withAllowPrivate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: testAccAzureRMManagedDisk_networkPolicy_create_withAllowPrivate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "AllowPrivate"),
			),
		},
		data.ImportStep(),
		{
			Config: testAccAzureRMManagedDisk_networkPolicy_update_withAllowPrivate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "AllowPrivate"),
			),
		},
		data.ImportStep(),
	})
}

func (ManagedDiskResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ManagedDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.DisksClient.Get(ctx, id.ResourceGroup, id.DiskName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Managed Disk %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ManagedDiskResource) destroyVirtualMachine(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	vmName := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	// this is a preview feature we don't want to use right now
	var forceDelete *bool = nil
	future, err := client.Compute.VMClient.Delete(ctx, resourceGroup, vmName, forceDelete)
	if err != nil {
		return fmt.Errorf("Bad: Delete on vmClient: %+v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Compute.VMClient.Client); err != nil {
		return fmt.Errorf("Bad: Delete on vmClient: %+v", err)
	}

	return nil
}

func (ManagedDiskResource) empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) requiresImport(data acceptance.TestData) string {
	template := ManagedDiskResource{}.empty(data)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "import" {
  name                 = azurerm_managed_disk.test.name
  location             = azurerm_managed_disk.test.location
  resource_group_name  = azurerm_managed_disk.test.resource_group_name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, template)
}

func (ManagedDiskResource) empty_withZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
  zones                = ["1"]

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) importConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Import"
  source_uri           = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
  storage_account_id   = azurerm_storage_account.test.id
  disk_size_gb         = "45"

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ManagedDiskResource) copy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "source" {
  name                 = "acctestd1-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd2-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Copy"
  source_resource_id   = azurerm_managed_disk.source.id
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ManagedDiskResource) empty_updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "2"

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) platformImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  os_type              = "Linux"
  create_option        = "FromImage"
  image_reference_id   = data.azurerm_platform_image.test.id
  storage_account_type = "Standard_LRS"
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) zeroGbFromPlatformImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  os_type              = "Linux"
  create_option        = "FromImage"
  disk_size_gb         = 0
  image_reference_id   = data.azurerm_platform_image.test.id
  storage_account_type = "Standard_LRS"
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
  sku_name            = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

    key_permissions = [
      "create",
      "delete",
      "get",
      "purge",
    ]

    secret_permissions = [
      "delete",
      "get",
      "set",
    ]
  }

  enabled_for_disk_encryption = true

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  encryption_settings {
    enabled = true

    disk_encryption_key {
      secret_url      = "${azurerm_key_vault_secret.test.id}"
      source_vault_id = "${azurerm_key_vault.test.id}"
    }

    key_encryption_key {
      key_url         = "${azurerm_key_vault_key.test.id}"
      source_vault_id = "${azurerm_key_vault.test.id}"
    }
  }

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger)
}

func (ManagedDiskResource) create_withUltraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "101"
  disk_mbps_read_write = "10"
  zones                = ["1"]

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) update_withUltraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "102"
  disk_mbps_read_write = "11"
  zones                = ["1"]

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ManagedDiskResource) import_withUltraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "import" {
  name                 = azurerm_managed_disk.test.name
  location             = azurerm_managed_disk.test.location
  resource_group_name  = azurerm_managed_disk.test.resource_group_name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "101"
  disk_mbps_read_write = "10"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, r.create_withUltraSSD(data))
}

func (ManagedDiskResource) diskEncryptionSetDependencies(data acceptance.TestData) string {
	// whilst this is in Preview it's only supported in: West Central US, Canada Central, North Europe
	// TODO: switch back to default location
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  soft_delete_enabled         = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "create",
    "delete",
    "get",
    "purge",
    "update",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
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
    "get",
    "wrapkey",
    "unwrapkey",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}
`, data.RandomInteger, location, data.RandomString, data.RandomInteger)
}

func (r ManagedDiskResource) diskEncryptionSetEncrypted(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                   = "acctestd-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  storage_account_type   = "Standard_LRS"
  create_option          = "Empty"
  disk_size_gb           = 1
  disk_encryption_set_id = azurerm_disk_encryption_set.test.id

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, r.diskEncryptionSetDependencies(data), data.RandomInteger)
}

func (r ManagedDiskResource) diskEncryptionSetUnencrypted(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 1

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, r.diskEncryptionSetDependencies(data), data.RandomInteger)
}

func (r ManagedDiskResource) managedDiskAttached(data acceptance.TestData, diskSize int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = %d
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, r.templateAttached(data), data.RandomInteger, diskSize)
}

func (r ManagedDiskResource) tierUpdateWhileAttached(data acceptance.TestData, tier string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
  tier                 = "%s"
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, r.templateAttached(data), data.RandomInteger, tier)
}

func (r ManagedDiskResource) storageTypeUpdateWhilstAttached(data acceptance.TestData, storageAccountType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctestdisk-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "%s"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, r.templateAttached(data), data.RandomInteger, storageAccountType)
}

func (ManagedDiskResource) templateAttached(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestvm-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_D2s_v3"
  admin_username                  = "adminuser"
  admin_password                  = "Password1234!"
  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagedDisk_networkPolicy_create(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                  = "acctestd-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  storage_account_type  = "Standard_LRS"
  create_option         = "Empty"
  disk_size_gb          = "4"
  zones                 = ["1"]
  network_access_policy = "DenyAll"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_networkPolicy_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                  = "acctestd-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  storage_account_type  = "Standard_LRS"
  create_option         = "Empty"
  disk_size_gb          = "4"
  zones                 = ["1"]
  network_access_policy = "DenyAll"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagedDisk_networkPolicy_create_withAllowPrivate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_disk_access" "test" {
  name                = "accda%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "staging"
  }
}

resource "azurerm_managed_disk" "test" {
  name                  = "acctestd-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  storage_account_type  = "Standard_LRS"
  create_option         = "Empty"
  disk_size_gb          = "4"
  zones                 = ["1"]
  network_access_policy = "AllowPrivate"
  disk_access_id        = azurerm_disk_access.test.id

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagedDisk_networkPolicy_update_withAllowPrivate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_disk_access" "test" {
  name                = "accda%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "staging"
  }
}

resource "azurerm_managed_disk" "test" {
  name                  = "acctestd-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  storage_account_type  = "Standard_LRS"
  create_option         = "Empty"
  disk_size_gb          = "4"
  zones                 = ["1"]
  network_access_policy = "AllowPrivate"
  disk_access_id        = azurerm_disk_access.test.id

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
