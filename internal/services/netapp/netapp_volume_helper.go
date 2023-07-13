// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

func expandNetAppVolumeGroupVolumeExportPolicyRule(input []ExportPolicyRule) *volumegroups.VolumePropertiesExportPolicy {

	if len(input) == 0 {
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
	if len(input) == 0 {
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
	if len(input) == 0 {
		return &volumegroups.VolumePropertiesDataProtection{}
	}

	snapshotObject := volumegroups.VolumeSnapshotProperties{}
	snapshotObject.SnapshotPolicyId = &input[0].DataProtectionSnapshotPolicy

	return &volumegroups.VolumePropertiesDataProtection{
		Snapshot: &snapshotObject,
	}
}

func expandNetAppVolumeGroupVolumes(input []NetAppVolumeGroupVolume) (*[]volumegroups.VolumeGroupVolumeProperties, error) {
	if len(input) == 0 {
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
	if len(input) == 0 {
		return &volumes.VolumePatchPropertiesExportPolicy{}
	}

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
	if len(input) == 0 {
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
	if len(input) == 0 {
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
	if len(input) == 0 {
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

	if input == nil || len(pointer.From(input)) == 0 {
		return results, fmt.Errorf("received empty volumegroups.VolumeGroupVolumeProperties slice")
	}

	for _, item := range *input {
		volumeGroupVolume := NetAppVolumeGroupVolume{}

		props := item.Properties
		volumeGroupVolume.Name = getUserDefinedVolumeName(item.Name)
		volumeGroupVolume.VolumePath = props.CreationToken
		volumeGroupVolume.ServiceLevel = string(pointer.From(props.ServiceLevel))
		volumeGroupVolume.SubnetId = props.SubnetId
		volumeGroupVolume.CapacityPoolId = utils.NormalizeNilableString(props.CapacityPoolResourceId)
		volumeGroupVolume.Protocols = pointer.From(props.ProtocolTypes)
		volumeGroupVolume.SecurityStyle = string(pointer.From(props.SecurityStyle))
		volumeGroupVolume.SnapshotDirectoryVisible = pointer.From(props.SnapshotDirectoryVisible)
		volumeGroupVolume.ThroughputInMibps = pointer.From(props.ThroughputMibps)
		volumeGroupVolume.Tags = pointer.From(item.Tags)
		volumeGroupVolume.ProximityPlacementGroupId = utils.NormalizeNilableString(props.ProximityPlacementGroup)
		volumeGroupVolume.VolumeSpecName = pointer.From(props.VolumeSpecName)

		if props.UsageThreshold > 0 {
			usageThreshold := props.UsageThreshold / 1073741824
			volumeGroupVolume.StorageQuotaInGB = usageThreshold
		}

		if props.ExportPolicy != nil && props.ExportPolicy.Rules != nil && len(pointer.From(props.ExportPolicy.Rules)) > 0 {
			volumeGroupVolume.ExportPolicy = flattenNetAppVolumeGroupVolumesExportPolicies(props.ExportPolicy.Rules)
		}

		if props.MountTargets != nil && len(pointer.From(props.MountTargets)) > 0 {
			volumeGroupVolume.MountIpAddresses = flattenNetAppVolumeGroupVolumesMountIpAddresses(props.MountTargets)
		}

		// Getting volume resource directly from standalone volume
		// since VolumeGroup Volumes don't return DataProtection information
		volumeClient := metadata.Client.NetApp.VolumeClient
		id, err := volumes.ParseVolumeID(pointer.From(item.Id))
		if err != nil {
			return []NetAppVolumeGroupVolume{}, err
		}

		standaloneVol, err := volumeClient.Get(ctx, pointer.From(id))
		if err != nil {
			return []NetAppVolumeGroupVolume{}, fmt.Errorf("retrieving %s: %v", id, err)
		}

		if standaloneVol.Model.Properties.DataProtection != nil && standaloneVol.Model.Properties.DataProtection.Replication != nil {
			volumeGroupVolume.DataProtectionReplication = flattenNetAppVolumeGroupVolumesDPReplication(standaloneVol.Model.Properties.DataProtection.Replication)
		}

		if standaloneVol.Model.Properties.DataProtection != nil && standaloneVol.Model.Properties.DataProtection.Snapshot != nil {
			volumeGroupVolume.DataProtectionSnapshotPolicy = flattenNetAppVolumeGroupVolumesDPSnapshotPolicy(standaloneVol.Model.Properties.DataProtection.Snapshot)
		}

		volumeGroupVolume.Id = pointer.From(standaloneVol.Model.Id)

		results = append(results, volumeGroupVolume)
	}

	return results, nil
}

func flattenNetAppVolumeGroupVolumesExportPolicies(input *[]volumegroups.ExportPolicyRule) []ExportPolicyRule {
	results := make([]ExportPolicyRule, 0)

	if input == nil || len(pointer.From(input)) == 0 {
		return results
	}

	for _, item := range pointer.From(input) {
		rule := ExportPolicyRule{}

		rule.RuleIndex = int(pointer.From(item.RuleIndex))
		rule.AllowedClients = pointer.From(item.AllowedClients)
		rule.Nfsv3Enabled = pointer.From(item.Nfsv3)
		rule.Nfsv41Enabled = pointer.From(item.Nfsv41)
		rule.UnixReadOnly = pointer.From(item.UnixReadOnly)
		rule.UnixReadWrite = pointer.From(item.UnixReadWrite)
		rule.RootAccessEnabled = pointer.From(item.HasRootAccess)

		results = append(results, rule)
	}

	return results
}

func flattenNetAppVolumeGroupVolumesMountIpAddresses(input *[]volumegroups.MountTargetProperties) []string {
	results := make([]string, 0)

	if input == nil || len(pointer.From(input)) == 0 {
		return results
	}

	for _, item := range pointer.From(input) {
		if item.IPAddress != nil {
			results = append(results, pointer.From(item.IPAddress))
		}
	}

	return results
}

func flattenNetAppVolumeGroupVolumesDPReplication(input *volumes.ReplicationObject) []DataProtectionReplication {
	if input == nil {
		return []DataProtectionReplication{}
	}
	if string(pointer.From(input.EndpointType)) == "" || !strings.EqualFold(string(pointer.From(input.EndpointType)), string(volumes.EndpointTypeDst)) {
		return []DataProtectionReplication{}
	}

	replicationFrequency := ""
	if input.ReplicationSchedule != nil {
		replicationFrequency = translateSDKSchedule(strings.ToLower(string(pointer.From(input.ReplicationSchedule))))
	}

	return []DataProtectionReplication{
		{
			EndpointType:           strings.ToLower(string(pointer.From(input.EndpointType))),
			RemoteVolumeLocation:   pointer.From(input.RemoteVolumeRegion),
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
			DataProtectionSnapshotPolicy: pointer.From(input.SnapshotPolicyId),
		},
	}
}

func getUserDefinedVolumeName(input *string) string {
	volumeName := pointer.From(input)

	if volumeName == "" {
		return ""
	}

	segments := len(strings.Split(volumeName, "/"))

	return strings.Split(volumeName, "/")[segments-1]
}

func deleteVolume(ctx context.Context, metadata sdk.ResourceMetaData, volumeId string) error {
	client := metadata.Client.NetApp.VolumeClient

	id, err := volumes.ParseVolumeID(volumeId)
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, pointer.From(id))
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
		if dataProtectionReplication.Replication.EndpointType != nil && !strings.EqualFold(string(pointer.From(dataProtectionReplication.Replication.EndpointType)), string(volumes.EndpointTypeDst)) {
			// This is the case where primary volume started the deletion, in this case, to be consistent we will remove replication from secondary
			replicaVolumeId, err = volumesreplication.ParseVolumeID(dataProtectionReplication.Replication.RemoteVolumeResourceId)
			if err != nil {
				return err
			}
		}

		replicationClient := metadata.Client.NetApp.VolumeReplicationClient
		// Checking replication status before deletion, it need to be broken before proceeding with deletion
		if res, err := replicationClient.VolumesReplicationStatus(ctx, pointer.From(replicaVolumeId)); err == nil {
			// Wait for replication state = "mirrored"
			if model := res.Model; model != nil {
				if model.MirrorState != nil && strings.ToLower(string(pointer.From(model.MirrorState))) == "uninitialized" {
					if err := waitForReplMirrorState(ctx, replicationClient, pointer.From(replicaVolumeId), "mirrored"); err != nil {
						return fmt.Errorf("waiting for replica %s to become 'mirrored': %+v", pointer.From(replicaVolumeId), err)
					}
				}
			}

			// Breaking replication
			if err = replicationClient.VolumesBreakReplicationThenPoll(ctx, pointer.From(replicaVolumeId), volumesreplication.BreakReplicationRequest{
				ForceBreakReplication: utils.Bool(true),
			}); err != nil {
				return fmt.Errorf("breaking replication for %s: %+v", pointer.From(replicaVolumeId), err)
			}

			// Waiting for replication be in broken state
			metadata.Logger.Infof("waiting for the replication of %s to be in broken state", pointer.From(replicaVolumeId))
			if err := waitForReplMirrorState(ctx, replicationClient, pointer.From(replicaVolumeId), "broken"); err != nil {
				return fmt.Errorf("waiting for the breaking of replication for %s: %+v", pointer.From(replicaVolumeId), err)
			}
		}

		// Deleting replication and waiting for it to fully complete the operation
		if err = replicationClient.VolumesDeleteReplicationThenPoll(ctx, pointer.From(replicaVolumeId)); err != nil {
			return fmt.Errorf("deleting replicate %s: %+v", pointer.From(replicaVolumeId), err)
		}

		if err := waitForReplicationDeletion(ctx, replicationClient, pointer.From(replicaVolumeId)); err != nil {
			return fmt.Errorf("waiting for the replica %s to be deleted: %+v", pointer.From(replicaVolumeId), err)
		}
	}

	// Deleting volume and waiting for it fo fully complete the operation
	if err = client.DeleteThenPoll(ctx, pointer.From(id), volumes.DeleteOperationOptions{
		ForceDelete: utils.Bool(true),
	}); err != nil {
		return fmt.Errorf("deleting %s: %+v", pointer.From(id), err)
	}

	if err = waitForVolumeDeletion(ctx, client, pointer.From(id)); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", pointer.From(id), err)
	}

	return nil
}

func waitForVolumeCreateOrUpdate(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
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
		return fmt.Errorf("internal-error: context had no deadline")
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
		return fmt.Errorf("internal-error: context had no deadline")
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
		return fmt.Errorf("internal-error: context had no deadline")
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
		return fmt.Errorf("internal-error: context had no deadline")
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
		return fmt.Errorf("internal-error: context had no deadline")
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
