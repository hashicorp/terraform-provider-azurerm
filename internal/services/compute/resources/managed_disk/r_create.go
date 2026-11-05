// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedDiskCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := commonids.NewManagedDiskID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Managed Disk %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_managed_disk", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	createOption := disks.DiskCreateOption(d.Get("create_option").(string))
	storageAccountType := d.Get("storage_account_type").(string)
	osType := disks.OperatingSystemTypes(d.Get("os_type").(string))
	maxShares := d.Get("max_shares").(int)

	t := d.Get("tags").(map[string]any)
	skuName := disks.DiskStorageAccountTypes(storageAccountType)

	props := &disks.DiskProperties{
		CreationData: disks.CreationData{
			CreateOption:    createOption,
			PerformancePlus: pointer.To(d.Get("performance_plus_enabled").(bool)),
		},
		OptimizedForFrequentAttach: pointer.To(d.Get("optimized_frequent_attach_enabled").(bool)),
		OsType:                     &osType,
		Encryption: &disks.Encryption{
			Type: pointer.To(disks.EncryptionTypeEncryptionAtRestWithPlatformKey),
		},
	}

	diskSizeGB := d.Get("disk_size_gb").(int)
	if diskSizeGB != 0 {
		props.DiskSizeGB = pointer.To(int64(diskSizeGB))
	}

	if maxShares != 0 {
		props.MaxShares = pointer.To(int64(maxShares))
	}

	// lintignore:R019 // deliberate subset: the else-if lists exactly the UltraSSD/PremiumV2-only fields so changes to them can be rejected on other SKUs
	if storageAccountType == string(disks.DiskStorageAccountTypesUltraSSDLRS) || storageAccountType == string(disks.DiskStorageAccountTypesPremiumVTwoLRS) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			props.DiskIOPSReadWrite = pointer.To(int64(v.(int)))
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			props.DiskMBpsReadWrite = pointer.To(int64(v.(int)))
		}

		if v, ok := d.GetOk("disk_iops_read_only"); ok {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_iops_read_only is only available for UltraSSD disks and PremiumV2 disks with shared disk enabled")
			}

			props.DiskIOPSReadOnly = pointer.To(int64(v.(int)))
		}

		if v, ok := d.GetOk("disk_mbps_read_only"); ok {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_mbps_read_only is only available for UltraSSD disks and PremiumV2 disks with shared disk enabled")
			}

			props.DiskMBpsReadOnly = pointer.To(int64(v.(int)))
		}

		if v, ok := d.GetOk("logical_sector_size"); ok {
			props.CreationData.LogicalSectorSize = pointer.To(int64(v.(int)))
		}
	} else if d.HasChanges(
		"disk_iops_read_write",
		"disk_mbps_read_write",
		"disk_iops_read_only",
		"disk_mbps_read_only",
		"logical_sector_size",
	) {
		return fmt.Errorf("[ERROR] disk_iops_read_write, disk_mbps_read_write, disk_iops_read_only, disk_mbps_read_only and logical_sector_size are only available for UltraSSD disks and PremiumV2 disks")
	}

	if createOption == disks.DiskCreateOptionImport || createOption == disks.DiskCreateOptionImportSecure {
		sourceUri := d.Get("source_uri").(string)
		if sourceUri == "" {
			return fmt.Errorf("`source_uri` must be specified when `create_option` is set to `Import` or `ImportSecure`")
		}

		storageAccountId := d.Get("storage_account_id").(string)
		if storageAccountId == "" {
			return fmt.Errorf("`storage_account_id` must be specified when `create_option` is set to `Import` or `ImportSecure`")
		}

		props.CreationData.StorageAccountId = pointer.To(storageAccountId)
		props.CreationData.SourceUri = pointer.To(sourceUri)
	}
	if createOption == disks.DiskCreateOptionCopy || createOption == disks.DiskCreateOptionRestore {
		sourceResourceId := d.Get("source_resource_id").(string)
		if sourceResourceId == "" {
			return fmt.Errorf("`source_resource_id` must be specified when `create_option` is set to `Copy` or `Restore`")
		}

		props.CreationData.SourceResourceId = pointer.To(sourceResourceId)
	}
	if createOption == disks.DiskCreateOptionFromImage {
		if imageReferenceId := d.Get("image_reference_id").(string); imageReferenceId != "" {
			props.CreationData.ImageReference = &disks.ImageDiskReference{
				Id: pointer.To(imageReferenceId),
			}
		} else if galleryImageReferenceId := d.Get("gallery_image_reference_id").(string); galleryImageReferenceId != "" {
			props.CreationData.GalleryImageReference = &disks.ImageDiskReference{
				Id: pointer.To(galleryImageReferenceId),
			}
		} else {
			return fmt.Errorf("`image_reference_id` or `gallery_image_reference_id` must be specified when `create_option` is set to `FromImage`")
		}
	}

	if createOption == disks.DiskCreateOptionUpload {
		if uploadSizeBytes := d.Get("upload_size_bytes").(int); uploadSizeBytes != 0 {
			props.CreationData.UploadSizeBytes = pointer.To(int64(uploadSizeBytes))
		} else {
			return fmt.Errorf("`upload_size_bytes` must be specified when `create_option` is set to `Upload`")
		}
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		props.EncryptionSettingsCollection = helpers.ExpandManagedDiskEncryptionSettings(v.([]any))
	}

	if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
		encryptionType, err := helpers.RetrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
		if err != nil {
			return err
		}

		props.Encryption = &disks.Encryption{
			Type:                encryptionType,
			DiskEncryptionSetId: pointer.To(diskEncryptionSetId),
		}
	}

	props.NetworkAccessPolicy = pointer.ToEnum[disks.NetworkAccessPolicy](d.Get("network_access_policy").(string))

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		switch {
		case *props.NetworkAccessPolicy == disks.NetworkAccessPolicyAllowPrivate:
			props.DiskAccessId = pointer.To(diskAccessID)
		case diskAccessID != "" && *props.NetworkAccessPolicy != disks.NetworkAccessPolicyAllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			props.DiskAccessId = nil
		}
	}

	if d.Get("public_network_access_enabled").(bool) {
		props.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessEnabled)
	} else {
		props.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessDisabled)
	}

	if tier := d.Get("tier").(string); tier != "" {
		if storageAccountType != string(disks.DiskStorageAccountTypesPremiumZRS) && storageAccountType != string(disks.DiskStorageAccountTypesPremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		props.Tier = &tier
	}

	if d.Get("trusted_launch_enabled").(bool) {
		props.SecurityProfile = &disks.DiskSecurityProfile{
			SecurityType: pointer.To(disks.DiskSecurityTypesTrustedLaunch),
		}

		switch createOption {
		case disks.DiskCreateOptionFromImage:
		case disks.DiskCreateOptionImport:
		case disks.DiskCreateOptionImportSecure:
		default:
			return fmt.Errorf("trusted_launch_enabled cannot be set to true with create_option %q. Supported Create Options when Trusted Launch is enabled are `FromImage`, `Import`, `ImportSecure`", createOption)
		}
	}

	securityType := d.Get("security_type").(string)
	secureVMDiskEncryptionId := d.Get("secure_vm_disk_encryption_set_id")
	if securityType != "" {
		if d.Get("trusted_launch_enabled").(bool) {
			return fmt.Errorf("`security_type` cannot be specified when `trusted_launch_enabled` is set to `true`")
		}

		switch createOption {
		case disks.DiskCreateOptionFromImage:
		case disks.DiskCreateOptionImport:
		case disks.DiskCreateOptionImportSecure:
		default:
			return fmt.Errorf("`security_type` can only be specified when `create_option` is set to `FromImage`, `Import` or `ImportSecure`")
		}

		if disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey == disks.DiskSecurityTypes(securityType) && secureVMDiskEncryptionId == "" {
			return fmt.Errorf("`secure_vm_disk_encryption_set_id` must be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`")
		}

		props.SecurityProfile = &disks.DiskSecurityProfile{
			SecurityType: pointer.ToEnum[disks.DiskSecurityTypes](securityType),
		}
	}

	if secureVMDiskEncryptionId != "" {
		if disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey != disks.DiskSecurityTypes(securityType) {
			return fmt.Errorf("`secure_vm_disk_encryption_set_id` can only be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`")
		}
		props.SecurityProfile.SecureVMDiskEncryptionSetId = pointer.To(secureVMDiskEncryptionId.(string))
	}

	if d.Get("on_demand_bursting_enabled").(bool) {
		switch storageAccountType {
		case string(disks.DiskStorageAccountTypesPremiumLRS):
		case string(disks.DiskStorageAccountTypesPremiumZRS):
		default:
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}

		if diskSizeGB != 0 && diskSizeGB <= 512 {
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `disk_size_gb` is larger than 512GB")
		}

		props.BurstingEnabled = pointer.To(true)
	}

	if v, ok := d.GetOk("hyper_v_generation"); ok {
		props.HyperVGeneration = pointer.ToEnum[disks.HyperVGeneration](v.(string))
	}

	createDisk := disks.Disk{
		Name:             &name,
		ExtendedLocation: helpers.ExpandManagedDiskEdgeZone(d.Get("edge_zone").(string)),
		Location:         location,
		Properties:       props,
		Sku: &disks.DiskSku{
			Name: &skuName,
		},
		Tags: tags.Expand(t),
	}

	if zone, ok := d.GetOk("zone"); ok {
		createDisk.Zones = &[]string{
			zone.(string),
		}
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, createDisk, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceManagedDiskRead(d, meta)
}
