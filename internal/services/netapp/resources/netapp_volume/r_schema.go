// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func netappVolumeSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VolumeName,
		},

		"location": commonschema.Location(),

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"pool_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.PoolName,
		},

		"volume_path": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VolumePath,
		},

		"service_level": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumes.ServiceLevelPremium),
				string(volumes.ServiceLevelStandard),
				string(volumes.ServiceLevelUltra),
				string(volumes.ServiceLevelFlexible),
			}, false),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"create_from_snapshot_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: snapshots.ValidateSnapshotID,
		},

		"accept_grow_capacity_pool_for_short_term_clone_split": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForAcceptGrowCapacityPoolForShortTermCloneSplit(), false),
			Description:  "While auto splitting the short term clone volume, if the parent pool does not have enough space to accommodate the volume after split, it will be automatically resized, which will lead to increased billing. To accept capacity pool size auto grow and create a short term clone volume, set the property as accepted. Can only be used in conjunction with `create_from_snapshot_resource_id`.",
		},

		"network_features": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because the service team is changing network features on the backend to upgrade everyone from Basic to Standard and there is a feature that allows customers to change network features from portal but not the API. This could cause drift that forces data loss that we want to avoid
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumes.NetworkFeaturesBasic),
				string(volumes.NetworkFeaturesStandard),
			}, false),
		},

		"protocols": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 2,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"NFSv3",
					"NFSv4.1",
					"CIFS",
				}, false),
			},
		},

		"kerberos_enabled": {
			// Due to large infrastructure requirements, there is not a reliable way to test Kerberos volumes in a shared environment at this time
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Description: "Enable to allow Kerberos secured volumes. Requires appropriate export rules as well as the parent `azurerm_netapp_account` having a defined AD connection.",
		},

		"smb_continuous_availability_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Description: "Continuous availability option should be used only for SQL and FSLogix workloads. Using it for any other SMB workloads is not supported.",
			ForceNew:    true,
		},

		"smb3_protocol_encryption_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Description: "SMB3 encryption option should be used only for SMB/DualProtocol volumes. Using it for any other workloads is not supported.",
			ForceNew:    true,
		},

		"security_style": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForSecurityStyle(), false),
		},

		"storage_quota_in_gb": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			// The upper bound covers the largest size any volume flavour supports (2,400 TiB, only reachable
			// with `breakthrough_mode_enabled`), the per-flavour limits are enforced in CustomizeDiff.
			ValidateFunc: validation.IntBetween(50, 2457600),
		},

		"throughput_in_mibps": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.FloatAtLeast(1.0),
		},

		"export_policy_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"rule_index": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 5),
					},

					"allowed_clients": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},

					"protocol": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NFSv3",
								"NFSv4.1",
								"CIFS",
							}, false),
						},
					},

					"unix_read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"unix_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"root_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"kerberos_5_read_only_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"kerberos_5_read_write_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"kerberos_5i_read_only_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"kerberos_5i_read_write_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"kerberos_5p_read_only_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"kerberos_5p_read_write_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"mount_target": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"smb_server_fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"snapshot_directory_visible": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"data_protection_replication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"endpoint_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "dst",
						ValidateFunc: validation.StringInSlice([]string{
							"dst",
						}, false),
					},

					"remote_volume_location": commonschema.Location(),

					"remote_volume_resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"replication_frequency": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"10minutes",
							"daily",
							"hourly",
						}, false),
					},
				},
			},
		},

		"data_protection_snapshot_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"snapshot_policy_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
		},

		"data_protection_backup_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backup_policy_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
						Description:  "The ID of the backup policy to associate with this volume.",
					},

					"policy_enabled": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						Description: "If set to false, the backup policy will not be enabled on this volume, thus disabling scheduled backups.",
					},

					"backup_vault_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
						Description:  "The ID of the backup vault to associate with this volume.",
					},
				},
			},
		},

		"data_protection_advanced_ransomware": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"protection_enabled": {
						Type:        pluginsdk.TypeBool,
						Required:    true,
						Description: "Enable or disable the Advanced Ransomware Protection feature.",
					},
				},
			},
		},

		"azure_vmware_data_store_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
			Default:  false,
		},

		"encryption_key_source": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForEncryptionKeySource(), false),
		},

		"key_vault_private_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: azure.ValidateResourceID,
			RequiredWith: []string{"encryption_key_source"},
		},

		"smb_non_browsable_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable non browsable share setting for SMB/Dual Protocol volume. When enabled, it restricts windows clients to browse the share",
		},

		"smb_access_based_enumeration_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable access based enumeration setting for SMB/Dual Protocol volume. When enabled, users who do not have permission to access a shared folder or file underneath it, do not see that shared resource displayed in their environment.",
		},

		"large_volume_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicates whether the volume is a large volume.",
		},

		"breakthrough_mode_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Indicates whether the large volume runs in Breakthrough Mode, placing it on dedicated capacity that delivers higher throughput and larger capacity.",
		},

		"cool_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"retrieval_policy": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForCoolAccessRetrievalPolicy(), false),
					},

					"tiering_policy": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForCoolAccessTieringPolicy(), false),
					},

					"coolness_period_in_days": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(2, 183),
					},
				},
			},
		},

		"unix_permissions": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "0770",
			ValidateFunc: validation.StringLenBetween(4, 4),
		},
	}
}
