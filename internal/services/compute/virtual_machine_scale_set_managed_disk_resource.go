// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetvms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name virtual_machine_scale_set_managed_disk -service-package-name compute -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary" -test-name "empty"

type VirtualMachineScaleSetManagedDiskResource struct{}

var (
	_ sdk.ResourceWithUpdate        = VirtualMachineScaleSetManagedDiskResource{}
	_ sdk.ResourceWithIdentity      = VirtualMachineScaleSetManagedDiskResource{}
	_ sdk.ResourceWithCustomizeDiff = VirtualMachineScaleSetManagedDiskResource{}
)

type VirtualMachineScaleSetManagedDiskResourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`
	EdgeZone          string `tfschema:"edge_zone"`
	Zone              string `tfschema:"zone"`

	StorageAccountType string              `tfschema:"storage_account_type"`
	Creation           []DiskCreationModel `tfschema:"creation"`

	DiskSizeGb        int64  `tfschema:"disk_size_gb"`
	Tier              string `tfschema:"tier"`
	DiskIOPSReadWrite int64  `tfschema:"disk_iops_read_write"`
	DiskMBpsReadWrite int64  `tfschema:"disk_mbps_read_write"`
	DiskIOPSReadOnly  int64  `tfschema:"disk_iops_read_only"`
	DiskMBpsReadOnly  int64  `tfschema:"disk_mbps_read_only"`
	MaxShares         int64  `tfschema:"max_shares"`

	OnDemandBurstingEnabled        bool `tfschema:"on_demand_bursting_enabled"`
	OptimizedFrequentAttachEnabled bool `tfschema:"optimized_frequent_attach_enabled"`

	OsType           string `tfschema:"os_type"`
	HyperVGeneration string `tfschema:"hyper_v_generation"`

	TrustedLaunchEnabled        bool   `tfschema:"trusted_launch_enabled"`
	SecurityType                string `tfschema:"security_type"`
	SecureVMDiskEncryptionSetId string `tfschema:"secure_vm_disk_encryption_set_id"`

	DiskEncryptionSetId string                       `tfschema:"disk_encryption_set_id"`
	EncryptionSettings  []DiskEncryptionSettingModel `tfschema:"encryption_settings"`
	DataAccessAuthMode  string                       `tfschema:"data_access_auth_mode"`

	NetworkAccessPolicy        string `tfschema:"network_access_policy"`
	DiskAccessId               string `tfschema:"disk_access_id"`
	PublicNetworkAccessEnabled bool   `tfschema:"public_network_access_enabled"`

	Tags map[string]string `tfschema:"tags"`

	DiskSizeBytes int64  `tfschema:"disk_size_bytes"`
	UniqueId      string `tfschema:"unique_id"`
}

type DiskCreationModel struct {
	Option                  string `tfschema:"option"`
	SourceResourceId        string `tfschema:"source_resource_id"`
	SourceUri               string `tfschema:"source_uri"`
	StorageAccountId        string `tfschema:"storage_account_id"`
	GalleryImageReferenceId string `tfschema:"gallery_image_reference_id"`
	ImageReferenceId        string `tfschema:"image_reference_id"`
	UploadSizeBytes         int64  `tfschema:"upload_size_bytes"`
	LogicalSectorSize       int64  `tfschema:"logical_sector_size"`
	PerformancePlusEnabled  bool   `tfschema:"performance_plus_enabled"`
}

type DiskEncryptionSettingModel struct {
	DiskEncryptionKey []DiskEncryptionKeyModel `tfschema:"disk_encryption_key"`
	KeyEncryptionKey  []KeyEncryptionKeyModel  `tfschema:"key_encryption_key"`
}

type DiskEncryptionKeyModel struct {
	SecretUrl     string `tfschema:"secret_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type KeyEncryptionKeyModel struct {
	KeyUrl        string `tfschema:"key_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

func (r VirtualMachineScaleSetManagedDiskResource) ResourceType() string {
	return "azurerm_virtual_machine_scale_set_managed_disk"
}

func (r VirtualMachineScaleSetManagedDiskResource) ModelObject() interface{} {
	return &VirtualMachineScaleSetManagedDiskResourceModel{}
}

func (r VirtualMachineScaleSetManagedDiskResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateManagedDiskID
}

func (r VirtualMachineScaleSetManagedDiskResource) Identity() resourceids.ResourceId {
	return &commonids.ManagedDiskId{}
}

func (r VirtualMachineScaleSetManagedDiskResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedDiskName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"creation": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"option": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(disks.DiskCreateOptionEmpty),
							string(disks.DiskCreateOptionCopy),
							string(disks.DiskCreateOptionRestore),
							string(disks.DiskCreateOptionFromImage),
							string(disks.DiskCreateOptionImport),
							string(disks.DiskCreateOptionImportSecure),
							string(disks.DiskCreateOptionUpload),
						}, false),
					},

					"source_resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"source_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"gallery_image_reference_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						ValidateFunc:  validate.SharedImageVersionID,
						ConflictsWith: []string{"creation.0.image_reference_id"},
					},

					"image_reference_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"creation.0.gallery_image_reference_id"},
					},

					"upload_size_bytes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"logical_sector_size": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ForceNew: true,
						// NOTE: O+C Azure computes the logical sector size based on the storage account type when this is omitted
						Computed:     true,
						ValidateFunc: validation.IntInSlice([]int{512, 4096}),
					},

					"performance_plus_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

		"storage_account_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForDiskStorageAccountTypes(), false),
		},

		"data_access_auth_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(disks.DataAccessAuthModeNone),
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForDataAccessAuthMode(), false),
		},

		"disk_access_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
			//       https://github.com/Azure/azure-rest-api-specs/issues/14192
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     diskaccesses.ValidateDiskAccessID,
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

		"disk_iops_read_only": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C Azure assigns a default based on the disk size and SKU when this is omitted
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_iops_read_write": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C Azure assigns a default based on the disk size and SKU when this is omitted
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_mbps_read_only": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C Azure assigns a default based on the disk size and SKU when this is omitted
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_mbps_read_write": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C Azure assigns a default based on the disk size and SKU when this is omitted
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C when the disk is created from a source (`Copy`, `Restore`, `FromImage` or `Import`) the size is inherited from that source, so Azure computes this when it is omitted
			Computed:     true,
			ValidateFunc: validate.ManagedDiskSizeGB,
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"encryption_settings": encryptionSettingsSchema(),

		"hyper_v_generation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForHyperVGeneration(), false),
		},

		"max_shares": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(2, 10),
		},

		"network_access_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(disks.NetworkAccessPolicyAllowAll),
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForNetworkAccessPolicy(), false),
		},

		"on_demand_bursting_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"optimized_frequent_attach_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"os_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForOperatingSystemTypes(), false),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
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

		"tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C Azure assigns a performance tier based on the disk size when this is omitted
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"trusted_launch_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"disk_size_bytes": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			storageAccountType := rd.Get("storage_account_type").(string)
			isUltraOrPremiumV2 := storageAccountType == string(disks.DiskStorageAccountTypesUltraSSDLRS) || storageAccountType == string(disks.DiskStorageAccountTypesPremiumVTwoLRS)
			isPremium := storageAccountType == string(disks.DiskStorageAccountTypesPremiumLRS) || storageAccountType == string(disks.DiskStorageAccountTypesPremiumZRS)

			rawConfig := rd.GetRawConfig()
			isConfigured := func(attr string) bool {
				if !rawConfig.IsKnown() || rawConfig.IsNull() {
					return false
				}
				return !rawConfig.GetAttr(attr).IsNull()
			}

			if !isUltraOrPremiumV2 {
				logicalSectorSizeConfigured := false
				if rawConfig.IsKnown() && !rawConfig.IsNull() {
					if creationRaw := rawConfig.GetAttr("creation"); creationRaw.IsKnown() && !creationRaw.IsNull() && creationRaw.LengthInt() > 0 {
						if block := creationRaw.AsValueSlice()[0]; block.IsKnown() && !block.IsNull() {
							if lss := block.AsValueMap()["logical_sector_size"]; !lss.IsNull() {
								logicalSectorSizeConfigured = true
							}
						}
					}
				}
				if isConfigured("disk_iops_read_write") || isConfigured("disk_mbps_read_write") || isConfigured("disk_iops_read_only") || isConfigured("disk_mbps_read_only") || logicalSectorSizeConfigured {
					return errors.New("`disk_iops_read_write`, `disk_mbps_read_write`, `disk_iops_read_only`, `disk_mbps_read_only` and `creation.0.logical_sector_size` are only available for `UltraSSD_LRS` and `PremiumV2_LRS` disks")
				}
			}

			if (isConfigured("disk_iops_read_only") || isConfigured("disk_mbps_read_only")) && !isConfigured("max_shares") {
				return errors.New("`disk_iops_read_only` and `disk_mbps_read_only` are only available for `UltraSSD_LRS` and `PremiumV2_LRS` disks with shared disk enabled (`max_shares` set)")
			}

			if isConfigured("tier") && !isPremium {
				return errors.New("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
			}

			if rd.Get("on_demand_bursting_enabled").(bool) {
				if !isPremium {
					return errors.New("`on_demand_bursting_enabled` can only be set to `true` when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
				}
				if diskSizeGb := rd.Get("disk_size_gb").(int); diskSizeGb != 0 && diskSizeGb <= 512 {
					return errors.New("`on_demand_bursting_enabled` can only be set to `true` when `disk_size_gb` is larger than 512")
				}
			}

			if diskAccessId := rd.Get("disk_access_id").(string); diskAccessId != "" && rd.Get("network_access_policy").(string) != string(disks.NetworkAccessPolicyAllowPrivate) {
				return errors.New("`disk_access_id` is only available when `network_access_policy` is set to `AllowPrivate`")
			}

			if oldSize, newSize := rd.GetChange("disk_size_gb"); oldSize.(int) != 0 && newSize.(int) != 0 && newSize.(int) < oldSize.(int) {
				return errors.New("`disk_size_gb` can only be increased - shrinking a Managed Disk is not supported by Azure")
			}

			// Azure Disk Encryption cannot be disabled once it has been enabled, so removing the block forces a new resource
			if oldRaw, newRaw := rd.GetChange("encryption_settings"); len(oldRaw.([]interface{})) > 0 && len(newRaw.([]interface{})) == 0 {
				if err := rd.ForceNew("encryption_settings"); err != nil {
					return fmt.Errorf("setting `encryption_settings` to force a new resource: %+v", err)
				}
			}

			return nil
		},
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.DisksClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config VirtualMachineScaleSetManagedDiskResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewManagedDiskID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			creation := config.Creation[0]
			createOption := disks.DiskCreateOption(creation.Option)
			encryptionTypePlatformKey := disks.EncryptionTypeEncryptionAtRestWithPlatformKey

			props := &disks.DiskProperties{
				CreationData: disks.CreationData{
					CreateOption:    createOption,
					PerformancePlus: pointer.To(creation.PerformancePlusEnabled),
				},
				OptimizedForFrequentAttach: pointer.To(config.OptimizedFrequentAttachEnabled),
				Encryption: &disks.Encryption{
					Type: &encryptionTypePlatformKey,
				},
			}

			if config.OsType != "" {
				props.OsType = pointer.ToEnum[disks.OperatingSystemTypes](config.OsType)
			}

			if config.DiskSizeGb != 0 {
				props.DiskSizeGB = pointer.To(config.DiskSizeGb)
			}

			if config.MaxShares != 0 {
				props.MaxShares = pointer.To(config.MaxShares)
			}

			if config.DiskIOPSReadWrite != 0 {
				props.DiskIOPSReadWrite = pointer.To(config.DiskIOPSReadWrite)
			}
			if config.DiskMBpsReadWrite != 0 {
				props.DiskMBpsReadWrite = pointer.To(config.DiskMBpsReadWrite)
			}
			if config.DiskIOPSReadOnly != 0 {
				props.DiskIOPSReadOnly = pointer.To(config.DiskIOPSReadOnly)
			}
			if config.DiskMBpsReadOnly != 0 {
				props.DiskMBpsReadOnly = pointer.To(config.DiskMBpsReadOnly)
			}
			if creation.LogicalSectorSize != 0 {
				props.CreationData.LogicalSectorSize = pointer.To(creation.LogicalSectorSize)
			}

			switch createOption {
			case disks.DiskCreateOptionEmpty:
				if config.DiskSizeGb == 0 {
					return errors.New("`disk_size_gb` must be specified when `creation.0.option` is set to `Empty`")
				}
			case disks.DiskCreateOptionImport, disks.DiskCreateOptionImportSecure:
				if creation.SourceUri == "" {
					return errors.New("`creation.0.source_uri` must be specified when `creation.0.option` is set to `Import` or `ImportSecure`")
				}
				if creation.StorageAccountId == "" {
					return errors.New("`creation.0.storage_account_id` must be specified when `creation.0.option` is set to `Import` or `ImportSecure`")
				}
				props.CreationData.SourceUri = pointer.To(creation.SourceUri)
				props.CreationData.StorageAccountId = pointer.To(creation.StorageAccountId)
			case disks.DiskCreateOptionCopy, disks.DiskCreateOptionRestore:
				if creation.SourceResourceId == "" {
					return errors.New("`creation.0.source_resource_id` must be specified when `creation.0.option` is set to `Copy` or `Restore`")
				}
				props.CreationData.SourceResourceId = pointer.To(creation.SourceResourceId)
			case disks.DiskCreateOptionFromImage:
				switch {
				case creation.GalleryImageReferenceId != "":
					props.CreationData.GalleryImageReference = &disks.ImageDiskReference{
						Id: pointer.To(creation.GalleryImageReferenceId),
					}
				case creation.ImageReferenceId != "":
					props.CreationData.ImageReference = &disks.ImageDiskReference{
						Id: pointer.To(creation.ImageReferenceId),
					}
				default:
					return errors.New("`creation.0.gallery_image_reference_id` or `creation.0.image_reference_id` must be specified when `creation.0.option` is set to `FromImage`")
				}
			case disks.DiskCreateOptionUpload:
				if creation.UploadSizeBytes == 0 {
					return errors.New("`creation.0.upload_size_bytes` must be specified when `creation.0.option` is set to `Upload`")
				}
				props.CreationData.UploadSizeBytes = pointer.To(creation.UploadSizeBytes)
			}

			if len(config.EncryptionSettings) > 0 {
				props.EncryptionSettingsCollection = expandVirtualMachineScaleSetManagedDiskEncryptionSettings(config.EncryptionSettings)
			}

			if config.DiskEncryptionSetId != "" {
				encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, metadata.Client.Compute.DiskEncryptionSetsClient, config.DiskEncryptionSetId)
				if err != nil {
					return err
				}
				props.Encryption = &disks.Encryption{
					Type:                encryptionType,
					DiskEncryptionSetId: pointer.To(config.DiskEncryptionSetId),
				}
			}

			props.DataAccessAuthMode = pointer.ToEnum[disks.DataAccessAuthMode](config.DataAccessAuthMode)

			props.NetworkAccessPolicy = pointer.ToEnum[disks.NetworkAccessPolicy](config.NetworkAccessPolicy)
			if config.DiskAccessId != "" {
				props.DiskAccessId = pointer.To(config.DiskAccessId)
			}

			props.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessDisabled)
			if config.PublicNetworkAccessEnabled {
				props.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessEnabled)
			}

			if config.Tier != "" {
				props.Tier = pointer.To(config.Tier)
			}

			if config.OnDemandBurstingEnabled {
				props.BurstingEnabled = pointer.To(true)
			}

			if config.HyperVGeneration != "" {
				props.HyperVGeneration = pointer.ToEnum[disks.HyperVGeneration](config.HyperVGeneration)
			}

			if config.TrustedLaunchEnabled {
				switch createOption {
				case disks.DiskCreateOptionFromImage, disks.DiskCreateOptionImport, disks.DiskCreateOptionImportSecure:
				default:
					return fmt.Errorf("`trusted_launch_enabled` cannot be set to `true` when `creation.0.option` is set to `%s` - supported options are `FromImage`, `Import` and `ImportSecure`", createOption)
				}
				props.SecurityProfile = &disks.DiskSecurityProfile{
					SecurityType: pointer.To(disks.DiskSecurityTypesTrustedLaunch),
				}
			}

			if config.SecurityType != "" {
				if config.TrustedLaunchEnabled {
					return errors.New("`security_type` cannot be specified when `trusted_launch_enabled` is set to `true`")
				}

				switch createOption {
				case disks.DiskCreateOptionFromImage, disks.DiskCreateOptionImport, disks.DiskCreateOptionImportSecure:
				default:
					return errors.New("`security_type` can only be specified when `creation.0.option` is set to `FromImage`, `Import` or `ImportSecure`")
				}

				if disks.DiskSecurityTypes(config.SecurityType) == disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey && config.SecureVMDiskEncryptionSetId == "" {
					return errors.New("`secure_vm_disk_encryption_set_id` must be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`")
				}

				props.SecurityProfile = &disks.DiskSecurityProfile{
					SecurityType: pointer.ToEnum[disks.DiskSecurityTypes](config.SecurityType),
				}
			}

			if config.SecureVMDiskEncryptionSetId != "" {
				if disks.DiskSecurityTypes(config.SecurityType) != disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey {
					return errors.New("`secure_vm_disk_encryption_set_id` can only be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`")
				}
				props.SecurityProfile.SecureVMDiskEncryptionSetId = pointer.To(config.SecureVMDiskEncryptionSetId)
			}

			storageAccountType := disks.DiskStorageAccountTypes(config.StorageAccountType)
			payload := disks.Disk{
				Location:         location.Normalize(config.Location),
				ExtendedLocation: expandManagedDiskEdgeZone(config.EdgeZone),
				Properties:       props,
				Sku: &disks.DiskSku{
					Name: &storageAccountType,
				},
				Tags: pointer.To(config.Tags),
			}

			if config.Zone != "" {
				payload.Zones = &[]string{config.Zone}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.DisksClient

			id, err := commonids.ParseManagedDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, *id, resp.Model)
		},
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) flatten(metadata sdk.ResourceMetaData, id commonids.ManagedDiskId, model *disks.Disk) error {
	state := VirtualMachineScaleSetManagedDiskResourceModel{
		Name:              id.DiskName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.EdgeZone = flattenManagedDiskEdgeZone(model.ExtendedLocation)

		if model.Zones != nil && len(*model.Zones) > 0 {
			state.Zone = (*model.Zones)[0]
		}

		if sku := model.Sku; sku != nil {
			state.StorageAccountType = string(pointer.From(sku.Name))
		}

		if props := model.Properties; props != nil {
			state.Creation = r.flattenCreation(props.CreationData)
			state.DiskSizeGb = pointer.From(props.DiskSizeGB)
			state.DiskSizeBytes = pointer.From(props.DiskSizeBytes)
			state.Tier = pointer.From(props.Tier)
			state.DiskIOPSReadWrite = pointer.From(props.DiskIOPSReadWrite)
			state.DiskMBpsReadWrite = pointer.From(props.DiskMBpsReadWrite)
			state.DiskIOPSReadOnly = pointer.From(props.DiskIOPSReadOnly)
			state.DiskMBpsReadOnly = pointer.From(props.DiskMBpsReadOnly)
			if v := pointer.From(props.MaxShares); v > 1 {
				state.MaxShares = v
			}
			state.OnDemandBurstingEnabled = pointer.From(props.BurstingEnabled)
			state.OptimizedFrequentAttachEnabled = pointer.From(props.OptimizedForFrequentAttach)
			state.OsType = string(pointer.From(props.OsType))
			state.HyperVGeneration = string(pointer.From(props.HyperVGeneration))
			state.DataAccessAuthMode = string(disks.DataAccessAuthModeNone)
			if v := pointer.From(props.DataAccessAuthMode); v != "" {
				state.DataAccessAuthMode = string(v)
			}
			state.NetworkAccessPolicy = string(pointer.From(props.NetworkAccessPolicy))
			state.DiskAccessId = pointer.From(props.DiskAccessId)
			state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == disks.PublicNetworkAccessEnabled
			state.UniqueId = pointer.From(props.UniqueId)

			if props.Encryption != nil {
				state.DiskEncryptionSetId = pointer.From(props.Encryption.DiskEncryptionSetId)
			}

			state.EncryptionSettings = flattenVirtualMachineScaleSetManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)

			if securityProfile := props.SecurityProfile; securityProfile != nil {
				if pointer.From(securityProfile.SecurityType) == disks.DiskSecurityTypesTrustedLaunch {
					state.TrustedLaunchEnabled = true
				} else {
					state.SecurityType = string(pointer.From(securityProfile.SecurityType))
				}
				state.SecureVMDiskEncryptionSetId = pointer.From(securityProfile.SecureVMDiskEncryptionSetId)
			}
		}

		state.Tags = pointer.From(model.Tags)
	}

	return metadata.Encode(&state)
}

