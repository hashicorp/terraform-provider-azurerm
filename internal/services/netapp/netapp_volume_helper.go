package netapp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumesreplication"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupVolume struct {
	Id                           string                         `tfschema:"id"`
	Name                         string                         `tfschema:"name"`
	VolumePath                   string                         `tfschema:"volume_path"`
	ServiceLevel                 string                         `tfschema:"service_level"`
	SubnetId                     string                         `tfschema:"subnet_id"`
	Protocols                    []string                       `tfschema:"protocols"`
	SecurityStyle                string                         `tfschema:"security_style"`
	StorageQuotaInGB             int64                          `tfschema:"storage_quota_in_gb"`
	ThroughputInMibps            float64                        `tfschema:"throughput_in_mibps"`
	Tags                         map[string]string              `tfschema:"tags"`
	SnapshotDirectoryVisible     bool                           `tfschema:"snapshot_directory_visible"`
	CapacityPoolId               string                         `tfschema:"capacity_pool_id"`
	ProximityPlacementGroupId    string                         `tfschema:"proximity_placement_group_id"`
	VolumeSpecName               string                         `tfschema:"volume_spec_name"`
	ExportPolicy                 []ExportPolicyRule             `tfschema:"export_policy_rule"`
	MountIpAddresses             []string                       `tfschema:"mount_ip_addresses"`
	DataProtectionReplication    []DataProtectionReplication    `tfschema:"data_protection_replication"`
	DataProtectionSnapshotPolicy []DataProtectionSnapshotPolicy `tfschema:"data_protection_snapshot_policy"`
}

type ExportPolicyRule struct {
	RuleIndex         int    `tfschema:"rule_index"`
	AllowedClients    string `tfschema:"allowed_clients"`
	Nfsv3Enabled      bool   `tfschema:"nfsv3_enabled"`
	Nfsv41Enabled     bool   `tfschema:"nfsv41_enabled"`
	UnixReadOnly      bool   `tfschema:"unix_read_only"`
	UnixReadWrite     bool   `tfschema:"unix_read_write"`
	RootAccessEnabled bool   `tfschema:"root_access_enabled"`
}

type DataProtectionReplication struct {
	EndpointType           string `tfschema:"endpoint_type"`
	RemoteVolumeLocation   string `tfschema:"remote_volume_location"`
	RemoteVolumeResourceId string `tfschema:"remote_volume_resource_id"`
	ReplicationFrequency   string `tfschema:"replication_frequency"`
}

type DataProtectionSnapshotPolicy struct {
	DataProtectionSnapshotPolicy string `tfschema:"snapshot_policy_id"`
}

type VolumeSpecNameSapHana string

const (
	VolumeSpecNameSapHanaData       VolumeSpecNameSapHana = "data"
	VolumeSpecNameSapHanaLog        VolumeSpecNameSapHana = "log"
	VolumeSpecNameSapHanaShared     VolumeSpecNameSapHana = "shared"
	VolumeSpecNameSapHanaDataBackup VolumeSpecNameSapHana = "data-backup"
	VolumeSpecNameSapHanaLogBackup  VolumeSpecNameSapHana = "log-backup"
)

func PossibleValuesForVolumeSpecNameSapHana() []string {
	return []string{
		string(VolumeSpecNameSapHanaData),
		string(VolumeSpecNameSapHanaLog),
		string(VolumeSpecNameSapHanaShared),
		string(VolumeSpecNameSapHanaDataBackup),
		string(VolumeSpecNameSapHanaLogBackup),
	}
}

func RequiredVolumesForSAPHANA() []string {
	return []string{
		string(VolumeSpecNameSapHanaData),
		string(VolumeSpecNameSapHanaLog),
	}
}

type ProtocolType string

const (
	ProtocolTypeNfsV41 ProtocolType = "NFSv4.1"
	ProtocolTypeNfsV3  ProtocolType = "NFSv3"
	ProtocolTypeCifs   ProtocolType = "CIFS"
)

