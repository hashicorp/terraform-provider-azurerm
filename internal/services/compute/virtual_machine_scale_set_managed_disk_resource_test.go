// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualMachineScaleSetManagedDiskResource struct{}

func TestAccVirtualMachineScaleSetManagedDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetManagedDisk_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_machine_scale_set_managed_disk"),
		},
	})
}

func TestAccVirtualMachineScaleSetManagedDisk_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

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

func TestAccVirtualMachineScaleSetManagedDisk_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetManagedDisk_osTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osType(data, "Linux"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("os_type").HasValue("Linux"),
			),
		},
		data.ImportStep(),
		{
			Config: r.osType(data, "Windows"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("os_type").HasValue("Windows"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetManagedDisk_diskSizeCannotShrink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSize(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.withSize(data, 1),
			ExpectError: regexp.MustCompile("`disk_size_gb` can only be increased"),
		},
	})
}

func TestAccVirtualMachineScaleSetManagedDisk_validationErrors(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.diskIopsWithoutUltraSku(data),
			ExpectError: regexp.MustCompile("are only available for `UltraSSD_LRS` and `PremiumV2_LRS`"),
		},
		{
			Config:      r.tierWithoutPremiumSku(data),
			ExpectError: regexp.MustCompile("`tier` can only be specified when `storage_account_type`"),
		},
		{
			Config:      r.onDemandBurstingWithoutPremiumSku(data),
			ExpectError: regexp.MustCompile("`on_demand_bursting_enabled` can only be set to `true` when `storage_account_type`"),
		},
		{
			Config:      r.diskAccessWithoutPrivatePolicy(data),
			ExpectError: regexp.MustCompile("`disk_access_id` is only available when `network_access_policy` is set to `AllowPrivate`"),
		},
	})
}

func TestAccVirtualMachineScaleSetManagedDisk_ultraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ultraSSD(data, 101, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_iops_read_write").HasValue("101"),
				check.That(data.ResourceName).Key("disk_mbps_read_write").HasValue("10"),
			),
		},
		data.ImportStep(),
		{
			Config: r.ultraSSD(data, 102, 11),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_iops_read_write").HasValue("102"),
				check.That(data.ResourceName).Key("disk_mbps_read_write").HasValue("11"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetManagedDisk_encryptionSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_settings.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (VirtualMachineScaleSetManagedDiskResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseManagedDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.DisksClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (VirtualMachineScaleSetManagedDiskResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualMachineScaleSetManagedDiskResource) empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "Standard_LRS"

  creation {
    option = "Empty"
  }

  disk_size_gb = 1
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetManagedDiskResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "import" {
  name                 = azurerm_virtual_machine_scale_set_managed_disk.test.name
  resource_group_name  = azurerm_virtual_machine_scale_set_managed_disk.test.resource_group_name
  location             = azurerm_virtual_machine_scale_set_managed_disk.test.location
  storage_account_type = azurerm_virtual_machine_scale_set_managed_disk.test.storage_account_type

  creation {
    option = "Empty"
  }

  disk_size_gb = azurerm_virtual_machine_scale_set_managed_disk.test.disk_size_gb
}
`, r.empty(data))
}

func (r VirtualMachineScaleSetManagedDiskResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "Standard_LRS"

  creation {
    option = "Empty"
  }

  disk_size_gb = 2

  tags = {
    environment = "staging"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetManagedDiskResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "Standard_LRS"
  zone                 = "1"

  creation {
    option = "Empty"
  }

  disk_size_gb                      = 4
  os_type                           = "Linux"
  network_access_policy             = "DenyAll"
  public_network_access_enabled     = false
  data_access_auth_mode             = "AzureActiveDirectory"
  optimized_frequent_attach_enabled = true

  tags = {
    environment = "staging"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetManagedDiskResource) osType(data acceptance.TestData, osType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "Standard_LRS"
  os_type              = "%s"

  creation {
    option = "Empty"
  }

  disk_size_gb = 1
}
`, r.template(data), data.RandomInteger, osType)
}

func (r VirtualMachineScaleSetManagedDiskResource) withSize(data acceptance.TestData, sizeGb int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "Standard_LRS"

  creation {
    option = "Empty"
  }

  disk_size_gb = %d
}
`, r.template(data), data.RandomInteger, sizeGb)
}

func (r VirtualMachineScaleSetManagedDiskResource) ultraSSD(data acceptance.TestData, iops int, mbps int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "UltraSSD_LRS"
  zone                 = "1"

  creation {
    option              = "Empty"
    logical_sector_size = 512
  }

  disk_size_gb         = 4
  disk_iops_read_write = %d
  disk_mbps_read_write = %d
}
`, r.template(data), data.RandomInteger, iops, mbps)
}

func (r VirtualMachineScaleSetManagedDiskResource) encryptionSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "Standard_LRS"

  creation {
    option = "Empty"
  }

  disk_size_gb = 1

  encryption_settings {
    disk_encryption_key {
      secret_url      = azurerm_key_vault_secret.test.id
      source_vault_id = azurerm_key_vault.test.id
    }

    key_encryption_key {
      key_url         = azurerm_key_vault_key.test.id
      source_vault_id = azurerm_key_vault.test.id
    }
  }
}
`, r.encryptionTemplate(data), data.RandomInteger)
}

func (VirtualMachineScaleSetManagedDiskResource) encryptionTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults       = false
      purge_soft_delete_on_destroy          = false
      purge_soft_deleted_keys_on_destroy    = false
      purge_soft_deleted_secrets_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }

  enabled_for_disk_encryption = true
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func (VirtualMachineScaleSetManagedDiskResource) diskIopsWithoutUltraSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = "acctestRG-%d"
  location             = %q
  storage_account_type = "Standard_LRS"
  disk_iops_read_write = 200

  creation {
    option = "Empty"
  }

  disk_size_gb = 1
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (VirtualMachineScaleSetManagedDiskResource) tierWithoutPremiumSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = "acctestRG-%d"
  location             = %q
  storage_account_type = "Standard_LRS"
  tier                 = "P10"

  creation {
    option = "Empty"
  }

  disk_size_gb = 1
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (VirtualMachineScaleSetManagedDiskResource) onDemandBurstingWithoutPremiumSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                       = "acctestd-%d"
  resource_group_name        = "acctestRG-%d"
  location                   = %q
  storage_account_type       = "Standard_LRS"
  on_demand_bursting_enabled = true

  creation {
    option = "Empty"
  }

  disk_size_gb = 1024
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (VirtualMachineScaleSetManagedDiskResource) diskAccessWithoutPrivatePolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%d"
  resource_group_name  = "acctestRG-%d"
  location             = %q
  storage_account_type = "Standard_LRS"
  disk_access_id       = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-%d/providers/Microsoft.Compute/diskAccesses/acctest-diskaccess"

  creation {
    option = "Empty"
  }

  disk_size_gb = 1
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