func (r VirtualMachineScaleSetManagedDiskResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.DisksClient

			id, err := commonids.ParseManagedDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.DisksClient

			id, err := commonids.ParseManagedDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config VirtualMachineScaleSetManagedDiskResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			requiresDetach := r.updateRequiresDetach(metadata, existing.Model, config.DiskSizeGb)

			payload := *existing.Model
			props := *existing.Model.Properties
			payload.Properties = &props

			if metadata.ResourceData.HasChange("storage_account_type") {
				storageAccountType := disks.DiskStorageAccountTypes(config.StorageAccountType)
				payload.Sku = &disks.DiskSku{
					Name: &storageAccountType,
				}
			}

			if metadata.ResourceData.HasChange("disk_size_gb") {
				props.DiskSizeGB = pointer.To(config.DiskSizeGb)
			}

			if metadata.ResourceData.HasChange("tier") {
				props.Tier = pointer.To(config.Tier)
			}

			if metadata.ResourceData.HasChange("disk_iops_read_write") {
				props.DiskIOPSReadWrite = pointer.To(config.DiskIOPSReadWrite)
			}

			if metadata.ResourceData.HasChange("disk_mbps_read_write") {
				props.DiskMBpsReadWrite = pointer.To(config.DiskMBpsReadWrite)
			}

			if metadata.ResourceData.HasChange("disk_iops_read_only") {
				props.DiskIOPSReadOnly = pointer.To(config.DiskIOPSReadOnly)
			}

			if metadata.ResourceData.HasChange("disk_mbps_read_only") {
				props.DiskMBpsReadOnly = pointer.To(config.DiskMBpsReadOnly)
			}

			if metadata.ResourceData.HasChange("max_shares") {
				props.MaxShares = pointer.To(config.MaxShares)
			}

			if metadata.ResourceData.HasChange("on_demand_bursting_enabled") {
				props.BurstingEnabled = pointer.To(config.OnDemandBurstingEnabled)
			}

			if metadata.ResourceData.HasChange("optimized_frequent_attach_enabled") {
				props.OptimizedForFrequentAttach = pointer.To(config.OptimizedFrequentAttachEnabled)
			}

			if metadata.ResourceData.HasChange("os_type") {
				props.OsType = nil
				if config.OsType != "" {
					props.OsType = pointer.ToEnum[disks.OperatingSystemTypes](config.OsType)
				}
			}

			if metadata.ResourceData.HasChange("encryption_settings") {
				props.EncryptionSettingsCollection = expandVirtualMachineScaleSetManagedDiskEncryptionSettings(config.EncryptionSettings)
			}

			if metadata.ResourceData.HasChange("data_access_auth_mode") {
				props.DataAccessAuthMode = pointer.ToEnum[disks.DataAccessAuthMode](config.DataAccessAuthMode)
			}

			if metadata.ResourceData.HasChange("network_access_policy") {
				props.NetworkAccessPolicy = pointer.ToEnum[disks.NetworkAccessPolicy](config.NetworkAccessPolicy)
			}

			if metadata.ResourceData.HasChange("disk_access_id") {
				props.DiskAccessId = nil
				if config.DiskAccessId != "" {
					props.DiskAccessId = pointer.To(config.DiskAccessId)
				}
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				props.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessDisabled)
				if config.PublicNetworkAccessEnabled {
					props.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessEnabled)
				}
			}

			if metadata.ResourceData.HasChange("disk_encryption_set_id") {
				if config.DiskEncryptionSetId == "" {
					return errors.New("once a customer-managed key is used, you cannot change the selection back to a platform-managed key")
				}
				encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, metadata.Client.Compute.DiskEncryptionSetsClient, config.DiskEncryptionSetId)
				if err != nil {
					return err
				}
				props.Encryption = &disks.Encryption{
					Type:                encryptionType,
					DiskEncryptionSetId: pointer.To(config.DiskEncryptionSetId),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if !requiresDetach {
				if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
				return nil
			}

			return r.updateWithDetach(ctx, metadata, *id, existing.Model.ManagedBy, payload)
		},
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) updateWithDetach(ctx context.Context, metadata sdk.ResourceMetaData, id commonids.ManagedDiskId, managedBy *string, payload disks.Disk) error {
	client := metadata.Client.Compute.DisksClient

	if managedBy == nil || *managedBy == "" {
		if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
		return nil
	}

	instanceId, err := virtualmachinescalesetvms.ParseVirtualMachineScaleSetVirtualMachineIDInsensitively(*managedBy)
	if err != nil {
		if _, vmErr := virtualmachines.ParseVirtualMachineIDInsensitively(*managedBy); vmErr == nil {
			return fmt.Errorf("%s is attached to the standalone Virtual Machine %q - this change requires the disk to be taken offline, which is only supported for disks attached to a Virtual Machine Scale Set instance by this resource; please manage this disk with `azurerm_managed_disk` instead", id, *managedBy)
		}
		if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
		return nil
	}

	vmClient := metadata.Client.Compute.VirtualMachineScaleSetVMsClient

	locks.ByName(instanceId.VirtualMachineScaleSetName, VirtualMachineScaleSetResourceName)
	defer locks.UnlockByName(instanceId.VirtualMachineScaleSetName, VirtualMachineScaleSetResourceName)

	instance, err := vmClient.Get(ctx, *instanceId, virtualmachinescalesetvms.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *instanceId, err)
	}

	var lun *int64
	var caching *virtualmachinescalesetvms.CachingTypes
	var writeAcceleratorEnabled *bool
	if model := instance.Model; model != nil && model.Properties != nil && model.Properties.StorageProfile != nil && model.Properties.StorageProfile.DataDisks != nil {
		for _, dataDisk := range *model.Properties.StorageProfile.DataDisks {
			if dataDisk.ManagedDisk == nil || dataDisk.ManagedDisk.Id == nil {
				continue
			}
			dataDiskId, err := commonids.ParseManagedDiskIDInsensitively(*dataDisk.ManagedDisk.Id)
			if err != nil {
				continue
			}
			if strings.EqualFold(dataDiskId.ID(), id.ID()) {
				lun = pointer.To(dataDisk.Lun)
				caching = dataDisk.Caching
				writeAcceleratorEnabled = dataDisk.WriteAcceleratorEnabled
				break
			}
		}
	}

	if lun == nil {
		return fmt.Errorf("locating %s as a data disk on %s", id, *instanceId)
	}

	detach := virtualmachinescalesetvms.AttachDetachDataDisksRequest{
		DataDisksToDetach: &[]virtualmachinescalesetvms.DataDisksToDetach{
			{
				DiskId: id.ID(),
			},
		},
	}
	if err := vmClient.AttachDetachDataDisksThenPoll(ctx, *instanceId, detach); err != nil {
		return fmt.Errorf("detaching %s from %s: %+v", id, *instanceId, err)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	attach := virtualmachinescalesetvms.AttachDetachDataDisksRequest{
		DataDisksToAttach: &[]virtualmachinescalesetvms.DataDisksToAttach{
			{
				DiskId:                  id.ID(),
				Lun:                     lun,
				Caching:                 caching,
				WriteAcceleratorEnabled: writeAcceleratorEnabled,
			},
		},
	}
	if err := vmClient.AttachDetachDataDisksThenPoll(ctx, *instanceId, attach); err != nil {
		return fmt.Errorf("re-attaching %s to %s: %+v", id, *instanceId, err)
	}

	return nil
}