func PossibleValuesForProtocolType() []string {
	return []string{
		string(ProtocolTypeNfsV41),
		string(ProtocolTypeNfsV3),
		string(ProtocolTypeCifs),
	}
}

func PossibleValuesForProtocolTypeVolumeGroupSapHana() []string {
	return []string{
		string(ProtocolTypeNfsV41),
		string(ProtocolTypeNfsV3),
	}
}

type ReplicationSchedule string

const (
	ReplicationSchedule10Minutes ReplicationSchedule = "10minutes"
	ReplicationScheduleDaily     ReplicationSchedule = "daily"
	ReplicationScheduleHourly    ReplicationSchedule = "hourly"
)

func PossibleValuesForReplicationSchedule() []string {
	return []string{
		string(ReplicationSchedule10Minutes),
		string(ReplicationScheduleDaily),
		string(ReplicationScheduleHourly),
	}
}

// Diverging from the SDK volumegroups.SecurityStyle since it is defined as lower case
// but the backend changes it to Pascal case on GET. Please refer to https://github.com/Azure/azure-sdk-for-go/issues/14684
type SecurityStyle string

const (
	SecurityStyleUnix SecurityStyle = "Unix"
	SecurityStyleNtfs SecurityStyle = "Ntfs"
)

func PossibleValuesForSecurityStyle() []string {
	return []string{
		string(SecurityStyleUnix),
	}
}

