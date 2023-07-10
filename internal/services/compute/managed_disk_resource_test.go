// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedDiskResource struct{}

func TestAccManagedDisk_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

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

func TestAccManagedDisk_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zeroGbFromPlatformImage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectNonEmptyPlan: true, // since the `disk_size_gb` will have changed
		},
	})
}

func TestAccManagedDisk_import(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_copy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.copy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_fromPlatformImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.platformImage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_fromGalleryImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.galleryImage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_upload(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upload(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upload_size_bytes").HasValue("21475885568"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
				check.That(data.ResourceName).Key("tags.cost-center").HasValue("ops"),
				check.That(data.ResourceName).Key("disk_size_gb").HasValue("1"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue(string(disks.DiskStorageAccountTypesStandardLRS)),
			),
		},
		{
			Config: r.empty_updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
				check.That(data.ResourceName).Key("disk_size_gb").HasValue("2"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue(string(disks.DiskStorageAccountTypesPremiumLRS)),
			),
		},
	})
}

func TestAccManagedDisk_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.0.secret_url").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.0.source_vault_id").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.0.key_url").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.0.source_vault_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.encryptionUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_importEmpty_withZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty_withZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_create_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withUltraSSD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_update_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withUltraSSD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_iops_read_write").HasValue("101"),
				check.That(data.ResourceName).Key("disk_mbps_read_write").HasValue("10"),
			),
		},
		{
			Config: r.update_withUltraSSD(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withUltraSSD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.import_withUltraSSD),
	})
}

