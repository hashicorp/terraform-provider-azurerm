// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenNetAppVolumeExportPolicyRule(input *volumes.VolumePropertiesExportPolicy) []any {
	results := make([]any, 0)
	if input == nil || input.Rules == nil {
		return results
	}

	for _, item := range *input.Rules {
		allowedClients := []string{}
		if v := item.AllowedClients; v != nil {
			allowedClients = strings.Split(*v, ",")
		}

		protocolsEnabled := []string{}
		if pointer.From(item.Cifs) {
			protocolsEnabled = append(protocolsEnabled, "CIFS")
		}
		if pointer.From(item.Nfsv3) {
			protocolsEnabled = append(protocolsEnabled, "NFSv3")
		}
		if pointer.From(item.Nfsv41) {
			protocolsEnabled = append(protocolsEnabled, "NFSv4.1")
		}

		results = append(results, map[string]any{
			"allowed_clients":                pluginsdk.FlattenSlice(&allowedClients),
			"kerberos_5_read_only_enabled":   pointer.From(item.Kerberos5ReadOnly),
			"kerberos_5_read_write_enabled":  pointer.From(item.Kerberos5ReadWrite),
			"kerberos_5i_read_only_enabled":  pointer.From(item.Kerberos5iReadOnly),
			"kerberos_5i_read_write_enabled": pointer.From(item.Kerberos5iReadWrite),
			"kerberos_5p_read_only_enabled":  pointer.From(item.Kerberos5pReadOnly),
			"kerberos_5p_read_write_enabled": pointer.From(item.Kerberos5pReadWrite),
			"protocol":                       pluginsdk.FlattenSlice(&protocolsEnabled),
			"root_access_enabled":            pointer.From(item.HasRootAccess),
			"rule_index":                     pointer.From(item.RuleIndex),
			"unix_read_only":                 pointer.From(item.UnixReadOnly),
			"unix_read_write":                pointer.From(item.UnixReadWrite),
		})
	}

	return results
}

func flattenNetAppVolumeMountTargets(input *[]volumes.MountTargetProperties) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"ip_address":      pointer.From(item.IPAddress),
			"smb_server_fqdn": pointer.From(item.SmbServerFqdn),
		})
	}

	return results
}

func flattenNetAppVolumeDataProtectionReplication(input *volumes.VolumePropertiesDataProtection) []any {
	if input == nil || input.Replication == nil || input.Replication.EndpointType == nil {
		return []any{}
	}

	if strings.ToLower(string(*input.Replication.EndpointType)) == "" || !strings.EqualFold(string(*input.Replication.EndpointType), "dst") {
		return []any{}
	}

	replicationFrequency := ""
	if input.Replication.ReplicationSchedule != nil {
		replicationFrequency = helpers.TranslateSDKSchedule(strings.ToLower(string(*input.Replication.ReplicationSchedule)))
	}

	return []any{
		map[string]any{
			"endpoint_type":             strings.ToLower(string(*input.Replication.EndpointType)),
			"remote_volume_location":    location.NormalizeNilable(input.Replication.RemoteVolumeRegion),
			"remote_volume_resource_id": input.Replication.RemoteVolumeResourceId,
			"replication_frequency":     replicationFrequency,
		},
	}
}

func flattenNetAppVolumeDataProtectionSnapshotPolicy(input *volumes.VolumePropertiesDataProtection) []any {
	if input == nil || input.Snapshot == nil || input.Snapshot.SnapshotPolicyId == nil || *input.Snapshot.SnapshotPolicyId == "" {
		return []any{}
	}

	return []any{
		map[string]any{
			"snapshot_policy_id": input.Snapshot.SnapshotPolicyId,
		},
	}
}

func flattenNetAppVolumeDataProtectionBackupPolicy(input *volumes.VolumePropertiesDataProtection) []any {
	if input == nil || input.Backup == nil {
		return []any{}
	}

	backupPolicyID := ""
	policyEnforced := false
	backupVaultID := ""

	if input.Backup.BackupPolicyId != nil {
		backupPolicyID = pointer.From(input.Backup.BackupPolicyId)
	}

	if input.Backup.PolicyEnforced != nil {
		policyEnforced = pointer.From(input.Backup.PolicyEnforced)
	}

	if input.Backup.BackupVaultId != nil {
		backupVaultID = pointer.From(input.Backup.BackupVaultId)
	}

	return []any{
		map[string]any{
			"backup_policy_id": backupPolicyID,
			"policy_enabled":   policyEnforced,
			"backup_vault_id":  backupVaultID,
		},
	}
}

func flattenNetAppVolumeDataProtectionAdvancedRansomwareProtection(input *volumes.VolumePropertiesDataProtection) []any {
	if input == nil || input.RansomwareProtection == nil {
		return []any{}
	}

	desiredState := ""
	if input.RansomwareProtection.DesiredRansomwareProtectionState != nil {
		desiredState = pointer.FromEnum(input.RansomwareProtection.DesiredRansomwareProtectionState)
	}

	// Only return the block if a desired state has been set
	if desiredState == "" {
		return []any{}
	}

	protectionEnabled := strings.EqualFold(desiredState, string(volumes.DesiredRansomwareProtectionStateEnabled))

	return []any{
		map[string]any{
			"protection_enabled": protectionEnabled,
		},
	}
}