func expandNetAppVolumeGroupVolumeExportPolicyRule(input []ExportPolicyRule) *volumegroups.VolumePropertiesExportPolicy {

	if len(input) == 0 || input == nil {
		return &volumegroups.VolumePropertiesExportPolicy{}
	}

	results := make([]volumegroups.ExportPolicyRule, 0)

	for _, item := range input {

		// Hard-Coded values, for AVG these cannot be set differently
		// they are not exposed as TF configuration
		// but PUT request requires those fields to succeed
		cifsEnabled := false
		kerberos5ReadOnly := false
		kerberos5ReadWrite := false
		kerberos5iReadOnly := false
		kerberos5iReadWrite := false
		kerberos5pReadOnly := false
		kerberos5pReadWrite := false

		result := volumegroups.ExportPolicyRule{
			AllowedClients:      utils.String(item.AllowedClients),
			Cifs:                utils.Bool(cifsEnabled),
			Nfsv3:               utils.Bool(item.Nfsv3Enabled),
			Nfsv41:              utils.Bool(item.Nfsv41Enabled),
			RuleIndex:           utils.Int64(int64(item.RuleIndex)),
			UnixReadOnly:        utils.Bool(item.UnixReadOnly),
			UnixReadWrite:       utils.Bool(item.UnixReadWrite),
			HasRootAccess:       utils.Bool(item.RootAccessEnabled),
			Kerberos5ReadOnly:   utils.Bool(kerberos5ReadOnly),
			Kerberos5ReadWrite:  utils.Bool(kerberos5ReadWrite),
			Kerberos5iReadOnly:  utils.Bool(kerberos5iReadOnly),
			Kerberos5iReadWrite: utils.Bool(kerberos5iReadWrite),
			Kerberos5pReadOnly:  utils.Bool(kerberos5pReadOnly),
			Kerberos5pReadWrite: utils.Bool(kerberos5pReadWrite),
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

func expandNetAppVolumeGroupVolumes(input []NetAppVolumeGroupVolume) (*[]volumegroups.VolumeGroupVolumeProperties, error) {
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
		snapshotDirectoryVisible := item.SnapshotDirectoryVisible
		securityStyle := volumegroups.SecurityStyle(item.SecurityStyle)
		storageQuotaInGB := item.StorageQuotaInGB * 1073741824
		proximityPlacementGroupId := utils.NormalizeNilableString(&item.ProximityPlacementGroupId)
		exportPolicyRule := expandNetAppVolumeGroupVolumeExportPolicyRule(item.ExportPolicy)
		dataProtectionReplication := expandNetAppVolumeGroupDataProtectionReplication(item.DataProtectionReplication)
		dataProtectionSnapshotPolicy := expandNetAppVolumeGroupDataProtectionSnapshotPolicy(item.DataProtectionSnapshotPolicy)

		volumeProperties := &volumegroups.VolumeGroupVolumeProperties{
			Name: utils.String(name),
			Properties: volumegroups.VolumeProperties{
				CapacityPoolResourceId:   utils.String(capacityPoolID),
				CreationToken:            volumePath,
				ServiceLevel:             &serviceLevel,
				SubnetId:                 subnetID,
				ProtocolTypes:            &protocols,
				SecurityStyle:            &securityStyle,
				UsageThreshold:           storageQuotaInGB,
				ExportPolicy:             exportPolicyRule,
				SnapshotDirectoryVisible: utils.Bool(snapshotDirectoryVisible),
				ThroughputMibps:          utils.Float(item.ThroughputInMibps),
				ProximityPlacementGroup:  &proximityPlacementGroupId,
				VolumeSpecName:           utils.String(item.VolumeSpecName),
				DataProtection: &volumegroups.VolumePropertiesDataProtection{
					Replication: dataProtectionReplication.Replication,
					Snapshot:    dataProtectionSnapshotPolicy.Snapshot,
				},
			},
			Tags: &item.Tags,
		}

		results = append(results, *volumeProperties)
	}

	return &results, nil
}

func expandNetAppVolumeGroupVolumeExportPolicyRulePatch(input []interface{}) *volumes.VolumePatchPropertiesExportPolicy {
	results := make([]volumes.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})

			ruleIndex := int64(v["rule_index"].(int))
			allowedClients := v["allowed_clients"].(string)
			nfsv3Enabled := v["nfsv3_enabled"].(bool)
			nfsv41Enabled := v["nfsv41_enabled"].(bool)
			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)
			rootAccessEnabled := v["root_access_enabled"].(bool)

			// Hard-Coded values, for AVG these cannot be set differently
			// they are not exposed as TF configuration
			// but PUT request requires those fields to succeed
			cifsEnabled := false
			kerberos5ReadOnly := false
			kerberos5ReadWrite := false
			kerberos5iReadOnly := false
			kerberos5iReadWrite := false
			kerberos5pReadOnly := false
			kerberos5pReadWrite := false

			result := volumes.ExportPolicyRule{
				AllowedClients:      utils.String(allowedClients),
				Cifs:                utils.Bool(cifsEnabled),
				Nfsv3:               utils.Bool(nfsv3Enabled),
				Nfsv41:              utils.Bool(nfsv41Enabled),
				RuleIndex:           utils.Int64(ruleIndex),
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

	return &volumes.VolumePatchPropertiesExportPolicy{
		Rules: &results,
	}
}

func expandNetAppVolumeDataProtectionReplication(input []interface{}) *volumes.VolumePropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumes.VolumePropertiesDataProtection{}
	}

	replicationObject := volumes.ReplicationObject{}

	replicationRaw := input[0].(map[string]interface{})

	if v, ok := replicationRaw["endpoint_type"]; ok {
		endpointType := volumes.EndpointType(v.(string))
		replicationObject.EndpointType = &endpointType
	}
	if v, ok := replicationRaw["remote_volume_location"]; ok {
		replicationObject.RemoteVolumeRegion = utils.String(v.(string))
	}
	if v, ok := replicationRaw["remote_volume_resource_id"]; ok {
		replicationObject.RemoteVolumeResourceId = v.(string)
	}
	if v, ok := replicationRaw["replication_frequency"]; ok {
		replicationSchedule := volumes.ReplicationSchedule(translateTFSchedule(v.(string)))
		replicationObject.ReplicationSchedule = &replicationSchedule
	}

	return &volumes.VolumePropertiesDataProtection{
		Replication: &replicationObject,
	}
}