func TestAccManagedDisk_diskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diskEncryptionSetEncrypted(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_diskEncryptionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diskEncryptionSetUnencrypted(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.diskEncryptionSetEncrypted(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_attachedDiskUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedDiskAttached(data, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.managedDiskAttached(data, 20),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageTypeUpdateWhilstAttached(data, "Standard_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.storageTypeUpdateWhilstAttached(data, "Premium_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_attachedTierUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tierUpdateWhileAttached(data, "P10"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("P10"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tierUpdateWhileAttached(data, "P20"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("P20"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_networkPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkPolicy_create(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_networkPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkPolicy_create(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "DenyAll"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkPolicy_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "DenyAll"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_networkPolicy_create_withAllowPrivate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkPolicy_create_withAllowPrivate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_networkPolicy_update_withAllowPrivate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkPolicy_create_withAllowPrivate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "AllowPrivate"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkPolicy_update_withAllowPrivate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "network_access_policy", "AllowPrivate"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_publicNetworkAccessDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_publicNetworkAccessDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_publicNetworkAccessUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_update_withMaxShares(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withMaxShares(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update_withMaxShares(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_create_withLogicalSectorSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withLogicalSectorSize(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_create_withTrustedLaunchEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withTrustedLaunchEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_create_withSecurityType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withSecurityType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_create_withSecureVMDiskEncryptionSetId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withSecureVMDiskEncryptionSetId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_update_withIOpsReadOnlyAndMBpsReadOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withIOpsReadOnlyAndMBpsReadOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update_withIOpsReadOnlyAndMBpsReadOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccManagedDisk_create_withOnDemandBurstingEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withOnDemandBurstingEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_update_withOnDemandBurstingEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update_withOnDemandBurstingEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_create_withHyperVGeneration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.create_withHyperVGeneration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.edgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_storageAccountType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountType(data, "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_onlineLiveResize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.onlineLiveResizing(data, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.onlineLiveResizing(data, 20),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClientForResource(r.checkLinuxVirtualMachineWasNotRestarted, "azurerm_linux_virtual_machine.test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_detachFromVM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.detachFromVm(data, 4000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.detachFromVm(data, 5000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_premiumV2WithIOpsReadWriteAndMBpsReadWrite(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumV2WithIOpsReadWriteAndMBpsReadWrite(data, "westeurope", 3000, 125),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.premiumV2WithIOpsReadWriteAndMBpsReadWrite(data, "westeurope", 4000, 200),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedDisk_premiumV2WithIOpsReadOnlyAndMBpsReadOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumV2WithIOpsReadOnlyAndMBpsReadOnly(data, "westeurope", 3000, 125),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.premiumV2WithIOpsReadOnlyAndMBpsReadOnly(data, "westeurope", 4000, 126),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ManagedDiskResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := disks.ParseDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.DisksClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Managed Disk %q", id.String())
	}

	return utils.Bool(resp.Model != nil), nil
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
  zone                 = "1"

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
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

# NOTE: using the legacy vm resource since this test requires an unmanaged disk
resource "azurerm_virtual_machine" "test" {
  name                          = "acctestvm-%[1]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  network_interface_ids         = [azurerm_network_interface.test.id]
  vm_size                       = "Standard_F2"
  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "acctestvm-%[1]d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
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
  name                 = "acctestd-%[1]d"
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

  depends_on = [
    azurerm_virtual_machine.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary)
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
  hyper_v_generation   = "V1"
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

func (r ManagedDiskResource) galleryImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

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

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_linux_virtual_machine.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    "foo" = "bar"
  }
}

resource "azurerm_managed_disk" "test" {
  name                       = "acctestd-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  os_type                    = "Linux"
  hyper_v_generation         = "V1"
  create_option              = "FromImage"
  gallery_image_reference_id = azurerm_shared_image_version.test.id
  storage_account_type       = "Standard_LRS"
}
`, LinuxVirtualMachineResource{}.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ManagedDiskResource) upload(data acceptance.TestData) string {
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
  create_option        = "Upload"
  storage_account_type = "Standard_LRS"
  upload_size_bytes    = 21475885568
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) encryptionTemplate(data acceptance.TestData) string {
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
  sku_name            = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func (r ManagedDiskResource) encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  encryption_settings {
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
`, r.encryptionTemplate(data), data.RandomInteger)
}

func (r ManagedDiskResource) encryptionUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test2" {
  name                = "acctestkv2-%[2]s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
  sku_name            = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

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

resource "azurerm_key_vault_secret" "test2" {
  name         = "secret2-%[2]s"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test2.id
}

resource "azurerm_key_vault_key" "test2" {
  name         = "key2-%[2]s"
  key_vault_id = azurerm_key_vault.test2.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%[3]d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  encryption_settings {
    disk_encryption_key {
      secret_url      = "${azurerm_key_vault_secret.test2.id}"
      source_vault_id = "${azurerm_key_vault.test2.id}"
    }

    key_encryption_key {
      key_url         = "${azurerm_key_vault_key.test2.id}"
      source_vault_id = "${azurerm_key_vault.test2.id}"
    }
  }

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, r.encryptionTemplate(data), data.RandomString, data.RandomInteger)
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
  zone                 = "1"

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
  zone                 = "1"

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
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
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

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
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
  address_prefixes     = ["10.0.2.0/24"]
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

func (ManagedDiskResource) networkPolicy_create(data acceptance.TestData) string {
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
  zone                  = "1"
  network_access_policy = "DenyAll"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) networkPolicy_update(data acceptance.TestData) string {
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
  zone                  = "1"
  network_access_policy = "DenyAll"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) networkPolicy_create_withAllowPrivate(data acceptance.TestData) string {
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
  zone                  = "1"
  network_access_policy = "AllowPrivate"
  disk_access_id        = azurerm_disk_access.test.id

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ManagedDiskResource) networkPolicy_update_withAllowPrivate(data acceptance.TestData) string {
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
  zone                  = "1"
  network_access_policy = "AllowPrivate"
  disk_access_id        = azurerm_disk_access.test.id

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ManagedDiskResource) publicNetworkAccessDefault(data acceptance.TestData) string {
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
  disk_size_gb         = "4"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) publicNetworkAccessDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                          = "acctestd-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  storage_account_type          = "Standard_LRS"
  create_option                 = "Empty"
  disk_size_gb                  = "4"
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) create_withMaxShares(data acceptance.TestData) string {
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
  disk_size_gb         = "256"
  max_shares           = 2
  zone                 = "1"
  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) create_withLogicalSectorSize(data acceptance.TestData) string {
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
  disk_size_gb         = "256"
  logical_sector_size  = 512
  zone                 = "1"
  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) update_withMaxShares(data acceptance.TestData) string {
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
  disk_size_gb         = "1024"
  max_shares           = 5
  zone                 = "1"
  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) create_withTrustedLaunchEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "18_04-LTS-gen2"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}


resource "azurerm_managed_disk" "test" {
  name                   = "acctestd-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  os_type                = "Linux"
  hyper_v_generation     = "V2"
  create_option          = "FromImage"
  image_reference_id     = data.azurerm_platform_image.test.id
  storage_account_type   = "Standard_LRS"
  trusted_launch_enabled = true
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) create_withSecurityType(data acceptance.TestData) string {
	// Confidential VM has limited region support
	data.Locations.Primary = "northeurope"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-confidential-vm-focal"
  sku       = "20_04-lts-cvm"
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
  hyper_v_generation   = "V2"
  create_option        = "FromImage"
  image_reference_id   = data.azurerm_platform_image.test.id
  storage_account_type = "Standard_LRS"

  security_type = "ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey"
}
`, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) create_withSecureVMDiskEncryptionSetId(data acceptance.TestData) string {
	// Confidential VM has limited region support
	data.Locations.Primary = "northeurope"
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

data "azurerm_platform_image" "test" {
  location  = "%[1]s"
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-confidential-vm-focal"
  sku       = "20_04-lts-cvm"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%[3]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "premium"
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  enabled_for_disk_encryption = true
  soft_delete_retention_days  = 7
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
  key_type     = "RSA-HSM"
  key_size     = 2048
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
  depends_on = [azurerm_key_vault_access_policy.service-principal]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id
  encryption_type     = "ConfidentialVmEncryptedWithCustomerKey"
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

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  os_type              = "Linux"
  hyper_v_generation   = "V2"
  create_option        = "FromImage"
  image_reference_id   = data.azurerm_platform_image.test.id
  storage_account_type = "Standard_LRS"

  security_type                    = "ConfidentialVM_DiskEncryptedWithCustomerKey"
  secure_vm_disk_encryption_set_id = azurerm_disk_encryption_set.test.id

  depends_on = [
    azurerm_key_vault_access_policy.disk-encryption,
  ]
}

`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (ManagedDiskResource) create_withIOpsReadOnlyAndMBpsReadOnly(data acceptance.TestData) string {
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
  disk_iops_read_only  = "101"
  disk_mbps_read_only  = "10"
  max_shares           = "2"
  zone                 = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) update_withIOpsReadOnlyAndMBpsReadOnly(data acceptance.TestData) string {
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
  disk_iops_read_only  = "102"
  disk_mbps_read_only  = "11"
  max_shares           = "2"
  zone                 = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) create_withOnDemandBurstingEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_managed_disk" "test" {
  name                       = "acctestd-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  storage_account_type       = "Premium_LRS"
  create_option              = "Empty"
  disk_size_gb               = "1024"
  on_demand_bursting_enabled = true
  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) update_withOnDemandBurstingEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_managed_disk" "test" {
  name                       = "acctestd-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  storage_account_type       = "Premium_LRS"
  create_option              = "Empty"
  disk_size_gb               = "1024"
  on_demand_bursting_enabled = true
  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) create_withHyperVGeneration(data acceptance.TestData) string {
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
  disk_size_gb         = "1024"
  zone                 = "1"
  hyper_v_generation   = "V2"
  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) edgeZone(data acceptance.TestData) string {
	// @tombuildsstuff: WestUS has an edge zone available - so hard-code to that for now
	data.Locations.Primary = "westus"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
  edge_zone            = data.azurerm_extended_locations.test.extended_locations[0]

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) storageAccountType(data acceptance.TestData, location string) string {
	// Limited regional availability for some storage account type
	data.Locations.Primary = location

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
  storage_account_type = "PremiumV2_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
  logical_sector_size  = 512

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskResource) onlineLiveResizing(data acceptance.TestData, diskSizeGB int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  random_integer   = %[1]d
  primary_location = %[2]q
  first_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.random_integer}"
  location = local.primary_location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-${local.random_integer}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-${local.random_integer}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

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

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-${local.random_integer}"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = %[3]d
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "10"
  caching            = "ReadWrite"
}
`, data.RandomInteger, data.Locations.Primary, diskSizeGB)
}

func (ManagedDiskResource) detachFromVm(data acceptance.TestData, diskSizeGB int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  random_integer   = %[1]d
  primary_location = %[2]q
  first_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.random_integer}"
  location = local.primary_location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-${local.random_integer}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-${local.random_integer}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

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

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-${local.random_integer}"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = %[3]d
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "10"
  caching            = "None"
}
`, data.RandomInteger, data.Locations.Primary, diskSizeGB)
}

func (ManagedDiskResource) checkLinuxVirtualMachineWasNotRestarted(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	activityLogsClient := client.Monitor.ActivityLogsClient
	filter := fmt.Sprintf("eventTimestamp ge '%s' and resourceUri eq '%s'", time.Now().Add(-1*time.Hour).Format(time.RFC3339), state.ID)
	logs, err := activityLogsClient.ListComplete(ctx, filter, "")
	if err != nil {
		return fmt.Errorf("retrieving activity logs for Virtual Machine %q: %+v", state.ID, err)
	}

	wasShutDown := false
	for logs.NotDone() {
		val := logs.Value()
		if val.Authorization != nil && val.Authorization.Action != nil {
			if strings.EqualFold(*val.Authorization.Action, "Microsoft.Compute/virtualMachines/powerOff/action") {
				wasShutDown = true
				break
			}
		}

		if err := logs.NextWithContext(ctx); err != nil {
			return fmt.Errorf("listing the next page of results: %+v", err)
		}
	}

	if wasShutDown {
		return fmt.Errorf("the Virtual Machine %q was shut down when resizing the Data Disk when it shouldn't have been", state.ID)
	}

	return nil
}

func (ManagedDiskResource) premiumV2WithIOpsReadWriteAndMBpsReadWrite(data acceptance.TestData, location string, diskIOpsReadWrite, diskMBpsReadWrite int) string {
	// Limited regional availability for some storage account type
	data.Locations.Primary = location

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
  storage_account_type = "PremiumV2_LRS"
  create_option        = "Empty"
  disk_size_gb         = 256
  disk_iops_read_write = %d
  disk_mbps_read_write = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, diskIOpsReadWrite, diskMBpsReadWrite)
}

func (ManagedDiskResource) premiumV2WithIOpsReadOnlyAndMBpsReadOnly(data acceptance.TestData, location string, diskIOpsReadOnly, diskMBpsReadOnly int) string {
	// Limited regional availability for some storage account type
	data.Locations.Primary = location

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
  storage_account_type = "PremiumV2_LRS"
  create_option        = "Empty"
  disk_size_gb         = 256
  disk_iops_read_only  = %d
  disk_mbps_read_only  = %d
  max_shares           = 5
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, diskIOpsReadOnly, diskMBpsReadOnly)
}
