// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func resourceManagedDisk() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedDiskCreate,
		Read:   resourceManagedDiskRead,
		Update: resourceManagedDiskUpdate,
		Delete: resourceManagedDiskDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ManagedDiskV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := disks.ParseDiskID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"storage_account_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disks.DiskStorageAccountTypesStandardLRS),
					string(disks.DiskStorageAccountTypesStandardSSDZRS),
					string(disks.DiskStorageAccountTypesPremiumLRS),
					string(disks.DiskStorageAccountTypesPremiumVTwoLRS),
					string(disks.DiskStorageAccountTypesPremiumZRS),
					string(disks.DiskStorageAccountTypesStandardSSDLRS),
					string(disks.DiskStorageAccountTypesUltraSSDLRS),
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"create_option": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disks.DiskCreateOptionCopy),
					string(disks.DiskCreateOptionEmpty),
					string(disks.DiskCreateOptionFromImage),
					string(disks.DiskCreateOptionImport),
					string(disks.DiskCreateOptionImportSecure),
					string(disks.DiskCreateOptionRestore),
					string(disks.DiskCreateOptionUpload),
				}, false),
			},

			"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

			"logical_sector_size": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					512,
					4096,
				}),
				Computed: true,
			},

			"optimized_frequent_attach_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"performance_plus_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"source_uri": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true, // Not supported by disk update
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"image_reference_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"gallery_image_reference_id"},
			},

			"gallery_image_reference_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validate.SharedImageVersionID,
				ConflictsWith: []string{"image_reference_id"},
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disks.OperatingSystemTypesWindows),
					string(disks.OperatingSystemTypesLinux),
				}, false),
			},

			"disk_size_gb": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ManagedDiskSizeGB,
			},

			"upload_size_bytes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"disk_iops_read_write": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"disk_mbps_read_write": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"disk_iops_read_only": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"disk_mbps_read_only": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
				//       https://github.com/Azure/azure-rest-api-specs/issues/8132
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.DiskEncryptionSetID,
				ConflictsWith:    []string{"secure_vm_disk_encryption_set_id"},
			},

			"encryption_settings": encryptionSettingsSchema(),

			"network_access_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disks.NetworkAccessPolicyAllowAll),
					string(disks.NetworkAccessPolicyAllowPrivate),
					string(disks.NetworkAccessPolicyDenyAll),
				}, false),
			},
			"disk_access_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
				//       https://github.com/Azure/azure-rest-api-specs/issues/14192
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     diskaccesses.ValidateDiskAccessID,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_shares": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2, 10),
			},

			"trusted_launch_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"secure_vm_disk_encryption_set_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validate.DiskEncryptionSetID,
				ConflictsWith: []string{"disk_encryption_set_id"},
			},

			"security_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disks.DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey),
					string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey),
					string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey),
				}, false),
			},

			"hyper_v_generation": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true, // Not supported by disk update
				ValidateFunc: validation.StringInSlice([]string{
					string(disks.HyperVGenerationVOne),
					string(disks.HyperVGenerationVTwo),
				}, false),
			},

			"on_demand_bursting_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"zone": commonschema.ZoneSingleOptionalForceNew(),

			"tags": commonschema.Tags(),
		},

		// Encryption Settings cannot be disabled once enabled
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("encryption_settings", func(ctx context.Context, old, new, meta interface{}) bool {
				if !features.FourPointOhBeta() {
					return false
				}
				return len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0
			}),
		),
	}
}

func resourceManagedDiskCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := disks.NewDiskID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
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

	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := disks.DiskCreateOption(d.Get("create_option").(string))
	storageAccountType := d.Get("storage_account_type").(string)
	osType := disks.OperatingSystemTypes(d.Get("os_type").(string))
	maxShares := d.Get("max_shares").(int)

	t := d.Get("tags").(map[string]interface{})
	skuName := disks.DiskStorageAccountTypes(storageAccountType)
	encryptionTypePlatformKey := disks.EncryptionTypeEncryptionAtRestWithPlatformKey

	props := &disks.DiskProperties{
		CreationData: disks.CreationData{
			CreateOption:    createOption,
			PerformancePlus: pointer.To(d.Get("performance_plus_enabled").(bool)),
		},
		OptimizedForFrequentAttach: pointer.To(d.Get("optimized_frequent_attach_enabled").(bool)),
		OsType:                     &osType,
		Encryption: &disks.Encryption{
			Type: &encryptionTypePlatformKey,
		},
	}

	diskSizeGB := d.Get("disk_size_gb").(int)
	if diskSizeGB != 0 {
		props.DiskSizeGB = utils.Int64(int64(diskSizeGB))
	}

	if maxShares != 0 {
		props.MaxShares = utils.Int64(int64(maxShares))
	}

	if storageAccountType == string(disks.DiskStorageAccountTypesUltraSSDLRS) || storageAccountType == string(disks.DiskStorageAccountTypesPremiumVTwoLRS) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			props.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int64(v.(int))
			props.DiskMBpsReadWrite = &diskMBps
		}

		if v, ok := d.GetOk("disk_iops_read_only"); ok {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_iops_read_only is only available for UltraSSD disks and PremiumV2 disks with shared disk enabled")
			}

			props.DiskIOPSReadOnly = utils.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("disk_mbps_read_only"); ok {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_mbps_read_only is only available for UltraSSD disks and PremiumV2 disks with shared disk enabled")
			}

			props.DiskMBpsReadOnly = utils.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("logical_sector_size"); ok {
			props.CreationData.LogicalSectorSize = utils.Int64(int64(v.(int)))
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") || d.HasChange("disk_iops_read_only") || d.HasChange("disk_mbps_read_only") || d.HasChange("logical_sector_size") {
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

		props.CreationData.StorageAccountId = utils.String(storageAccountId)
		props.CreationData.SourceUri = utils.String(sourceUri)
	}
	if createOption == disks.DiskCreateOptionCopy || createOption == disks.DiskCreateOptionRestore {
		sourceResourceId := d.Get("source_resource_id").(string)
		if sourceResourceId == "" {
			return fmt.Errorf("`source_resource_id` must be specified when `create_option` is set to `Copy` or `Restore`")
		}

		props.CreationData.SourceResourceId = utils.String(sourceResourceId)
	}
	if createOption == disks.DiskCreateOptionFromImage {
		if imageReferenceId := d.Get("image_reference_id").(string); imageReferenceId != "" {
			props.CreationData.ImageReference = &disks.ImageDiskReference{
				Id: utils.String(imageReferenceId),
			}
		} else if galleryImageReferenceId := d.Get("gallery_image_reference_id").(string); galleryImageReferenceId != "" {
			props.CreationData.GalleryImageReference = &disks.ImageDiskReference{
				Id: utils.String(galleryImageReferenceId),
			}
		} else {
			return fmt.Errorf("`image_reference_id` or `gallery_image_reference_id` must be specified when `create_option` is set to `FromImage`")
		}
	}

	if createOption == disks.DiskCreateOptionUpload {
		if uploadSizeBytes := d.Get("upload_size_bytes").(int); uploadSizeBytes != 0 {
			props.CreationData.UploadSizeBytes = utils.Int64(int64(uploadSizeBytes))
		} else {
			return fmt.Errorf("`upload_size_bytes` must be specified when `create_option` is set to `Upload`")
		}
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		props.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(v.([]interface{}))
	}

	if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
		encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
		if err != nil {
			return err
		}

		props.Encryption = &disks.Encryption{
			Type:                encryptionType,
			DiskEncryptionSetId: utils.String(diskEncryptionSetId),
		}
	}

	if networkAccessPolicy := d.Get("network_access_policy").(string); networkAccessPolicy != "" {
		policy := disks.NetworkAccessPolicy(networkAccessPolicy)
		props.NetworkAccessPolicy = &policy
	} else {
		allowAllPolicy := disks.NetworkAccessPolicyAllowAll
		props.NetworkAccessPolicy = &allowAllPolicy
	}

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		switch {
		case *props.NetworkAccessPolicy == disks.NetworkAccessPolicyAllowPrivate:
			props.DiskAccessId = utils.String(diskAccessID)
		case diskAccessID != "" && *props.NetworkAccessPolicy != disks.NetworkAccessPolicyAllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			props.DiskAccessId = nil
		}
	}

	if d.Get("public_network_access_enabled").(bool) {
		networkAccessEnabled := disks.PublicNetworkAccessEnabled
		props.PublicNetworkAccess = &networkAccessEnabled
	} else {
		networkAccessDisabled := disks.PublicNetworkAccessDisabled
		props.PublicNetworkAccess = &networkAccessDisabled
	}

	if tier := d.Get("tier").(string); tier != "" {
		if storageAccountType != string(disks.DiskStorageAccountTypesPremiumZRS) && storageAccountType != string(disks.DiskStorageAccountTypesPremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		props.Tier = &tier
	}

	if d.Get("trusted_launch_enabled").(bool) {
		diskSecurityTypeTrustedLaunch := disks.DiskSecurityTypesTrustedLaunch
		props.SecurityProfile = &disks.DiskSecurityProfile{
			SecurityType: &diskSecurityTypeTrustedLaunch,
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

		diskSecurityType := disks.DiskSecurityTypes(securityType)
		props.SecurityProfile = &disks.DiskSecurityProfile{
			SecurityType: &diskSecurityType,
		}
	}

	if secureVMDiskEncryptionId != "" {
		if disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey != disks.DiskSecurityTypes(securityType) {
			return fmt.Errorf("`secure_vm_disk_encryption_set_id` can only be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`")
		}
		props.SecurityProfile.SecureVMDiskEncryptionSetId = utils.String(secureVMDiskEncryptionId.(string))
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

		props.BurstingEnabled = utils.Bool(true)
	}

	if v, ok := d.GetOk("hyper_v_generation"); ok {
		hyperVGeneration := disks.HyperVGeneration(v.(string))
		props.HyperVGeneration = &hyperVGeneration
	}

	createDisk := disks.Disk{
		Name:             &name,
		ExtendedLocation: expandManagedDiskEdgeZone(d.Get("edge_zone").(string)),
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

	err := client.CreateOrUpdateThenPoll(ctx, id, createDisk)
	if err != nil {
		return fmt.Errorf("creating/updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.Model == nil {
		return fmt.Errorf("reading Managed Disk %s (Resource Group %q): ID was nil", name, resourceGroup)
	}

	d.SetId(id.ID())

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	virtualMachinesClient := meta.(*clients.Client).Compute.VirtualMachinesClient
	skusClient := meta.(*clients.Client).Compute.SkusClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	maxShares := d.Get("max_shares").(int)
	storageAccountType := d.Get("storage_account_type").(string)
	diskSizeGB := d.Get("disk_size_gb").(int)
	onDemandBurstingEnabled := d.Get("on_demand_bursting_enabled").(bool)
	shouldShutDown := false
	shouldDetach := false
	expandedDisk := compute.DataDisk{}

	id, err := disks.ParseDiskID(d.Id())
	if err != nil {
		return err
	}

	disk, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(disk.HttpResponse) {
			return fmt.Errorf("Managed Disk %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("making Read request on Azure Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	diskUpdate := disks.DiskUpdate{
		Properties: &disks.DiskUpdateProperties{},
	}

	if d.HasChange("max_shares") {
		diskUpdate.Properties.MaxShares = utils.Int64(int64(maxShares))
		var skuName disks.DiskStorageAccountTypes
		for _, v := range disks.PossibleValuesForDiskStorageAccountTypes() {
			if strings.EqualFold(storageAccountType, v) {
				skuName = disks.DiskStorageAccountTypes(v)
			}
		}
		diskUpdate.Sku = &disks.DiskSku{
			Name: &skuName,
		}
	}

	if d.HasChange("tier") {
		if storageAccountType != string(disks.DiskStorageAccountTypesPremiumZRS) && storageAccountType != string(disks.DiskStorageAccountTypesPremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		shouldShutDown = true
		tier := d.Get("tier").(string)
		diskUpdate.Properties.Tier = &tier
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		diskUpdate.Tags = tags.Expand(t)
	}

	if d.HasChange("storage_account_type") {
		shouldShutDown = true
		var skuName disks.DiskStorageAccountTypes
		for _, v := range disks.PossibleValuesForDiskStorageAccountTypes() {
			if strings.EqualFold(storageAccountType, v) {
				skuName = disks.DiskStorageAccountTypes(v)
			}
		}
		diskUpdate.Sku = &disks.DiskSku{
			Name: &skuName,
		}
	}

	if strings.EqualFold(storageAccountType, string(disks.DiskStorageAccountTypesUltraSSDLRS)) || storageAccountType == string(disks.DiskStorageAccountTypesPremiumVTwoLRS) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			diskUpdate.Properties.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int64(v.(int))
			diskUpdate.Properties.DiskMBpsReadWrite = &diskMBps
		}

		if d.HasChange("disk_iops_read_only") {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_iops_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			v := d.Get("disk_iops_read_only")
			diskUpdate.Properties.DiskIOPSReadOnly = utils.Int64(int64(v.(int)))
		}

		if d.HasChange("disk_mbps_read_only") {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_mbps_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			v := d.Get("disk_mbps_read_only")
			diskUpdate.Properties.DiskMBpsReadOnly = utils.Int64(int64(v.(int)))
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") || d.HasChange("disk_iops_read_only") || d.HasChange("disk_mbps_read_only") {
		return fmt.Errorf("[ERROR] disk_iops_read_write, disk_mbps_read_write, disk_iops_read_only and disk_mbps_read_only are only available for UltraSSD disks and PremiumV2 disks")
	}

	if d.HasChange("optimized_frequent_attach_enabled") {
		diskUpdate.Properties.OptimizedForFrequentAttach = pointer.To(d.Get("optimized_frequent_attach_enabled").(bool))
	}

	if d.HasChange("os_type") {
		operatingSystemType := disks.OperatingSystemTypes(d.Get("os_type").(string))
		diskUpdate.Properties.OsType = &operatingSystemType
	}

	if d.HasChange("disk_size_gb") {
		if oldSize, newSize := d.GetChange("disk_size_gb"); newSize.(int) > oldSize.(int) {
			canBeResizedWithoutDowntime := false
			if meta.(*clients.Client).Features.ManagedDisk.ExpandWithoutDowntime {
				diskSupportsNoDowntimeResize := determineIfDataDiskSupportsNoDowntimeResize(disk.Model, oldSize.(int), newSize.(int))

				vmSkuSupportsNoDowntimeResize, err := determineIfVirtualMachineSkuSupportsNoDowntimeResize(ctx, disk.Model.ManagedBy, virtualMachinesClient, skusClient)
				if err != nil {
					return fmt.Errorf("determining if the Virtual Machine the Disk is attached to supports no-downtime-resize: %+v", err)
				}

				// If a disk is 4 TiB or less, you can't expand it beyond 4 TiB without detaching it from the VM.
				shouldDetach = oldSize.(int) < 4096 && newSize.(int) >= 4096

				canBeResizedWithoutDowntime = *vmSkuSupportsNoDowntimeResize && *diskSupportsNoDowntimeResize
			}
			if !canBeResizedWithoutDowntime {
				log.Printf("[INFO] The %s, or the Virtual Machine that it's attached to, doesn't support no-downtime-resizing - requiring that the VM should be shutdown", *id)
				shouldShutDown = true
			}
			diskUpdate.Properties.DiskSizeGB = utils.Int64(int64(newSize.(int)))
		} else {
			return fmt.Errorf("- New size must be greater than original size. Shrinking disks is not supported on Azure")
		}
	}

	if d.HasChange("encryption_settings") {
		diskUpdate.Properties.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(d.Get("encryption_settings").([]interface{}))
	}

	if d.HasChange("disk_encryption_set_id") {
		shouldShutDown = true
		if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
			if err != nil {
				return err
			}

			diskUpdate.Properties.Encryption = &disks.Encryption{
				Type:                encryptionType,
				DiskEncryptionSetId: utils.String(diskEncryptionSetId),
			}
		} else {
			return fmt.Errorf("Once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if networkAccessPolicy := d.Get("network_access_policy").(string); networkAccessPolicy != "" {
		policy := disks.NetworkAccessPolicy(networkAccessPolicy)
		diskUpdate.Properties.NetworkAccessPolicy = &policy
	} else {
		allowAllPolicy := disks.NetworkAccessPolicyAllowAll
		diskUpdate.Properties.NetworkAccessPolicy = &allowAllPolicy
	}

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		switch {
		case *diskUpdate.Properties.NetworkAccessPolicy == disks.NetworkAccessPolicyAllowPrivate:
			diskUpdate.Properties.DiskAccessId = utils.String(diskAccessID)
		case diskAccessID != "" && *diskUpdate.Properties.NetworkAccessPolicy != disks.NetworkAccessPolicyAllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			diskUpdate.Properties.DiskAccessId = nil
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if d.Get("public_network_access_enabled").(bool) {
			networkAccessEnabled := disks.PublicNetworkAccessEnabled
			diskUpdate.Properties.PublicNetworkAccess = &networkAccessEnabled
		} else {
			networkAccessDisabled := disks.PublicNetworkAccessDisabled
			diskUpdate.Properties.PublicNetworkAccess = &networkAccessDisabled
		}
	}

	if onDemandBurstingEnabled {
		switch storageAccountType {
		case string(disks.DiskStorageAccountTypesPremiumLRS):
		case string(disks.DiskStorageAccountTypesPremiumZRS):
		default:
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}

		if diskSizeGB != 0 && diskSizeGB <= 512 {
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `disk_size_gb` is larger than 512GB")
		}
	}

	if d.HasChange("on_demand_bursting_enabled") {
		shouldShutDown = true
		diskUpdate.Properties.BurstingEnabled = utils.Bool(onDemandBurstingEnabled)
	}

	// whilst we need to shut this down, if we're not attached to anything there's no point
	if shouldShutDown && disk.Model.ManagedBy == nil {
		shouldShutDown = false
	}

	// if we are attached to a VM we bring down the VM as necessary for the operations which are not allowed while it's online
	if shouldShutDown {
		virtualMachine, err := parse.VirtualMachineID(*disk.Model.ManagedBy)
		if err != nil {
			return fmt.Errorf("parsing VMID %q for disk attachment: %+v", *disk.Model.ManagedBy, err)
		}
		// check instanceView State
		vmClient := meta.(*clients.Client).Compute.VMClient

		locks.ByName(name, VirtualMachineResourceName)
		defer locks.UnlockByName(name, VirtualMachineResourceName)

		vm, err := vmClient.Get(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, "")
		if err != nil {
			return fmt.Errorf("retrieving Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		instanceView, err := vmClient.InstanceView(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
		if err != nil {
			return fmt.Errorf("retrieving InstanceView for Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		shouldTurnBackOn := virtualMachineShouldBeStarted(instanceView)
		shouldDeallocate := true

		if instanceView.Statuses != nil {
			for _, status := range *instanceView.Statuses {
				if status.Code == nil {
					continue
				}

				// could also be the provisioning state which we're not bothered with here
				state := strings.ToLower(*status.Code)
				if !strings.HasPrefix(state, "powerstate/") {
					continue
				}

				state = strings.TrimPrefix(state, "powerstate/")
				switch strings.ToLower(state) {
				case "deallocated":
					// VM already deallocated, no shutdown and deallocation needed anymore
					shouldShutDown = false
					shouldDeallocate = false
				case "deallocating":
					// VM is deallocating
					// To make sure we do not start updating before this action has finished,
					// only skip the shutdown and send another deallocation request if shouldDeallocate == true
					shouldShutDown = false
				case "stopped":
					shouldShutDown = false
				}
			}
		}

		// Detach
		if shouldDetach {
			dataDisks := make([]compute.DataDisk, 0)
			if vm.StorageProfile != nil && vm.StorageProfile.DataDisks != nil {
				for _, dataDisk := range *vm.StorageProfile.DataDisks {
					// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
					if dataDisk.Name != nil && *dataDisk.Name != id.DiskName {
						dataDisks = append(dataDisks, dataDisk)
					} else {
						if dataDisk.Caching != compute.CachingTypesNone {
							return fmt.Errorf("`disk_size_gb` can't be increased above 4095GB when `caching` is set to anything other than `None`")
						}
						expandedDisk = dataDisk
					}
				}

				vm.StorageProfile.DataDisks = &dataDisks

				// fixes #2485
				vm.Identity = nil
				// fixes #1600
				vm.Resources = nil

				future, err := vmClient.CreateOrUpdate(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, vm)
				if err != nil {
					return fmt.Errorf("removing Disk %q from Virtual Machine %q : %+v", id.DiskName, virtualMachine.String(), err)
				}

				if err = future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
					return fmt.Errorf("waiting for Disk %q to be removed from Virtual Machine %q : %+v", id.DiskName, virtualMachine.String(), err)
				}
			}
		}

		// Shutdown
		if shouldShutDown {
			log.Printf("[DEBUG] Shutting Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			forceShutdown := false
			future, err := vmClient.PowerOff(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, utils.Bool(forceShutdown))
			if err != nil {
				return fmt.Errorf("sending Power Off to Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
				return fmt.Errorf("waiting for Power Off of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Shut Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}

		// De-allocate
		if shouldDeallocate {
			log.Printf("[DEBUG] Deallocating Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			// Upgrading to 2021-07-01 exposed a new hibernate paramater to the Deallocate method
			deAllocFuture, err := vmClient.Deallocate(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, utils.Bool(false))
			if err != nil {
				return fmt.Errorf("Deallocating to Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := deAllocFuture.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
				return fmt.Errorf("waiting for Deallocation of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Deallocated Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}

		// Update Disk
		err = client.UpdateThenPoll(ctx, *id, diskUpdate)
		if err != nil {
			return fmt.Errorf("updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		// Reattach DataDisk
		if shouldDetach && vm.StorageProfile != nil {
			disks := *vm.StorageProfile.DataDisks

			expandedDisk.DiskSizeGB = pointer.To(int32(pointer.From(diskUpdate.Properties.DiskSizeGB)))
			disks = append(disks, expandedDisk)

			vm.StorageProfile.DataDisks = &disks

			// fixes #2485
			vm.Identity = nil
			// fixes #1600
			vm.Resources = nil

			// if there's too many disks we get a 409 back with:
			//   `The maximum number of data disks allowed to be attached to a VM of this size is 1.`
			// which we're intentionally not wrapping, since the errors good.
			future, err := vmClient.CreateOrUpdate(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, vm)
			if err != nil {
				return fmt.Errorf("updating Virtual Machine %q to reattach Disk %q: %+v", virtualMachine.String(), name, err)
			}

			if err = future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
				return fmt.Errorf("waiting for Virtual Machine %q to finish reattaching Disk %q: %+v", virtualMachine.String(), name, err)
			}
		}

		if shouldTurnBackOn && (shouldShutDown || shouldDeallocate) {
			log.Printf("[DEBUG] Starting Linux Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			future, err := vmClient.Start(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
			if err != nil {
				return fmt.Errorf("starting Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
				return fmt.Errorf("waiting for start of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Started Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}
	} else { // otherwise, just update it
		err := client.UpdateThenPoll(ctx, *id, diskUpdate)
		if err != nil {
			return fmt.Errorf("expanding managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disks.ParseDiskID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Disk %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Managed Disk %s (resource group %s): %s", id.DiskName, id.ResourceGroupName, err)
	}

	d.Set("name", id.DiskName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		d.Set("edge_zone", flattenManagedDiskEdgeZone(model.ExtendedLocation))

		zone := ""
		if model.Zones != nil && len(*model.Zones) > 0 {
			z := *model.Zones
			zone = z[0]
		}
		d.Set("zone", zone)

		if sku := model.Sku; sku != nil {
			d.Set("storage_account_type", string(*sku.Name))
		}

		if props := model.Properties; props != nil {
			creationData := props.CreationData
			d.Set("create_option", string(creationData.CreateOption))
			if creationData.LogicalSectorSize != nil {
				d.Set("logical_sector_size", creationData.LogicalSectorSize)
			}

			// imageReference is returned as well when galleryImageRefernece is used, only check imageReference when galleryImageReference is not returned
			galleryImageReferenceId := ""
			imageReferenceId := ""
			if galleryImageReference := creationData.GalleryImageReference; galleryImageReference != nil && galleryImageReference.Id != nil {
				galleryImageReferenceId = *galleryImageReference.Id
			} else if imageReference := creationData.ImageReference; imageReference != nil && imageReference.Id != nil {
				imageReferenceId = *imageReference.Id
			}
			d.Set("gallery_image_reference_id", galleryImageReferenceId)
			d.Set("image_reference_id", imageReferenceId)

			d.Set("performance_plus_enabled", creationData.PerformancePlus)
			d.Set("source_resource_id", creationData.SourceResourceId)
			d.Set("source_uri", creationData.SourceUri)
			d.Set("storage_account_id", creationData.StorageAccountId)
			d.Set("upload_size_bytes", creationData.UploadSizeBytes)

			d.Set("disk_size_gb", props.DiskSizeGB)
			d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
			d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
			d.Set("disk_iops_read_only", props.DiskIOPSReadOnly)
			d.Set("disk_mbps_read_only", props.DiskMBpsReadOnly)
			d.Set("optimized_frequent_attach_enabled", props.OptimizedForFrequentAttach)
			d.Set("os_type", string(pointer.From(props.OsType)))
			d.Set("tier", props.Tier)
			d.Set("max_shares", props.MaxShares)
			d.Set("hyper_v_generation", string(pointer.From(props.HyperVGeneration)))

			if networkAccessPolicy := props.NetworkAccessPolicy; *networkAccessPolicy != disks.NetworkAccessPolicyAllowAll {
				d.Set("network_access_policy", string(pointer.From(props.NetworkAccessPolicy)))
			}
			d.Set("disk_access_id", props.DiskAccessId)

			d.Set("public_network_access_enabled", *props.PublicNetworkAccess == disks.PublicNetworkAccessEnabled)

			diskEncryptionSetId := ""
			if props.Encryption != nil && props.Encryption.DiskEncryptionSetId != nil {
				diskEncryptionSetId = *props.Encryption.DiskEncryptionSetId
			}
			d.Set("disk_encryption_set_id", diskEncryptionSetId)

			if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
				return fmt.Errorf("setting `encryption_settings`: %+v", err)
			}

			trustedLaunchEnabled := false
			securityType := ""
			secureVMDiskEncryptionSetId := ""
			if securityProfile := props.SecurityProfile; securityProfile != nil {
				if *securityProfile.SecurityType == disks.DiskSecurityTypesTrustedLaunch {
					trustedLaunchEnabled = true
				} else {
					securityType = string(*securityProfile.SecurityType)
				}

				if securityProfile.SecureVMDiskEncryptionSetId != nil {
					secureVMDiskEncryptionSetId = *securityProfile.SecureVMDiskEncryptionSetId
				}
			}
			d.Set("trusted_launch_enabled", trustedLaunchEnabled)
			d.Set("security_type", securityType)
			d.Set("secure_vm_disk_encryption_set_id", secureVMDiskEncryptionSetId)

			onDemandBurstingEnabled := false
			if props.BurstingEnabled != nil {
				onDemandBurstingEnabled = *props.BurstingEnabled
			}
			d.Set("on_demand_bursting_enabled", onDemandBurstingEnabled)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceManagedDiskDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disks.ParseDiskID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroupName, err)
	}

	return nil
}