func (r VirtualMachineScaleSetManagedDiskResource) updateRequiresDetach(metadata sdk.ResourceMetaData, existing *disks.Disk, newSizeGb int64) bool {
	if metadata.ResourceData.HasChange("storage_account_type") {
		return true
	}

	if metadata.ResourceData.HasChange("tier") {
		return true
	}

	if metadata.ResourceData.HasChange("disk_encryption_set_id") {
		return true
	}

	if metadata.ResourceData.HasChange("on_demand_bursting_enabled") {
		return true
	}

	if metadata.ResourceData.HasChange("max_shares") {
		return true
	}

	if metadata.ResourceData.HasChange("disk_size_gb") && existing != nil && existing.Sku != nil && existing.Properties != nil {
		skuName := string(pointer.From(existing.Sku.Name))
		supportsOnlineExpandAbove4TiB := strings.EqualFold(skuName, string(disks.DiskStorageAccountTypesPremiumVTwoLRS)) || strings.EqualFold(skuName, string(disks.DiskStorageAccountTypesUltraSSDLRS))
		if !supportsOnlineExpandAbove4TiB {
			oldSizeGb := pointer.From(existing.Properties.DiskSizeGB)
			if oldSizeGb < 4096 && newSizeGb >= 4096 {
				return true
			}
		}
	}

	return false
}

