package netapp

import (
	"fmt"
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
						Computed: true,
					},

					"unix_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"root_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"kerberos5_read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"kerberos5_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"kerberos5i_read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"kerberos5i_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"kerberos5p_read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"kerberos5p_read_write": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
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
	StorageQuotaInGB             int64                          `tfschema:"storage_quota_in_gb"`
	ThroughputInMibps            float64                        `tfschema:"throughput_in_mibps"`
	Tags                         map[string]string              `tfschema:"tags"`
	MountIpAddresses             []string                       `tfschema:"mount_ip_addresses"`
	SnapshotDirectoryVisible     bool                           `tfschema:"snapshot_directory_visible"`
	CapacityPoolId               string                         `tfschema:"capacity_pool_id"`
	ProximityPlacementGroupId    string                         `tfschema:"proximity_placement_group_id"`
	VolumeSpecName               string                         `tfschema:"volume_spec_name"`
	ExportPolicy                 []ExportPolicyRule             `tfschema:"export_policy_rule"`
	DataProtectionReplication    []DataProtectionReplication    `tfschema:"data_protection_replication"`
	DataProtectionSnapshotPolicy []DataProtectionSnapshotPolicy `tfschema:"data_protection_snapshot_policy"`
}

type ExportPolicyRule struct {
	RuleIndex           int      `tfschema:"rule_index"`
	AllowedClients      []string `tfschema:"allowed_clients"`
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

func expandNetAppVolumeGroupExportPolicyRule(input []ExportPolicyRule) *volumegroups.VolumePropertiesExportPolicy {

	if len(input) == 0 || input == nil {
		return &volumegroups.VolumePropertiesExportPolicy{}
	}

	results := make([]volumegroups.ExportPolicyRule, 0)

	for _, item := range input {
		cifsEnabled := false
		nfsv3Enabled := false
		nfsv41Enabled := false

		if len(item.ProtocolsEnabled) != 0 {
			for _, protocol := range item.ProtocolsEnabled {
				if protocol != "" {
					switch strings.ToLower(protocol) {
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

		result := volumegroups.ExportPolicyRule{
			AllowedClients:      utils.String(strings.Join(item.AllowedClients, ",")),
			Cifs:                utils.Bool(cifsEnabled),
			Nfsv3:               utils.Bool(nfsv3Enabled),
			Nfsv41:              utils.Bool(nfsv41Enabled),
			RuleIndex:           utils.Int64(int64(item.RuleIndex)),
			UnixReadOnly:        utils.Bool(item.UnixReadOnly),
			UnixReadWrite:       utils.Bool(item.UnixReadWrite),
			HasRootAccess:       utils.Bool(item.RootAccessEnabled),
			Kerberos5ReadOnly:   utils.Bool(item.Kerberos5ReadOnly),
			Kerberos5ReadWrite:  utils.Bool(item.Kerberos5ReadWrite),
			Kerberos5iReadOnly:  utils.Bool(item.Kerberos5iReadOnly),
			Kerberos5iReadWrite: utils.Bool(item.Kerberos5iReadWrite),
			Kerberos5pReadOnly:  utils.Bool(item.Kerberos5pReadOnly),
			Kerberos5pReadWrite: utils.Bool(item.Kerberos5ReadWrite),
		}

		results = append(results, result)
	}

	return &volumegroups.VolumePropertiesExportPolicy{
		Rules: &results,
	}
}

func expandNetAppVolumeGroupDataProtectionReplication(input []DataProtectionReplication) *volumegroups.VolumePropertiesDataProtection {
	if len(input) == 0 || input == nil {
		return &volumegroups.VolumePropertiesDataProtection{}
	}

	replicationObject := volumegroups.ReplicationObject{}

	endpointType := volumegroups.EndpointType(input[0].EndpointType)
	replicationObject.EndpointType = &endpointType

	replicationObject.RemoteVolumeRegion = &input[0].RemoteVolumeLocation
	replicationObject.RemoteVolumeResourceId = input[0].RemoteVolumeResourceId

	replicationSchedule := volumegroups.ReplicationSchedule(translateTFSchedule(input[0].ReplicationFrequency))
	replicationObject.ReplicationSchedule = &replicationSchedule

	return &volumegroups.VolumePropertiesDataProtection{
		Replication: &replicationObject,
	}
}

func expandNetAppVolumeGroupDataProtectionSnapshotPolicy(input []DataProtectionSnapshotPolicy) *volumegroups.VolumePropertiesDataProtection {
	if len(input) == 0 || input == nil {
		return &volumegroups.VolumePropertiesDataProtection{}
	}

	snapshotObject := volumegroups.VolumeSnapshotProperties{}
	snapshotObject.SnapshotPolicyId = &input[0].DataProtectionSnapshotPolicy

	return &volumegroups.VolumePropertiesDataProtection{
		Snapshot: &snapshotObject,
	}
}

func expandNetAppVolumeGroupVolumes(input []NetAppVolumeGroupVolume, id volumegroups.VolumeGroupId) (*[]volumegroups.VolumeGroupVolumeProperties, error) {

	if len(input) == 0 || input == nil {
		return &[]volumegroups.VolumeGroupVolumeProperties{}, fmt.Errorf("received empty NetAppVolumeGroupVolume slice")
	}

	results := make([]volumegroups.VolumeGroupVolumeProperties, 0)

	for _, item := range input {
		name := item.Name
		volumePath := item.VolumePath
		serviceLevel := volumegroups.ServiceLevel(item.ServiceLevel)
		subnetID := item.SubnetId
		capacityPoolID := item.CapacityPoolId

		networkFeatures := volumegroups.NetworkFeatures(item.NetworkFeatures)
		if networkFeatures == "" {
			networkFeatures = volumegroups.NetworkFeaturesBasic
		}

		protocols := item.Protocols
		if len(protocols) == 0 {
			protocols = append(protocols, "NFSv3")
		}

		// Handling security style property
		securityStyle := volumegroups.SecurityStyle(item.SecurityStyle)
		if strings.EqualFold(string(securityStyle), "unix") && len(protocols) == 1 && strings.EqualFold(protocols[0], "cifs") {
			return &[]volumegroups.VolumeGroupVolumeProperties{}, fmt.Errorf("unix security style cannot be used in a CIFS enabled volume for %s", id)

		}
		if strings.EqualFold(string(securityStyle), "ntfs") && len(protocols) == 1 && (strings.EqualFold(protocols[0], "nfsv3") || strings.EqualFold(protocols[0], "nfsv4.1")) {
			return &[]volumegroups.VolumeGroupVolumeProperties{}, fmt.Errorf("ntfs security style cannot be used in a NFSv3/NFSv4.1 enabled volume for %s", id)
		}

		storageQuotaInGB := int64(item.StorageQuotaInGB * 1073741824)
		exportPolicyRule := expandNetAppVolumeGroupExportPolicyRule(item.ExportPolicy)
		dataProtectionReplication := expandNetAppVolumeGroupDataProtectionReplication(item.DataProtectionReplication)
		dataProtectionSnapshotPolicy := expandNetAppVolumeGroupDataProtectionSnapshotPolicy(item.DataProtectionSnapshotPolicy)

		volumeType := ""
		if dataProtectionReplication != nil && dataProtectionReplication.Replication != nil && strings.ToLower(string(*dataProtectionReplication.Replication.EndpointType)) == "dst" {
			volumeType = "DataProtection"
		}

		// Validating that snapshot policies are not being created in a data protection volume
		if dataProtectionSnapshotPolicy != nil && volumeType != "" {
			return &[]volumegroups.VolumeGroupVolumeProperties{}, fmt.Errorf("snapshot policy cannot be enabled on a data protection volume for %s", id)
		}

		volumeProperties := &volumegroups.VolumeGroupVolumeProperties{
			Name: utils.String(name),
			Properties: volumegroups.VolumeProperties{
				CapacityPoolResourceId:  utils.String(capacityPoolID),
				CreationToken:           volumePath,
				ServiceLevel:            &serviceLevel,
				SubnetId:                subnetID,
				NetworkFeatures:         &networkFeatures,
				ProtocolTypes:           &protocols,
				SecurityStyle:           &securityStyle,
				UsageThreshold:          storageQuotaInGB,
				ExportPolicy:            exportPolicyRule,
				VolumeType:              utils.String(volumeType),
				ThroughputMibps:         utils.Float(float64(item.ThroughputInMibps)),
				ProximityPlacementGroup: utils.String(item.ProximityPlacementGroupId),
				VolumeSpecName:          utils.String(item.VolumeSpecName),
				DataProtection: &volumegroups.VolumePropertiesDataProtection{
					Replication: dataProtectionReplication.Replication,
					Snapshot:    dataProtectionSnapshotPolicy.Snapshot,
				},
				SnapshotDirectoryVisible: &item.SnapshotDirectoryVisible,
			},
			Tags: &item.Tags,
		}

		results = append(results, *volumeProperties)
	}

	return &results, nil
}

func flattenNetAppVolumeGroupVolumes(input *[]volumegroups.VolumeGroupVolumeProperties) ([]NetAppVolumeGroupVolume, error) {
	results := make([]NetAppVolumeGroupVolume, 0)

	if len(*input) == 0 || input == nil {
		return results, fmt.Errorf("received empty volumegroups.VolumeGroupVolumeProperties slice")
	}

	for _, item := range *input {
		volumeGroupVolume := NetAppVolumeGroupVolume{}

		volumeGroupVolume.Name = getResourceNameString(item.Name)

		props := item.Properties

		volumeGroupVolume.VolumePath = props.CreationToken
		volumeGroupVolume.ServiceLevel = string(*props.ServiceLevel)
		volumeGroupVolume.SubnetId = props.SubnetId
		volumeGroupVolume.NetworkFeatures = getNetworkFeaturesString(props.NetworkFeatures)
		volumeGroupVolume.Protocols = *props.ProtocolTypes
		volumeGroupVolume.SecurityStyle = string(*props.SecurityStyle)
		volumeGroupVolume.SnapshotDirectoryVisible = *props.SnapshotDirectoryVisible
		volumeGroupVolume.ThroughputInMibps = float64(*props.ThroughputMibps)

		if int64(props.UsageThreshold) > 0 {
			usageThreshold := int64(props.UsageThreshold) / 1073741824
			volumeGroupVolume.StorageQuotaInGB = usageThreshold
		}

		if props.ExportPolicy != nil && len(*props.ExportPolicy.Rules) > 0 {
			volumeGroupVolume.ExportPolicy = flattenNetAppVolumeGroupVolumesExportPolicies(props.ExportPolicy.Rules)
		}

		if props.MountTargets != nil && len(*props.MountTargets) > 0 {
			volumeGroupVolume.MountIpAddresses = flattenNetAppVolumeGroupVolumesMountIpAddresses(props.MountTargets)
		}

		if props.DataProtection != nil && props.DataProtection.Replication != nil {
			volumeGroupVolume.DataProtectionReplication = flattenNetAppVolumeGroupVolumesDPReplication(props.DataProtection.Replication)
		}

		if props.DataProtection != nil && props.DataProtection.Snapshot != nil {
			volumeGroupVolume.DataProtectionSnapshotPolicy = flattenNetAppVolumeGroupVolumesDPSnapshotPolicy(props.DataProtection.Snapshot)
		}

		results = append(results, volumeGroupVolume)
	}

	return results, nil
}

func flattenNetAppVolumeGroupVolumesExportPolicies(input *[]volumegroups.ExportPolicyRule) []ExportPolicyRule {
	results := make([]ExportPolicyRule, 0)

	if len(*input) == 0 || input == nil {
		return results
	}

	for _, item := range *input {
		rule := ExportPolicyRule{}

		rule.RuleIndex = int(*item.RuleIndex)
		rule.AllowedClients = strings.Split(*item.AllowedClients, ",")

		protocolsEnabled := []string{}
		if *item.Cifs {
			protocolsEnabled = append(protocolsEnabled, "CIFS")
		}
		if *item.Nfsv3 {
			protocolsEnabled = append(protocolsEnabled, "NFSv3")
		}
		if *item.Nfsv41 {
			protocolsEnabled = append(protocolsEnabled, "NFSv4.1")
		}
		rule.ProtocolsEnabled = protocolsEnabled

		rule.UnixReadOnly = *item.UnixReadOnly
		rule.UnixReadWrite = *item.UnixReadWrite
		rule.Kerberos5ReadOnly = *item.Kerberos5ReadOnly
		rule.Kerberos5ReadWrite = *item.Kerberos5ReadWrite
		rule.Kerberos5iReadOnly = *item.Kerberos5iReadOnly
		rule.Kerberos5iReadWrite = *item.Kerberos5iReadWrite
		rule.Kerberos5pReadOnly = *item.Kerberos5pReadOnly
		rule.Kerberos5pReadWrite = *item.Kerberos5pReadWrite

		rule.RootAccessEnabled = *item.HasRootAccess

		results = append(results, rule)
	}

	return results
}

func flattenNetAppVolumeGroupVolumesMountIpAddresses(input *[]volumegroups.MountTargetProperties) []string {
	results := make([]string, 0)

	if len(*input) == 0 || input == nil {
		return results
	}

	for _, item := range *input {
		if item.IpAddress != nil {
			results = append(results, *item.IpAddress)
		}
	}

	return results
}

func flattenNetAppVolumeGroupVolumesDPReplication(input *volumegroups.ReplicationObject) []DataProtectionReplication {
	if input == nil {
		return []DataProtectionReplication{}
	}

	if strings.ToLower(string(*input.EndpointType)) == "" || strings.ToLower(string(*input.EndpointType)) != "dst" {
		return []DataProtectionReplication{}
	}

	return []DataProtectionReplication{
		{
			EndpointType:           string(*input.EndpointType),
			RemoteVolumeLocation:   *input.RemoteVolumeRegion,
			RemoteVolumeResourceId: input.RemoteVolumeResourceId,
			ReplicationFrequency:   string(*input.ReplicationSchedule),
		},
	}
}

func flattenNetAppVolumeGroupVolumesDPSnapshotPolicy(input *volumegroups.VolumeSnapshotProperties) []DataProtectionSnapshotPolicy {
	if input == nil {
		return []DataProtectionSnapshotPolicy{}
	}

	return []DataProtectionSnapshotPolicy{
		{
			DataProtectionSnapshotPolicy: *input.SnapshotPolicyId,
		},
	}
}

func getNetworkFeaturesString(input *volumegroups.NetworkFeatures) string {
	if input == nil {
		return string(volumegroups.NetworkFeaturesBasic)
	}

	return string(*input)
}

func getResourceNameString(input *string) string {
	segments := len(strings.Split(*input, "/"))
	if segments == 0 {
		return ""
	}

	return strings.Split(*input, "/")[segments-1]
}
