package compute

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceManagedDisk() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedDiskCreate,
		Read:   resourceManagedDiskRead,
		Update: resourceManagedDiskUpdate,
		Delete: resourceManagedDiskDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedDiskID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.StorageAccountTypesStandardLRS),
					string(compute.StorageAccountTypesStandardSSDZRS),
					string(compute.StorageAccountTypesPremiumLRS),
					string(compute.StorageAccountTypesPremiumZRS),
					string(compute.StorageAccountTypesStandardSSDLRS),
					string(compute.StorageAccountTypesUltraSSDLRS),
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"create_option": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.DiskCreateOptionCopy),
					string(compute.DiskCreateOptionEmpty),
					string(compute.DiskCreateOptionFromImage),
					string(compute.DiskCreateOptionImport),
					string(compute.DiskCreateOptionRestore),
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
				ValidateFunc: azure.ValidateResourceID,
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
					string(compute.OperatingSystemTypesWindows),
					string(compute.OperatingSystemTypesLinux),
				}, false),
			},

			"disk_size_gb": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ManagedDiskSizeGB,
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
					string(compute.NetworkAccessPolicyAllowAll),
					string(compute.NetworkAccessPolicyAllowPrivate),
					string(compute.NetworkAccessPolicyDenyAll),
				}, false),
			},
			"disk_access_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
				//       https://github.com/Azure/azure-rest-api-specs/issues/14192
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
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
					string(compute.DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey),
					string(compute.DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey),
					string(compute.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey),
				}, false),
			},

			"hyper_v_generation": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true, // Not supported by disk update
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.HyperVGenerationV1),
					string(compute.HyperVGenerationV2),
				}, false),
			},

			"on_demand_bursting_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"zone": commonschema.ZoneSingleOptionalForceNew(),

			"tags": tags.Schema(),
		},
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

	id := parse.NewManagedDiskID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.DiskName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Managed Disk %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_managed_disk", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := compute.DiskCreateOption(d.Get("create_option").(string))
	storageAccountType := d.Get("storage_account_type").(string)
	osType := d.Get("os_type").(string)
	maxShares := d.Get("max_shares").(int)

	t := d.Get("tags").(map[string]interface{})
	skuName := compute.DiskStorageAccountTypes(storageAccountType)

	props := &compute.DiskProperties{
		CreationData: &compute.CreationData{
			CreateOption: createOption,
		},
		OsType: compute.OperatingSystemTypes(osType),
		Encryption: &compute.Encryption{
			Type: compute.EncryptionTypeEncryptionAtRestWithPlatformKey,
		},
	}

	diskSizeGB := d.Get("disk_size_gb").(int)
	if diskSizeGB != 0 {
		props.DiskSizeGB = utils.Int32(int32(diskSizeGB))
	}

	if maxShares != 0 {
		props.MaxShares = utils.Int32(int32(maxShares))
	}

	if storageAccountType == string(compute.StorageAccountTypesUltraSSDLRS) {
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
				return fmt.Errorf("[ERROR] disk_iops_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			props.DiskIOPSReadOnly = utils.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("disk_mbps_read_only"); ok {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_mbps_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			props.DiskMBpsReadOnly = utils.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("logical_sector_size"); ok {
			props.CreationData.LogicalSectorSize = utils.Int32(int32(v.(int)))
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") || d.HasChange("disk_iops_read_only") || d.HasChange("disk_mbps_read_only") || d.HasChange("logical_sector_size") {
		return fmt.Errorf("[ERROR] disk_iops_read_write, disk_mbps_read_write, disk_iops_read_only, disk_mbps_read_only and logical_sector_size are only available for UltraSSD disks")
	}

	if createOption == compute.DiskCreateOptionImport {
		sourceUri := d.Get("source_uri").(string)
		if sourceUri == "" {
			return fmt.Errorf("`source_uri` must be specified when `create_option` is set to `Import`")
		}

		storageAccountId := d.Get("storage_account_id").(string)
		if storageAccountId == "" {
			return fmt.Errorf("`storage_account_id` must be specified when `create_option` is set to `Import`")
		}

		props.CreationData.StorageAccountID = utils.String(storageAccountId)
		props.CreationData.SourceURI = utils.String(sourceUri)
	}
	if createOption == compute.DiskCreateOptionCopy || createOption == compute.DiskCreateOptionRestore {
		sourceResourceId := d.Get("source_resource_id").(string)
		if sourceResourceId == "" {
			return fmt.Errorf("`source_resource_id` must be specified when `create_option` is set to `Copy` or `Restore`")
		}

		props.CreationData.SourceResourceID = utils.String(sourceResourceId)
	}
	if createOption == compute.DiskCreateOptionFromImage {
		if imageReferenceId := d.Get("image_reference_id").(string); imageReferenceId != "" {
			props.CreationData.ImageReference = &compute.ImageDiskReference{
				ID: utils.String(imageReferenceId),
			}
		} else if galleryImageReferenceId := d.Get("gallery_image_reference_id").(string); galleryImageReferenceId != "" {
			props.CreationData.GalleryImageReference = &compute.ImageDiskReference{
				ID: utils.String(galleryImageReferenceId),
			}
		} else {
			return fmt.Errorf("`image_reference_id` or `gallery_image_reference_id` must be specified when `create_option` is set to `FromImage`")
		}
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		encryptionSettings := v.([]interface{})
		settings := encryptionSettings[0].(map[string]interface{})
		props.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(settings)
	}

	if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
		encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
		if err != nil {
			return err
		}

		props.Encryption = &compute.Encryption{
			Type:                *encryptionType,
			DiskEncryptionSetID: utils.String(diskEncryptionSetId),
		}
	}

	if networkAccessPolicy := d.Get("network_access_policy").(string); networkAccessPolicy != "" {
		props.NetworkAccessPolicy = compute.NetworkAccessPolicy(networkAccessPolicy)
	} else {
		props.NetworkAccessPolicy = compute.NetworkAccessPolicyAllowAll
	}

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		switch {
		case props.NetworkAccessPolicy == compute.NetworkAccessPolicyAllowPrivate:
			props.DiskAccessID = utils.String(diskAccessID)
		case diskAccessID != "" && props.NetworkAccessPolicy != compute.NetworkAccessPolicyAllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			props.DiskAccessID = nil
		}
	}

	if d.Get("public_network_access_enabled").(bool) {
		props.PublicNetworkAccess = compute.PublicNetworkAccessEnabled
	} else {
		props.PublicNetworkAccess = compute.PublicNetworkAccessDisabled
	}

	if tier := d.Get("tier").(string); tier != "" {
		if storageAccountType != string(compute.StorageAccountTypesPremiumZRS) && storageAccountType != string(compute.StorageAccountTypesPremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		props.Tier = &tier
	}

	if d.Get("trusted_launch_enabled").(bool) {
		props.SecurityProfile = &compute.DiskSecurityProfile{
			SecurityType: compute.DiskSecurityTypesTrustedLaunch,
		}

		switch createOption {
		case compute.DiskCreateOptionFromImage:
		case compute.DiskCreateOptionImport:
		default:
			return fmt.Errorf("trusted_launch_enabled cannot be set to true with create_option %q. Supported Create Options when Trusted Launch is enabled are FromImage, Import", createOption)
		}
	}

	securityType := d.Get("security_type").(string)
	secureVMDiskEncryptionId := d.Get("secure_vm_disk_encryption_set_id")
	if securityType != "" {
		if d.Get("trusted_launch_enabled").(bool) {
			return fmt.Errorf("`security_type` cannot be specified when `trusted_launch_enabled` is set to `true`")
		}

		switch createOption {
		case compute.DiskCreateOptionFromImage:
		case compute.DiskCreateOptionImport:
		default:
			return fmt.Errorf("`security_type` can only be specified when `create_option` is set to `FromImage` or `Import`")
		}

		if compute.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey == compute.DiskSecurityTypes(securityType) && secureVMDiskEncryptionId == "" {
			return fmt.Errorf("`secure_vm_disk_encryption_set_id` must be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`")
		}

		props.SecurityProfile = &compute.DiskSecurityProfile{
			SecurityType: compute.DiskSecurityTypes(securityType),
		}
	}

	if secureVMDiskEncryptionId != "" {
		if compute.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey != compute.DiskSecurityTypes(securityType) {
			return fmt.Errorf("`secure_vm_disk_encryption_set_id` can only be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`")
		}
		props.SecurityProfile.SecureVMDiskEncryptionSetID = utils.String(secureVMDiskEncryptionId.(string))
	}

	if d.Get("on_demand_bursting_enabled").(bool) {
		switch storageAccountType {
		case string(compute.StorageAccountTypesPremiumLRS):
		case string(compute.StorageAccountTypesPremiumZRS):
		default:
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}

		if diskSizeGB != 0 && diskSizeGB <= 512 {
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `disk_size_gb` is larger than 512GB")
		}

		props.BurstingEnabled = utils.Bool(true)
	}

	if v, ok := d.GetOk("hyper_v_generation"); ok {
		props.HyperVGeneration = compute.HyperVGeneration(v.(string))
	}

	createDisk := compute.Disk{
		Name:             &name,
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         &location,
		DiskProperties:   props,
		Sku: &compute.DiskSku{
			Name: skuName,
		},
		Tags: tags.Expand(t),
	}

	if zone, ok := d.GetOk("zone"); ok {
		createDisk.Zones = &[]string{
			zone.(string),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, createDisk)
	if err != nil {
		return fmt.Errorf("creating/updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("reading Managed Disk %s (Resource Group %q): ID was nil", name, resourceGroup)
	}

	d.SetId(id.ID())

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
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

	disk, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(disk.Response) {
			return fmt.Errorf("Managed Disk %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("making Read request on Azure Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	diskUpdate := compute.DiskUpdate{
		DiskUpdateProperties: &compute.DiskUpdateProperties{},
	}

	if d.HasChange("max_shares") {
		diskUpdate.MaxShares = utils.Int32(int32(maxShares))
		var skuName compute.DiskStorageAccountTypes
		for _, v := range compute.PossibleDiskStorageAccountTypesValues() {
			if strings.EqualFold(storageAccountType, string(v)) {
				skuName = v
			}
		}
		diskUpdate.Sku = &compute.DiskSku{
			Name: skuName,
		}
	}

	if d.HasChange("tier") {
		if storageAccountType != string(compute.StorageAccountTypesPremiumZRS) && storageAccountType != string(compute.StorageAccountTypesPremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		shouldShutDown = true
		tier := d.Get("tier").(string)
		diskUpdate.Tier = &tier
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		diskUpdate.Tags = tags.Expand(t)
	}

	if d.HasChange("storage_account_type") {
		shouldShutDown = true
		var skuName compute.DiskStorageAccountTypes
		for _, v := range compute.PossibleDiskStorageAccountTypesValues() {
			if strings.EqualFold(storageAccountType, string(v)) {
				skuName = v
			}
		}
		diskUpdate.Sku = &compute.DiskSku{
			Name: skuName,
		}
	}

	if strings.EqualFold(storageAccountType, string(compute.StorageAccountTypesUltraSSDLRS)) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			diskUpdate.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int64(v.(int))
			diskUpdate.DiskMBpsReadWrite = &diskMBps
		}

		if d.HasChange("disk_iops_read_only") {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_iops_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			v := d.Get("disk_iops_read_only")
			diskUpdate.DiskIOPSReadOnly = utils.Int64(int64(v.(int)))
		}

		if d.HasChange("disk_mbps_read_only") {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_mbps_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			v := d.Get("disk_mbps_read_only")
			diskUpdate.DiskMBpsReadOnly = utils.Int64(int64(v.(int)))
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") || d.HasChange("disk_iops_read_only") || d.HasChange("disk_mbps_read_only") {
		return fmt.Errorf("[ERROR] disk_iops_read_write, disk_mbps_read_write, disk_iops_read_only and disk_mbps_read_only are only available for UltraSSD disks")
	}

	if d.HasChange("os_type") {
		diskUpdate.DiskUpdateProperties.OsType = compute.OperatingSystemTypes(d.Get("os_type").(string))
	}

	if d.HasChange("disk_size_gb") {
		if old, new := d.GetChange("disk_size_gb"); new.(int) > old.(int) {
			shouldShutDown = true
			diskUpdate.DiskUpdateProperties.DiskSizeGB = utils.Int32(int32(new.(int)))
		} else {
			return fmt.Errorf("- New size must be greater than original size. Shrinking disks is not supported on Azure")
		}
	}

	if d.HasChange("disk_encryption_set_id") {
		shouldShutDown = true
		if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
			if err != nil {
				return err
			}

			diskUpdate.Encryption = &compute.Encryption{
				Type:                *encryptionType,
				DiskEncryptionSetID: utils.String(diskEncryptionSetId),
			}
		} else {
			return fmt.Errorf("Once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if networkAccessPolicy := d.Get("network_access_policy").(string); networkAccessPolicy != "" {
		diskUpdate.NetworkAccessPolicy = compute.NetworkAccessPolicy(networkAccessPolicy)
	} else {
		diskUpdate.NetworkAccessPolicy = compute.NetworkAccessPolicyAllowAll
	}

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		switch {
		case diskUpdate.NetworkAccessPolicy == compute.NetworkAccessPolicyAllowPrivate:
			diskUpdate.DiskAccessID = utils.String(diskAccessID)
		case diskAccessID != "" && diskUpdate.NetworkAccessPolicy != compute.NetworkAccessPolicyAllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			diskUpdate.DiskAccessID = nil
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if d.Get("public_network_access_enabled").(bool) {
			diskUpdate.PublicNetworkAccess = compute.PublicNetworkAccessEnabled
		} else {
			diskUpdate.PublicNetworkAccess = compute.PublicNetworkAccessDisabled
		}
	}

	if onDemandBurstingEnabled {
		switch storageAccountType {
		case string(compute.StorageAccountTypesPremiumLRS):
		case string(compute.StorageAccountTypesPremiumZRS):
		default:
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}

		if diskSizeGB != 0 && diskSizeGB <= 512 {
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `disk_size_gb` is larger than 512GB")
		}
	}

	if d.HasChange("on_demand_bursting_enabled") {
		shouldShutDown = true
		diskUpdate.BurstingEnabled = utils.Bool(onDemandBurstingEnabled)
	}

	// whilst we need to shut this down, if we're not attached to anything there's no point
	if shouldShutDown && disk.ManagedBy == nil {
		shouldShutDown = false
	}

	// if we are attached to a VM we bring down the VM as necessary for the operations which are not allowed while it's online
	if shouldShutDown {
		virtualMachine, err := parse.VirtualMachineID(*disk.ManagedBy)
		if err != nil {
			return fmt.Errorf("parsing VMID %q for disk attachment: %+v", *disk.ManagedBy, err)
		}
		// check instanceView State
		vmClient := meta.(*clients.Client).Compute.VMClient

		locks.ByName(name, VirtualMachineResourceName)
		defer locks.UnlockByName(name, VirtualMachineResourceName)

		instanceView, err := vmClient.InstanceView(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
		if err != nil {
			return fmt.Errorf("retrieving InstanceView for Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		shouldTurnBackOn := true
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
				case "deallocating":
					shouldTurnBackOn = false
					shouldShutDown = false
					shouldDeallocate = false
				case "stopping":
				case "stopped":
					shouldShutDown = false
					shouldTurnBackOn = false
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

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

			if err := deAllocFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Deallocation of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Deallocated Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}

		// Update Disk
		updateFuture, err := client.Update(ctx, resourceGroup, name, diskUpdate)
		if err != nil {
			return fmt.Errorf("updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if shouldTurnBackOn {
			log.Printf("[DEBUG] Starting Linux Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			future, err := vmClient.Start(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
			if err != nil {
				return fmt.Errorf("starting Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for start of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Started Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}
	} else { // otherwise, just update it
		diskFuture, err := client.Update(ctx, resourceGroup, name, diskUpdate)
		if err != nil {
			return fmt.Errorf("expanding managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		err = diskFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("waiting for expand operation on managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DiskName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Disk %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Managed Disk %s (resource group %s): %s", id.DiskName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("edge_zone", flattenEdgeZone(resp.ExtendedLocation))

	zone := ""
	if resp.Zones != nil && len(*resp.Zones) > 0 {
		z := *resp.Zones
		zone = z[0]
	}
	d.Set("zone", zone)

	if sku := resp.Sku; sku != nil {
		d.Set("storage_account_type", string(sku.Name))
	}

	if props := resp.DiskProperties; props != nil {
		if creationData := props.CreationData; creationData != nil {
			d.Set("create_option", string(creationData.CreateOption))
			if creationData.LogicalSectorSize != nil {
				d.Set("logical_sector_size", creationData.LogicalSectorSize)
			}

			// imageReference is returned as well when galleryImageRefernece is used, only check imageReference when galleryImageReference is not returned
			galleryImageReferenceId := ""
			imageReferenceId := ""
			if galleryImageReference := creationData.GalleryImageReference; galleryImageReference != nil && galleryImageReference.ID != nil {
				galleryImageReferenceId = *galleryImageReference.ID
			} else if imageReference := creationData.ImageReference; imageReference != nil && imageReference.ID != nil {
				imageReferenceId = *imageReference.ID
			}
			d.Set("gallery_image_reference_id", galleryImageReferenceId)
			d.Set("image_reference_id", imageReferenceId)

			d.Set("source_resource_id", creationData.SourceResourceID)
			d.Set("source_uri", creationData.SourceURI)
			d.Set("storage_account_id", creationData.StorageAccountID)
		}

		d.Set("disk_size_gb", props.DiskSizeGB)
		d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
		d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
		d.Set("disk_iops_read_only", props.DiskIOPSReadOnly)
		d.Set("disk_mbps_read_only", props.DiskMBpsReadOnly)
		d.Set("os_type", props.OsType)
		d.Set("tier", props.Tier)
		d.Set("max_shares", props.MaxShares)
		d.Set("hyper_v_generation", props.HyperVGeneration)

		if networkAccessPolicy := props.NetworkAccessPolicy; networkAccessPolicy != compute.NetworkAccessPolicyAllowAll {
			d.Set("network_access_policy", props.NetworkAccessPolicy)
		}
		d.Set("disk_access_id", props.DiskAccessID)

		d.Set("public_network_access_enabled", props.PublicNetworkAccess == compute.PublicNetworkAccessEnabled)

		diskEncryptionSetId := ""
		if props.Encryption != nil && props.Encryption.DiskEncryptionSetID != nil {
			diskEncryptionSetId = *props.Encryption.DiskEncryptionSetID
		}
		d.Set("disk_encryption_set_id", diskEncryptionSetId)

		if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
			return fmt.Errorf("setting `encryption_settings`: %+v", err)
		}

		trustedLaunchEnabled := false
		securityType := ""
		secureVMDiskEncryptionSetId := ""
		if securityProfile := props.SecurityProfile; securityProfile != nil {
			if securityProfile.SecurityType == compute.DiskSecurityTypesTrustedLaunch {
				trustedLaunchEnabled = true
			} else {
				securityType = string(securityProfile.SecurityType)
			}

			if securityProfile.SecureVMDiskEncryptionSetID != nil {
				secureVMDiskEncryptionSetId = *securityProfile.SecureVMDiskEncryptionSetID
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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceManagedDiskDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.DiskName)
	if err != nil {
		return fmt.Errorf("deleting Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroup, err)
	}

	return nil
}