func expandNetAppVolumeDataProtectionSnapshotPolicy(input []interface{}) *volumes.VolumePropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumes.VolumePropertiesDataProtection{}
	}

	snapshotObject := volumes.VolumeSnapshotProperties{}

	snapshotRaw := input[0].(map[string]interface{})

	if v, ok := snapshotRaw["snapshot_policy_id"]; ok {
		snapshotObject.SnapshotPolicyId = utils.String(v.(string))
	}

	return &volumes.VolumePropertiesDataProtection{
		Snapshot: &snapshotObject,
	}
}

func expandNetAppVolumeDataProtectionSnapshotPolicyPatch(input []interface{}) *volumes.VolumePatchPropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumes.VolumePatchPropertiesDataProtection{}
	}

	snapshotObject := volumes.VolumeSnapshotProperties{}

	snapshotRaw := input[0].(map[string]interface{})

	if v, ok := snapshotRaw["snapshot_policy_id"]; ok {
		snapshotObject.SnapshotPolicyId = utils.String(v.(string))
	}

	return &volumes.VolumePatchPropertiesDataProtection{
		Snapshot: &snapshotObject,
	}
}

func flattenNetAppVolumeGroupVolumes(ctx context.Context, input *[]volumegroups.VolumeGroupVolumeProperties, metadata sdk.ResourceMetaData) ([]NetAppVolumeGroupVolume, error) {
	results := make([]NetAppVolumeGroupVolume, 0)

	if len(*input) == 0 || input == nil {
		return results, fmt.Errorf("received empty volumegroups.VolumeGroupVolumeProperties slice")
	}

	for _, item := range *input {
		volumeGroupVolume := NetAppVolumeGroupVolume{}

		props := item.Properties

		volumeGroupVolume.Name = getResourceNameString(item.Name)
		volumeGroupVolume.VolumePath = props.CreationToken
		volumeGroupVolume.ServiceLevel = string(*props.ServiceLevel)
		volumeGroupVolume.SubnetId = props.SubnetId
		volumeGroupVolume.CapacityPoolId = utils.NormalizeNilableString(props.CapacityPoolResourceId)
		volumeGroupVolume.Protocols = *props.ProtocolTypes
		volumeGroupVolume.SecurityStyle = string(*props.SecurityStyle)
		volumeGroupVolume.SnapshotDirectoryVisible = *props.SnapshotDirectoryVisible
		volumeGroupVolume.ThroughputInMibps = *props.ThroughputMibps
		volumeGroupVolume.Tags = *item.Tags
		volumeGroupVolume.ProximityPlacementGroupId = utils.NormalizeNilableString(props.ProximityPlacementGroup)
		volumeGroupVolume.VolumeSpecName = *props.VolumeSpecName

		if props.UsageThreshold > 0 {
			usageThreshold := props.UsageThreshold / 1073741824
			volumeGroupVolume.StorageQuotaInGB = usageThreshold
		}

		if props.ExportPolicy != nil && len(*props.ExportPolicy.Rules) > 0 {
			volumeGroupVolume.ExportPolicy = flattenNetAppVolumeGroupVolumesExportPolicies(props.ExportPolicy.Rules)
		}

		if props.MountTargets != nil && len(*props.MountTargets) > 0 {
			volumeGroupVolume.MountIpAddresses = flattenNetAppVolumeGroupVolumesMountIpAddresses(props.MountTargets)
		}

		// Getting volume resource directly from standalone volume
		// since VolumeGroup Volumes don't return DataProtection information
		volumeClient := metadata.Client.NetApp.VolumeClient
		id, err := volumes.ParseVolumeID(*item.Id)
		if err != nil {
			return []NetAppVolumeGroupVolume{}, err
		}

		standaloneVol, err := volumeClient.Get(ctx, *id)
		if err != nil {
			return []NetAppVolumeGroupVolume{}, fmt.Errorf("retrieving %s: %v", id, err)
		}

		if standaloneVol.Model.Properties.DataProtection != nil && standaloneVol.Model.Properties.DataProtection.Replication != nil {
			volumeGroupVolume.DataProtectionReplication = flattenNetAppVolumeGroupVolumesDPReplication(standaloneVol.Model.Properties.DataProtection.Replication)
		}

		if standaloneVol.Model.Properties.DataProtection != nil && standaloneVol.Model.Properties.DataProtection.Snapshot != nil {
			volumeGroupVolume.DataProtectionSnapshotPolicy = flattenNetAppVolumeGroupVolumesDPSnapshotPolicy(standaloneVol.Model.Properties.DataProtection.Snapshot)
		}

		volumeGroupVolume.Id = *standaloneVol.Model.Id

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
		rule.AllowedClients = *item.AllowedClients
		rule.Nfsv3Enabled = *item.Nfsv3
		rule.Nfsv41Enabled = *item.Nfsv41
		rule.UnixReadOnly = *item.UnixReadOnly
		rule.UnixReadWrite = *item.UnixReadWrite
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
		if item.IPAddress != nil {
			results = append(results, *item.IPAddress)
		}
	}

	return results
}