func expandVirtualMachineScaleSetManagedDiskEncryptionSettings(input []DiskEncryptionSettingModel) *disks.EncryptionSettingsCollection {
	if len(input) == 0 {
		return &disks.EncryptionSettingsCollection{
			Enabled: false,
		}
	}

	setting := input[0]
	config := &disks.EncryptionSettingsCollection{
		Enabled: true,
	}

	var diskEncryptionKey *disks.KeyVaultAndSecretReference
	if len(setting.DiskEncryptionKey) > 0 {
		dek := setting.DiskEncryptionKey[0]
		diskEncryptionKey = &disks.KeyVaultAndSecretReference{
			SecretURL: dek.SecretUrl,
			SourceVault: disks.SourceVault{
				Id: pointer.To(dek.SourceVaultId),
			},
		}
	}

	var keyEncryptionKey *disks.KeyVaultAndKeyReference
	if len(setting.KeyEncryptionKey) > 0 {
		kek := setting.KeyEncryptionKey[0]
		keyEncryptionKey = &disks.KeyVaultAndKeyReference{
			KeyURL: kek.KeyUrl,
			SourceVault: disks.SourceVault{
				Id: pointer.To(kek.SourceVaultId),
			},
		}
	}

	config.EncryptionSettings = &[]disks.EncryptionSettingsElement{
		{
			DiskEncryptionKey: diskEncryptionKey,
			KeyEncryptionKey:  keyEncryptionKey,
		},
	}

	return config
}

