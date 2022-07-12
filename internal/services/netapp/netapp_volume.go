package netapp

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func netAppVolumeCommonSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeName,
		},

		"volume_path": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumePath,
		},

		"service_level": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumes.ServiceLevelPremium),
				string(volumes.ServiceLevelStandard),
				string(volumes.ServiceLevelUltra),
			}, false),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"create_from_snapshot_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: snapshots.ValidateSnapshotID,
		},

		"network_features": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumes.NetworkFeaturesBasic),
				string(volumes.NetworkFeaturesStandard),
			}, false),
		},

		"protocols": {
			Type:     pluginsdk.TypeSet,
			ForceNew: true,
			Optional: true,
			Computed: true,
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

		"security_style": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Unix", // Using hardcoded values instead of SDK enum since no matter what case is passed,
				"Ntfs", // ANF changes casing to Pascal case in the backend. Please refer to https://github.com/Azure/azure-sdk-for-go/issues/14684
			}, false),
		},

		"storage_quota_in_gb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(100, 102400),
		},

		"throughput_in_mibps": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
			Computed: true,
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
							ValidateFunc: validate.CIDR,
						},
					},

					"protocols_enabled": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
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

					"kerberos5_read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"kerberos5_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"kerberos5i_read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"kerberos5i_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"kerberos5p_read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"kerberos5p_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"mount_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"snapshot_directory_visible": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
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

					"remote_volume_location": azure.SchemaLocation(),

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
	}
}

func netAppVolumeSchema() map[string]*pluginsdk.Schema {
	return mergeSchemas(netAppVolumeCommonSchema(), map[string]*pluginsdk.Schema{
		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"pool_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.PoolName,
		},
	})
}

func netAppVolumeGroupVolumeSchema() map[string]*pluginsdk.Schema {
	return mergeSchemas(netAppVolumeCommonSchema(), map[string]*pluginsdk.Schema{
		"capacity_pool_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     azure.ValidateResourceID,
		},

		"proximity_placement_group_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     azure.ValidateResourceID,
		},

		"volume_spec_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	})
}