func flattenNetAppVolumeGroupVolumesDPReplication(input *volumes.ReplicationObject) []DataProtectionReplication {
	if input == nil {
		return []DataProtectionReplication{}
	}
	if string(*input.EndpointType) == "" || !strings.EqualFold(string(*input.EndpointType), string(volumes.EndpointTypeDst)) {
		return []DataProtectionReplication{}
	}

	replicationFrequency := ""
	if input.ReplicationSchedule != nil {
		replicationFrequency = translateSDKSchedule(strings.ToLower(string(*input.ReplicationSchedule)))
	}

	return []DataProtectionReplication{
		{
			EndpointType:           strings.ToLower(string(*input.EndpointType)),
			RemoteVolumeLocation:   *input.RemoteVolumeRegion,
			RemoteVolumeResourceId: input.RemoteVolumeResourceId,
			ReplicationFrequency:   replicationFrequency,
		},
	}
}

func flattenNetAppVolumeGroupVolumesDPSnapshotPolicy(input *volumes.VolumeSnapshotProperties) []DataProtectionSnapshotPolicy {
	if input == nil {
		return []DataProtectionSnapshotPolicy{}
	}

	return []DataProtectionSnapshotPolicy{
		{
			DataProtectionSnapshotPolicy: *input.SnapshotPolicyId,
		},
	}
}

func getResourceNameString(input *string) string {
	segments := len(strings.Split(*input, "/"))
	if segments == 0 {
		return ""
	}

	return strings.Split(*input, "/")[segments-1]
}

