// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the NetApp Volume.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VolumeName,
		},

		"location": commonschema.Location(),

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		"account_name": {
			Description:  "The name of the NetApp account in which the NetApp Pool should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"pool_name": {
			Description:  "The name of the NetApp pool in which the NetApp Volume should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.PoolName,
		},

		"volume_path": {
			Description:  "A unique file path for the volume. Used when creating mount targets.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VolumePath,
		},

		"service_level": {
			Description: "The target performance of the file system.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumes.ServiceLevelPremium),
				string(volumes.ServiceLevelStandard),
				string(volumes.ServiceLevelUltra),
				string(volumes.ServiceLevelFlexible),
			}, false),
		},

		"subnet_id": {
			Description:  "The ID of the Subnet the NetApp Volume resides in, which must have the `Microsoft.NetApp/volumes` delegation.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"create_from_snapshot_resource_id": {
			Description:  "Creates volume from snapshot. Following properties must be the same as the original volume where the snapshot was taken from: `protocols`, `subnet_id`, `location`, `service_level`, `resource_group_name` and `account_name`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: snapshots.ValidateSnapshotID,
		},

		"accept_grow_capacity_pool_for_short_term_clone_split": {
			Description:  "While auto splitting the short term clone volume, if the parent pool does not have enough space to accommodate the volume after split, it will be automatically resized, which will lead to increased billing. To accept capacity pool size auto grow and create a short term clone volume, set the property as accepted. Can only be used in conjunction with `create_from_snapshot_resource_id`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForAcceptGrowCapacityPoolForShortTermCloneSplit(), false),
		},

		"network_features": {
			Description: "Which network feature to use. This is a feature in public preview and for more information about it and how to register, refer to [Configure network features for an Azure NetApp Files volume](https://learn.microsoft.com/azure/azure-netapp-files/configure-network-features).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because the service team is changing network features on the backend to upgrade everyone from Basic to Standard and there is a feature that allows customers to change network features from portal but not the API. This could cause drift that forces data loss that we want to avoid
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumes.NetworkFeaturesBasic),
				string(volumes.NetworkFeaturesStandard),
			}, false),
		},

		"protocols": {
			Description: "The target volume protocol expressed as a list. Protocol conversion between `NFSv3` and `NFSv4.1` and vice-versa is supported without recreating the volume, however export policy rules must be updated accordingly to avoid configuration drift (e.g., when converting from `NFSv3` to `NFSv4.1`, set `nfsv3_enabled = false` and `nfsv41_enabled = true` in export policy rules). Dual protocol scenario is supported for CIFS and NFSv3, for more information, refer to [Create a dual-protocol volume for Azure NetApp Files](https://learn.microsoft.com/azure/azure-netapp-files/create-volumes-dual-protocol) document.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    2,
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
			Description: "Enable to allow Kerberos secured volumes. Requires appropriate export rules as well as the parent `azurerm_netapp_account` having a defined AD connection.",
			// Due to large infrastructure requirements, there is not a reliable way to test Kerberos volumes in a shared environment at this time
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"smb_continuous_availability_enabled": {
			Description: "Continuous availability option should be used only for SQL and FSLogix workloads. Using it for any other SMB workloads is not supported.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"smb3_protocol_encryption_enabled": {
			Description: "SMB3 encryption option should be used only for SMB/DualProtocol volumes. Using it for any other workloads is not supported.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"security_style": {
			Description:  "Volume security style. In a dual-protocol volume, if not provided, its value will be `ntfs`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForSecurityStyle(), false),
		},

		"storage_quota_in_gb": {
			Description: "The maximum Storage Quota allowed for a file system in Gigabytes.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
			// The upper bound covers the largest size any volume flavour supports (2,400 TiB, only reachable
			// with `breakthrough_mode_enabled`), the per-flavour limits are enforced in CustomizeDiff.
			ValidateFunc: validation.IntBetween(50, 2457600),
		},

		"throughput_in_mibps": {
			Description:  "Throughput of this volume in Mibps.",
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
						Description:  "The index number of the rule.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 5),
					},

					"allowed_clients": {
						Description: "A list of allowed clients IPv4 addresses.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},

					"protocol": {
						Description: "A list of allowed protocols. Only a single element is supported at this time.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MinItems:    1,
						MaxItems:    1,
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
						Description: "Whether the unix filesystem is read-only.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"unix_read_write": {
						Description: "Whether the unix filesystem is read-write.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"root_access_enabled": {
						Description: "Whether root access is permitted to this volume.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"kerberos_5_read_only_enabled": {
						Description: "Whether Kerberos 5 read-only access is permitted to this volume.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"kerberos_5_read_write_enabled": {
						Description: "Whether Kerberos 5 read/write is permitted to this volume.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"kerberos_5i_read_only_enabled": {
						Description: "Whether Kerberos 5i read-only is permitted to this volume.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"kerberos_5i_read_write_enabled": {
						Description: "Whether Kerberos 5i read/write is permitted to this volume.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"kerberos_5p_read_only_enabled": {
						Description: "Whether Kerberos 5p read-only is permitted to this volume.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"kerberos_5p_read_write_enabled": {
						Description: "Whether Kerberos 5p read/write is permitted to this volume.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
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
						Description: "The IP address of the mount target.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"smb_server_fqdn": {
						Description: "The SMB server's Fully Qualified Domain Name (FQDN). This value is populated when the volume's `protocols` include `CIFS`; otherwise, it is empty.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"snapshot_directory_visible": {
			Description: "Whether the .snapshot (NFS clients) or ~snapshot (SMB clients) path of a volume is visible.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"data_protection_replication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"endpoint_type": {
						Description: "The endpoint type, default value is `dst` for destination.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     "dst",
						ValidateFunc: validation.StringInSlice([]string{
							"dst",
						}, false),
					},

					"remote_volume_location": commonschema.Location(),

					"remote_volume_resource_id": {
						Description:  "Resource ID of the primary volume.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"replication_frequency": {
						Description: "Replication frequency.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description:  "Resource ID of the snapshot policy to apply to the volume.",
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
						Description:  "The ID of the backup policy to associate with this volume.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"policy_enabled": {
						Description: "If set to false, the backup policy will not be enabled on this volume, thus disabling scheduled backups.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"backup_vault_id": {
						Description:  "The ID of the Backup Vault to associate with this volume.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
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
						Description: "Enable or disable the Advanced Ransomware Protection feature.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
				},
			},
		},

		"azure_vmware_data_store_enabled": {
			Description: "Whether the NetApp Volume is enabled for Azure VMware Solution (AVS) datastore purpose.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"encryption_key_source": {
			Description:  "The encryption key source, it can be `Microsoft.NetApp` for platform managed keys or `Microsoft.KeyVault` for customer-managed keys.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForEncryptionKeySource(), false),
		},

		"key_vault_private_endpoint_id": {
			Description:  "The Private Endpoint ID for Key Vault, which is required when using customer-managed keys.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
			RequiredWith: []string{"encryption_key_source"},
		},

		"smb_non_browsable_enabled": {
			Description: "Enable non browsable share setting for SMB/Dual Protocol volume. When enabled, it restricts windows clients to browse the share.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"smb_access_based_enumeration_enabled": {
			Description: "Enable access based enumeration setting for SMB/Dual Protocol volume. When enabled, users who do not have permission to access a shared folder or file underneath it, do not see that shared resource displayed in their environment.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"large_volume_enabled": {
			Description: "Whether the volume is a large volume.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"breakthrough_mode_enabled": {
			Description: "Whether the large volume runs in Breakthrough Mode, placing it on dedicated capacity that delivers higher throughput and larger capacity.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"cool_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"retrieval_policy": {
						Description:  "The cool access retrieval policy for the volume.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForCoolAccessRetrievalPolicy(), false),
					},

					"tiering_policy": {
						Description:  "The cool access tiering policy for the volume.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForCoolAccessTieringPolicy(), false),
					},

					"coolness_period_in_days": {
						Description:  "The coolness period in days for the volume. Possible vales are between `2` and `183`.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(2, 183),
					},
				},
			},
		},

		"unix_permissions": {
			Description: "The unix permissions for an NFS volume are accepted in a four-digit octal format. The first digit selects the set user ID (4), set group ID (2), and sticky bit (1) attributes. The second digit selects the permissions for the owner of the file: read (4), write (2), and execute (1). The third digit selects the permissions for users in the same group. The fourth digit selects the permissions for users not in the group.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because Azure returns a default of `0770` for some volume types, while others reject this default value with a 400 when sent to the API
			// unexpected status 400 (400 Bad Request) with error: InvalidInputUnixPermissionForNtfsSecurityStyleVolume: Unix permissions are not allowed for specified volume configuration, Create SMB or dual protocol volume with NTFS security style, without unix permissions.
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(4, 4),
		},
	}
}