func mergeSchemas(schemas ...map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {

	result := map[string]*pluginsdk.Schema{}

	for _, schema := range schemas {
		for k, v := range schema {
			if result[k] == nil {
				result[k] = v
			}
		}
	}

	return result
}

type NetAppVolumeGroupVolume struct {
	Name                         string                         `tfschema:"name"`
	VolumePath                   string                         `tfschema:"volume_path"`
	ServiceLevel                 string                         `tfschema:"service_level"`
	SubnetId                     string                         `tfschema:"subnet_id"`
	CreateFromSnapshotResourceId string                         `tfschema:"create_from_snapshot_resource_id"`
	NetworkFeatures              string                         `tfschema:"network_features"`
	Protocols                    []string                       `tfschema:"protocols"`
	SecurityStyle                string                         `tfschema:"security_style"`
	StorageQuotaInGB             int                            `tfschema:"storage_quota_in_gb"`
	ThroughputInMibps            float32                        `tfschema:"throughput_in_mibps"`
	Tags                         map[string]string              `tfschema:"tags"`
	MountIpAddresses             []string                       `tfschema:"mount_ip_addresses"`
	SnapshotDirectoryVisible     string                         `tfschema:"snapshot_directory_visible"`
	CapacityPoolId               string                         `tfschema:"capacity_pool_id"`
	ProximityPlacementGroupId    string                         `tfschema:"proximity_placement_group_id"`
	VolumeSpecName               string                         `tfschema:"volume_spec_name"`
	ExportPolicy                 []ExportPolicyRule             `tfschema:"export_policy_rule"`
	DataProtectionReplication    []DataProtectionReplication    `tfschema:"data_protection_replication"`
	DataProtectionSnapshotPolicy []DataProtectionSnapshotPolicy `tfschema:"data_protection_snapshot_policy"`
}

type ExportPolicyRule struct {
	RuleIndex           int      `tfschema:"rule_index"`
	AllowedClients      string   `tfschema:"allowed_clients"`
	ProtocolsEnabled    []string `tfschema:"protocols_enabled"`
	UnixReadOnly        bool     `tfschema:"unix_read_only"`
	UnixReadWrite       bool     `tfschema:"unix_read_write"`
	RootAccessEnabled   bool     `tfschema:"root_access_enabled"`
	Kerberos5ReadOnly   bool     `tfschema:"kerberos5_read_only"`
	Kerberos5ReadWrite  bool     `tfschema:"kerberos5_read_write"`
	Kerberos5iReadOnly  bool     `tfschema:"kerberos5i_read_only"`
	Kerberos5iReadWrite bool     `tfschema:"kerberos5i_read_write"`
	Kerberos5pReadOnly  bool     `tfschema:"kerberos5p_read_only"`
	Kerberos5pReadWrite bool     `tfschema:"kerberos5p_read_write"`
}

type DataProtectionReplication struct {
	EndpointType           string `tfschema:"endpoint_type"`
	RemoteVolumeLocation   string `tfschema:"remote_volume_location"`
	RemoteVolumeResourceId string `tfschema:"remote_volume_resource_id"`
	ReplicationFrequency   string `tfschema:"replication_frequency"`
}

type DataProtectionSnapshotPolicy struct {
	DataProtectionSnapshotPolicy string `tfschema:"data_protection_snapshot_policy"`
}

func convertExportPolicyToInterface(exportPolicyList []ExportPolicyRule) []interface{} {
	items := make([]interface{}, len(exportPolicyList))
	for i, v := range exportPolicyList {
		items[i] = v
	}

	return items
}

type SupportedObjects interface {
	ExportPolicyRule | DataProtectionReplication | DataProtectionSnapshotPolicy
}

func convertSliceToInterface[T SupportedObjects](slice []T) []interface{} {

	items := make([]interface{}, len(slice))
	for i, v := range slice {
		items[i] = v
	}

	return items
}

func expandNetAppVolumeGroupExportPolicyRule(input []interface{}) *volumegroups.VolumePropertiesExportPolicy {

	if len(input) == 0 || input[0] == nil {
		return &volumegroups.VolumePropertiesExportPolicy{}
	}

	results := make([]volumegroups.ExportPolicyRule, 0)

	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			ruleIndex := int32(v["rule_index"].(int))
			allowedClients := strings.Join(*utils.ExpandStringSlice(v["allowed_clients"].(*pluginsdk.Set).List()), ",")

			cifsEnabled := false
			nfsv3Enabled := false
			nfsv41Enabled := false

			if vpe := v["protocols_enabled"]; vpe != nil {
				protocolsEnabled := vpe.([]interface{})
				if len(protocolsEnabled) != 0 {
					for _, protocol := range protocolsEnabled {
						if protocol != nil {
							switch strings.ToLower(protocol.(string)) {
							case "cifs":
								cifsEnabled = true
							case "nfsv3":
								nfsv3Enabled = true
							case "nfsv4.1":
								nfsv41Enabled = true
							}
						}
					}
				}
			}

			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)
			rootAccessEnabled := v["root_access_enabled"].(bool)
			kerberos5ReadOnly := v["kerberos5_read_only"].(bool)
			kerberos5ReadWrite := v["kerberos5_read_write"].(bool)
			kerberos5iReadOnly := v["kerberos5i_read_only"].(bool)
			kerberos5iReadWrite := v["kerberos5i_read_write"].(bool)
			kerberos5pReadOnly := v["kerberos5p_read_only"].(bool)
			kerberos5pReadWrite := v["kerberos5p_read_write"].(bool)

			result := volumegroups.ExportPolicyRule{
				AllowedClients:      utils.String(allowedClients),
				Cifs:                utils.Bool(cifsEnabled),
				Nfsv3:               utils.Bool(nfsv3Enabled),
				Nfsv41:              utils.Bool(nfsv41Enabled),
				RuleIndex:           utils.Int64(int64(ruleIndex)),
				UnixReadOnly:        utils.Bool(unixReadOnly),
				UnixReadWrite:       utils.Bool(unixReadWrite),
				HasRootAccess:       utils.Bool(rootAccessEnabled),
				Kerberos5ReadOnly:   utils.Bool(kerberos5ReadOnly),
				Kerberos5ReadWrite:  utils.Bool(kerberos5ReadWrite),
				Kerberos5iReadOnly:  utils.Bool(kerberos5iReadOnly),
				Kerberos5iReadWrite: utils.Bool(kerberos5iReadWrite),
				Kerberos5pReadOnly:  utils.Bool(kerberos5pReadOnly),
				Kerberos5pReadWrite: utils.Bool(kerberos5pReadWrite),
			}

			results = append(results, result)
		}
	}

	return &volumegroups.VolumePropertiesExportPolicy{
		Rules: &results,
	}
}

func expandNetAppVolumeGroupDataProtectionReplication(input []interface{}) *volumegroups.VolumePropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumegroups.VolumePropertiesDataProtection{}
	}

	replicationObject := volumegroups.ReplicationObject{}

	replicationRaw := input[0].(map[string]interface{})

	if v, ok := replicationRaw["endpoint_type"]; ok {
		endpointType := volumegroups.EndpointType(v.(string))
		replicationObject.EndpointType = &endpointType
	}
	if v, ok := replicationRaw["remote_volume_location"]; ok {
		replicationObject.RemoteVolumeRegion = utils.String(v.(string))
	}
	if v, ok := replicationRaw["remote_volume_resource_id"]; ok {
		replicationObject.RemoteVolumeResourceId = v.(string)
	}
	if v, ok := replicationRaw["replication_frequency"]; ok {
		replicationSchedule := volumegroups.ReplicationSchedule(translateTFSchedule(v.(string)))
		replicationObject.ReplicationSchedule = &replicationSchedule
	}

	return &volumegroups.VolumePropertiesDataProtection{
		Replication: &replicationObject,
	}
}

func expandNetAppVolumeGroupDataProtectionSnapshotPolicy(input []interface{}) *volumegroups.VolumePropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumegroups.VolumePropertiesDataProtection{}
	}

	snapshotObject := volumegroups.VolumeSnapshotProperties{}

	snapshotRaw := input[0].(map[string]interface{})

	if v, ok := snapshotRaw["snapshot_policy_id"]; ok {
		snapshotObject.SnapshotPolicyId = utils.String(v.(string))
	}

	return &volumegroups.VolumePropertiesDataProtection{
		Snapshot: &snapshotObject,
	}
}