func deleteVolume(ctx context.Context, metadata sdk.ResourceMetaData, volumeId string) error {
	client := metadata.Client.NetApp.VolumeClient

	id, err := volumes.ParseVolumeID(volumeId)
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if existing.HttpResponse.StatusCode == http.StatusNotFound {
			return metadata.MarkAsGone(id)
		}
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	// Removing replication if present
	if existing.Model.Properties.DataProtection != nil && existing.Model.Properties.DataProtection.Replication != nil {
		dataProtectionReplication := existing.Model.Properties.DataProtection
		replicaVolumeId, err := volumesreplication.ParseVolumeID(id.ID())
		if err != nil {
			return err
		}
		if dataProtectionReplication.Replication.EndpointType != nil && !strings.EqualFold(string(*dataProtectionReplication.Replication.EndpointType), string(volumes.EndpointTypeDst)) {
			// This is the case where primary volume started the deletion, in this case, to be consistent we will remove replication from secondary
			replicaVolumeId, err = volumesreplication.ParseVolumeID(dataProtectionReplication.Replication.RemoteVolumeResourceId)
			if err != nil {
				return err
			}
		}

		replicationClient := metadata.Client.NetApp.VolumeReplicationClient
		// Checking replication status before deletion, it need to be broken before proceeding with deletion
		if res, err := replicationClient.VolumesReplicationStatus(ctx, *replicaVolumeId); err == nil {
			// Wait for replication state = "mirrored"
			if model := res.Model; model != nil {
				if model.MirrorState != nil && strings.ToLower(string(*model.MirrorState)) == "uninitialized" {
					if err := waitForReplMirrorState(ctx, replicationClient, *replicaVolumeId, "mirrored"); err != nil {
						return fmt.Errorf("waiting for replica %s to become 'mirrored': %+v", *replicaVolumeId, err)
					}
				}
			}

			// Breaking replication
			if err = replicationClient.VolumesBreakReplicationThenPoll(ctx, *replicaVolumeId, volumesreplication.BreakReplicationRequest{
				ForceBreakReplication: utils.Bool(true),
			}); err != nil {
				return fmt.Errorf("breaking replication for %s: %+v", *replicaVolumeId, err)
			}

			// Waiting for replication be in broken state
			metadata.Logger.Infof("waiting for the replication of %s to be in broken state", *replicaVolumeId)
			if err := waitForReplMirrorState(ctx, replicationClient, *replicaVolumeId, "broken"); err != nil {
				return fmt.Errorf("waiting for the breaking of replication for %s: %+v", *replicaVolumeId, err)
			}
		}

		// Deleting replication and waiting for it to fully complete the operation
		if err = replicationClient.VolumesDeleteReplicationThenPoll(ctx, *replicaVolumeId); err != nil {
			return fmt.Errorf("deleting replicate %s: %+v", *replicaVolumeId, err)
		}

		if err := waitForReplicationDeletion(ctx, replicationClient, *replicaVolumeId); err != nil {
			return fmt.Errorf("waiting for the replica %s to be deleted: %+v", *replicaVolumeId, err)
		}
	}

	// Deleting volume and waiting for it fo fully complete the operation
	if err = client.DeleteThenPoll(ctx, *id, volumes.DeleteOperationOptions{
		ForceDelete: utils.Bool(true),
	}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = waitForVolumeDeletion(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func waitForVolumeCreateOrUpdate(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish creating: %+v", id, err)
	}

	return nil
}

func waitForVolumeGroupCreateOrUpdate(ctx context.Context, client *volumegroups.VolumeGroupsClient, id volumegroups.VolumeGroupId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeGroupStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish creating: %+v", id, err)
	}

	return nil
}

func waitForReplAuthorization(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404", "400"}, // TODO: Remove 400 when bug is fixed on RP side, where replicationStatus returns 400 at some point during authorization process
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeReplicationStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for replication authorization %s to complete: %+v", id, err)
	}

	return nil
}

func waitForReplMirrorState(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId, desiredState string) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200"}, // 200 means mirror state is still Mirrored
		Target:                    []string{"204"}, // 204 means mirror state is <> than Mirrored
		Refresh:                   netappVolumeReplicationMirrorStateRefreshFunc(ctx, client, id, desiredState),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be in the state %q: %+v", id, desiredState, err)
	}

	return nil
}

func waitForReplicationDeletion(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202", "400"}, // TODO: Remove 400 when bug is fixed on RP side, where replicationStatus returns 400 while it is in "Deleting" state
		Target:                    []string{"404"},
		Refresh:                   netappVolumeReplicationStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Replication of %s to be deleted: %+v", id, err)
	}

	return nil
}

