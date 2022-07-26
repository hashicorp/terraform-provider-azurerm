package netapp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumesreplication"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupVolume struct {
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
	RuleIndex           int    `tfschema:"rule_index"`
	AllowedClients      string `tfschema:"allowed_clients"`
	CifsEnabled         bool   `tfschema:"cifs_enabled"`
	Nfsv3Enabled        bool   `tfschema:"nfsv3_enabled"`
	Nfsv41Enabled       bool   `tfschema:"nfsv41_enabled"`
	UnixReadOnly        bool   `tfschema:"unix_read_only"`
	UnixReadWrite       bool   `tfschema:"unix_read_write"`
	RootAccessEnabled   bool   `tfschema:"root_access_enabled"`
	Kerberos5ReadOnly   bool   `tfschema:"kerberos5_read_only"`
	Kerberos5ReadWrite  bool   `tfschema:"kerberos5_read_write"`
	Kerberos5iReadOnly  bool   `tfschema:"kerberos5i_read_only"`
	Kerberos5iReadWrite bool   `tfschema:"kerberos5i_read_write"`
	Kerberos5pReadOnly  bool   `tfschema:"kerberos5p_read_only"`
	Kerberos5pReadWrite bool   `tfschema:"kerberos5p_read_write"`
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

func expandNetAppVolumeGroupExportPolicyRule(input []ExportPolicyRule) *volumegroups.VolumePropertiesExportPolicy {

	if len(input) == 0 || input == nil {
		return &volumegroups.VolumePropertiesExportPolicy{}
	}

	results := make([]volumegroups.ExportPolicyRule, 0)

	for _, item := range input {
		result := volumegroups.ExportPolicyRule{
			AllowedClients:      utils.String(item.AllowedClients),
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
		snapshotDirectoryVisible := item.SnapshotDirectoryVisible

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
				ThroughputMibps:          utils.Float(float64(item.ThroughputInMibps)),
				ProximityPlacementGroup:  utils.String(item.ProximityPlacementGroupId),
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

func flattenNetAppVolumeGroupVolumes(input *[]volumegroups.VolumeGroupVolumeProperties) ([]NetAppVolumeGroupVolume, error) {
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
		rule.AllowedClients = *item.AllowedClients
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
	dataProtectionReplication := existing.Model.Properties.DataProtection

	if dataProtectionReplication != nil && dataProtectionReplication.Replication != nil {
		replicaVolumeId, err := volumesreplication.ParseVolumeID(id.ID())
		if err != nil {
			return err
		}
		if dataProtectionReplication.Replication.EndpointType != nil && strings.ToLower(string(*dataProtectionReplication.Replication.EndpointType)) != "dst" {
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

func netappVolumeReplicationStateRefreshFunc(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.VolumesReplicationStatus(ctx, id)
		if err != nil {
			if httpResponse := res.HttpResponse; httpResponse != nil {
				if httpResponse.StatusCode == 400 && (strings.Contains(strings.ToLower(err.Error()), "deleting") || strings.Contains(strings.ToLower(err.Error()), "volume replication missing or deleted")) {
					// This error can be ignored until a bug is fixed on RP side that it is returning 400 while the replication is in "Deleting" process
					// TODO: remove this workaround when above bug is fixed
				} else if !response.WasNotFound(httpResponse) {
					return nil, "", fmt.Errorf("retrieving replication status from %s: %s", id, err)
				}
			}
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
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

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}
