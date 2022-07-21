package netapp

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupVolume struct {
	Name                      string             `tfschema:"name"`
	VolumePath                string             `tfschema:"volume_path"`
	ServiceLevel              string             `tfschema:"service_level"`
	SubnetId                  string             `tfschema:"subnet_id"`
	Protocols                 []string           `tfschema:"protocols"`
	SecurityStyle             string             `tfschema:"security_style"`
	StorageQuotaInGB          int64              `tfschema:"storage_quota_in_gb"`
	ThroughputInMibps         float64            `tfschema:"throughput_in_mibps"`
	Tags                      map[string]string  `tfschema:"tags"`
	SnapshotDirectoryVisible  bool               `tfschema:"snapshot_directory_visible"`
	CapacityPoolId            string             `tfschema:"capacity_pool_id"`
	ProximityPlacementGroupId string             `tfschema:"proximity_placement_group_id"`
	VolumeSpecName            string             `tfschema:"volume_spec_name"`
	ExportPolicy              []ExportPolicyRule `tfschema:"export_policy_rule"`
}

type ExportPolicyRule struct {
	RuleIndex           int      `tfschema:"rule_index"`
	AllowedClients      []string `tfschema:"allowed_clients"`
	CifsEnabled         bool     `tfschema:"cifs_enabled"`
	Nfsv3Enabled        bool     `tfschema:"nfsv3_enabled"`
	Nfsv41Enabled       bool     `tfschema:"nfsv41_enabled"`
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
		result := volumegroups.ExportPolicyRule{
			AllowedClients:      utils.String(strings.Join(item.AllowedClients, ",")),
			Cifs:                utils.Bool(item.CifsEnabled),
			Nfsv3:               utils.Bool(item.Nfsv3Enabled),
			Nfsv41:              utils.Bool(item.Nfsv41Enabled),
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
		protocols := item.Protocols

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

		volumeProperties := &volumegroups.VolumeGroupVolumeProperties{
			Name: utils.String(name),
			Properties: volumegroups.VolumeProperties{
				CapacityPoolResourceId:  utils.String(capacityPoolID),
				CreationToken:           volumePath,
				ServiceLevel:            &serviceLevel,
				SubnetId:                subnetID,
				ProtocolTypes:           &protocols,
				SecurityStyle:           &securityStyle,
				UsageThreshold:          storageQuotaInGB,
				ExportPolicy:            exportPolicyRule,
				ThroughputMibps:         utils.Float(float64(item.ThroughputInMibps)),
				ProximityPlacementGroup: utils.String(item.ProximityPlacementGroupId),
				VolumeSpecName:          utils.String(item.VolumeSpecName),
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

		props := item.Properties

		volumeGroupVolume.Name = string(*item.Name)
		volumeGroupVolume.VolumePath = props.CreationToken
		volumeGroupVolume.ServiceLevel = string(*props.ServiceLevel)
		volumeGroupVolume.SubnetId = props.SubnetId
		volumeGroupVolume.CapacityPoolId = utils.NormalizeNilableString(props.CapacityPoolResourceId)
		volumeGroupVolume.Protocols = *props.ProtocolTypes
		volumeGroupVolume.SecurityStyle = string(*props.SecurityStyle)
		volumeGroupVolume.SnapshotDirectoryVisible = *props.SnapshotDirectoryVisible
		volumeGroupVolume.ThroughputInMibps = float64(*props.ThroughputMibps)
		volumeGroupVolume.Tags = *item.Tags
		volumeGroupVolume.ProximityPlacementGroupId = utils.NormalizeNilableString(props.ProximityPlacementGroup)
		volumeGroupVolume.VolumeSpecName = string(*props.VolumeSpecName)

		if int64(props.UsageThreshold) > 0 {
			usageThreshold := int64(props.UsageThreshold) / 1073741824
			volumeGroupVolume.StorageQuotaInGB = usageThreshold
		}

		if props.ExportPolicy != nil && len(*props.ExportPolicy.Rules) > 0 {
			volumeGroupVolume.ExportPolicy = flattenNetAppVolumeGroupVolumesExportPolicies(props.ExportPolicy.Rules)
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
		rule.CifsEnabled = *item.Cifs
		rule.Nfsv3Enabled = *item.Nfsv3
		rule.Nfsv41Enabled = *item.Nfsv41
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

func getVolumeIndexbyName(input []NetAppVolumeGroupVolume, volumeName string) int {
	for i, item := range input {
		if item.Name == volumeName {
			return i
		}
	}

	return -1
}