func waitForVolumeDeletion(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappVolumeStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func netappVolumeStateRefreshFunc(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}

func netappVolumeGroupStateRefreshFunc(ctx context.Context, client *volumegroups.VolumeGroupsClient, id volumegroups.VolumeGroupId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.VolumeGroupsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}

func netappVolumeReplicationMirrorStateRefreshFunc(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId, desiredState string) pluginsdk.StateRefreshFunc {
	validStates := []string{"mirrored", "broken", "uninitialized"}

	return func() (interface{}, string, error) {
		// Possible Mirror States to be used as desiredStates:
		// mirrored, broken or uninitialized
		if !utils.SliceContainsValue(validStates, strings.ToLower(desiredState)) {
			return nil, "", fmt.Errorf("invalid desired mirror state was passed to check mirror replication state (%s), possible values: (%+v)", desiredState, volumesreplication.PossibleValuesForMirrorState())
		}

		res, err := client.VolumesReplicationStatus(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving replication status information from %s: %s", id, err)
			}
		}

		// TODO: fix this refresh function to use strings instead of fake status codes
		// Setting 200 as default response
		response := 200
		if res.Model != nil && res.Model.MirrorState != nil && strings.EqualFold(string(*res.Model.MirrorState), desiredState) {
			// return 204 if state matches desired state
			response = 204
		}

		return res, strconv.Itoa(response), nil
	}
}

func netappVolumeReplicationStateRefreshFunc(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.VolumesReplicationStatus(ctx, id)
		if err != nil {
			if response.WasBadRequest(res.HttpResponse) && (strings.Contains(strings.ToLower(err.Error()), "deleting") || strings.Contains(strings.ToLower(err.Error()), "volume replication missing or deleted")) {
				// This error can be ignored until a bug is fixed on RP side that it is returning 400 while the replication is in "Deleting" process
				// TODO: remove this workaround when above bug is fixed
			} else if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving replication status from %s: %s", id, err)
			}
		}
		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}

func translateTFSchedule(scheduleName string) string {
	if strings.EqualFold(scheduleName, string(ReplicationSchedule10Minutes)) {
		return string(volumegroups.ReplicationScheduleOneZerominutely)
	}

	return scheduleName
}

func translateSDKSchedule(scheduleName string) string {
	if strings.EqualFold(scheduleName, string(volumegroups.ReplicationScheduleOneZerominutely)) {
		return string(ReplicationSchedule10Minutes)
	}

	return scheduleName
}