func flattenVirtualMachineScaleSetManagedDiskEncryptionSettings(input *disks.EncryptionSettingsCollection) []DiskEncryptionSettingModel {
	if input == nil || input.EncryptionSettings == nil || len(*input.EncryptionSettings) == 0 {
		return []DiskEncryptionSettingModel{}
	}

	setting := (*input.EncryptionSettings)[0]

	diskEncryptionKeys := make([]DiskEncryptionKeyModel, 0)
	if key := setting.DiskEncryptionKey; key != nil {
		diskEncryptionKeys = append(diskEncryptionKeys, DiskEncryptionKeyModel{
			SecretUrl:     key.SecretURL,
			SourceVaultId: pointer.From(key.SourceVault.Id),
		})
	}

	keyEncryptionKeys := make([]KeyEncryptionKeyModel, 0)
	if key := setting.KeyEncryptionKey; key != nil {
		keyEncryptionKeys = append(keyEncryptionKeys, KeyEncryptionKeyModel{
			KeyUrl:        key.KeyURL,
			SourceVaultId: pointer.From(key.SourceVault.Id),
		})
	}

	if len(diskEncryptionKeys) == 0 {
		return []DiskEncryptionSettingModel{}
	}

	return []DiskEncryptionSettingModel{
		{
			DiskEncryptionKey: diskEncryptionKeys,
			KeyEncryptionKey:  keyEncryptionKeys,
		},
	}
}

func (r VirtualMachineScaleSetManagedDiskResource) flattenCreation(input disks.CreationData) []DiskCreationModel {
	creation := DiskCreationModel{
		Option:                 string(input.CreateOption),
		PerformancePlusEnabled: pointer.From(input.PerformancePlus),
		SourceResourceId:       pointer.From(input.SourceResourceId),
		SourceUri:              pointer.From(input.SourceUri),
		StorageAccountId:       pointer.From(input.StorageAccountId),
		UploadSizeBytes:        pointer.From(input.UploadSizeBytes),
	}

	creation.LogicalSectorSize = pointer.From(input.LogicalSectorSize)

	if input.GalleryImageReference != nil {
		creation.GalleryImageReferenceId = pointer.From(input.GalleryImageReference.Id)
	} else if input.ImageReference != nil {
		creation.ImageReferenceId = pointer.From(input.ImageReference.Id)
	}

	return []DiskCreationModel{creation}
}
