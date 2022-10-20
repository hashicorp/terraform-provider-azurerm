package migration

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = ManagedDiskV0ToV1{}

type ManagedDiskV0ToV1 struct{}

func (ManagedDiskV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId, err := disks.ParseDiskIDInsensitively(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		rawState["id"] = oldId.ID()
		return rawState, nil
	}
}

func (ManagedDiskV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
				string(disks.DiskCreateOptionRestore),
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
			ForceNew:     true,
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
			Type:             pluginsdk.TypeString,
			Optional:         true,
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
			Type:             pluginsdk.TypeString,
			Optional:         true,
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
				string(disks.DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey),
				string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey),
				string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey),
			}, false),
		},

		"hyper_v_generation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
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
	}
}

func encryptionSettingsSchema() *pluginsdk.Schema {
	if !features.FourPointOhBeta() {
		return &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:       pluginsdk.TypeBool,
						Optional:   true,
						Default:    true,
						ForceNew:   true,
						Deprecated: "Deprecated, Azure Disk Encryption is now configured directly by `disk_encryption_key` and `key_encryption_key`. To disable Azure Disk Encryption, please remove `encryption_settings` block. To enabled, specify a `encryption_settings` block`",
					},
					"disk_encryption_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"secret_url": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"source_vault_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"key_encryption_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_url": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"source_vault_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		}
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"disk_encryption_key": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_url": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},

							"source_vault_id": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
						},
					},
				},
				"key_encryption_key": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key_url": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},

							"source_vault_id": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}