func validateNetAppVolumeGroupVolumes(volumeList *[]volumegroups.VolumeGroupVolumeProperties, applicationType volumegroups.ApplicationType) []error {
	errors := make([]error, 0)
	volumeSpecRepeatCount := make(map[string]int)

	if applicationType == volumegroups.ApplicationTypeSAPNegativeHANA {

		// Validating maximum number of volumes
		if len(*volumeList) > 5 {
			errors = append(errors, fmt.Errorf("'`volume` list cannot be greater than 5 for %v'", applicationType))
		}

		// Validating each volume
		for _, volume := range *volumeList {
			// Validating maximum number of protocols
			if len(*volume.Properties.ProtocolTypes) > 1 {
				errors = append(errors, fmt.Errorf("'`protocols` list cannot be greater than 1 for %v on volume %v'", applicationType, *volume.Name))
			}

			// Validating protocol, it supports only one and that is enforced by the schema

			// Can't be CIFS at all times
			if strings.EqualFold((*volume.Properties.ProtocolTypes)[0], string(ProtocolTypeCifs)) {
				errors = append(errors, fmt.Errorf("'cifs is not supported for %v on volume %v'", applicationType, *volume.Name))
			}

			// Can't be nfsv3 on data, log and share volumes
			if strings.EqualFold((*volume.Properties.ProtocolTypes)[0], string(ProtocolTypeNfsV3)) &&
				(strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaData)) ||
					strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaShared)) ||
					strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaLog))) {

				errors = append(errors, fmt.Errorf("'nfsv3 on data, log and shared volumes for %v is not supported on volume %v'", applicationType, *volume.Name))
			}

			// Validating export policies
			if volume.Properties.ExportPolicy != nil {
				for _, rule := range *volume.Properties.ExportPolicy.Rules {
					errors = append(errors, validateNetAppVolumeGroupExportPolicyRule(rule, (*volume.Properties.ProtocolTypes)[0])...)
				}
			}

			// Checking CRR rule that log cannot be DataProtection type
			if strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaLog)) &&
				volume.Properties.DataProtection != nil &&
				volume.Properties.DataProtection.Replication != nil &&
				strings.EqualFold(string(*volume.Properties.DataProtection.Replication.EndpointType), string(volumegroups.EndpointTypeDst)) {

				errors = append(errors, fmt.Errorf("'log volume spec type cannot be DataProtection type for %v on volume %v'", applicationType, *volume.Name))
			}

			// Validating that snapshot policies are not being created in a data protection volume
			if volume.Properties.DataProtection != nil &&
				volume.Properties.DataProtection.Snapshot != nil &&
				(volume.Properties.DataProtection.Replication != nil && strings.EqualFold(string(*volume.Properties.DataProtection.Replication.EndpointType), string(volumegroups.EndpointTypeDst))) {

				errors = append(errors, fmt.Errorf("'snapshot policy cannot be enabled on a data protection volume for %v on volume %v'", applicationType, *volume.Name))
			}

			// Validating that data-backup and log-backup don't have PPG defined
			if (strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaDataBackup)) ||
				strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaLogBackup))) &&
				utils.NormalizeNilableString(volume.Properties.ProximityPlacementGroup) != "" {

				errors = append(errors, fmt.Errorf("'%v volume spec type cannot have PPG defined for %v on volume %v'", *volume.Properties.VolumeSpecName, applicationType, *volume.Name))
			}

			// Validating that data, log and shared have PPG defined.
			if (strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaData)) ||
				strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaLog)) ||
				strings.EqualFold(*volume.Properties.VolumeSpecName, string(VolumeSpecNameSapHanaShared))) &&
				utils.NormalizeNilableString(volume.Properties.ProximityPlacementGroup) == "" {

				errors = append(errors, fmt.Errorf("'%v volume spec type must have PPG defined for %v on volume %v'", *volume.Properties.VolumeSpecName, applicationType, *volume.Name))
			}

			// Adding volume spec name to hashmap for post volume loop check
			volumeSpecRepeatCount[*volume.Properties.VolumeSpecName] += 1
		}

		// Validating required volume spec types
		for _, requiredVolumeSpec := range RequiredVolumesForSAPHANA() {
			if _, ok := volumeSpecRepeatCount[requiredVolumeSpec]; !ok {
				errors = append(errors, fmt.Errorf("'required volume spec type %v is not present for %v'", requiredVolumeSpec, applicationType))
			}
		}

		// Validating that volume spec does not repeat
		for volumeSpecName, count := range volumeSpecRepeatCount {
			if count > 1 {
				errors = append(errors, fmt.Errorf("'volume spec type %v cannot be repeated for %v'", volumeSpecName, applicationType))
			}
		}
	}

	return errors
}

func validateNetAppVolumeGroupExportPolicyRule(rule volumegroups.ExportPolicyRule, protocolType string) []error {
	errors := make([]error, 0)

	// Validating that nfsv3 and nfsv4.1 are not enabled in the same rule
	if *rule.Nfsv3 && *rule.Nfsv41 {
		errors = append(errors, fmt.Errorf("'nfsv3 and nfsv4.1 cannot be enabled at the same time'"))
	}

	// Validating that nfsv4.1 export policy is not set on nfsv3 volume
	if *rule.Nfsv41 && strings.EqualFold(protocolType, string(ProtocolTypeNfsV3)) {
		errors = append(errors, fmt.Errorf("'nfsv4.1 export policy cannot be enabled on nfsv3 volume'"))
	}

	// Validating that nfsv3 export policy is not set on nfsv4.1 volume
	if *rule.Nfsv3 && strings.EqualFold(protocolType, string(ProtocolTypeNfsV41)) {
		errors = append(errors, fmt.Errorf("'nfsv3 export policy cannot be enabled on nfsv4.1 volume'"))
	}

	return errors
}

func FindStringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
