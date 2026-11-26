// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandNetAppVolumeDataProtectionSnapshotPolicy(input []any) *volumes.VolumeSnapshotProperties {
	if len(input) == 0 {
		return nil
	}

	snapshotObject := volumes.VolumeSnapshotProperties{}

	snapshotRaw := input[0].(map[string]any)

	if v, ok := snapshotRaw["snapshot_policy_id"]; ok {
		snapshotObject.SnapshotPolicyId = pointer.To(v.(string))
	}

	return &snapshotObject
}

func expandNetAppVolumeDataProtectionBackupPolicy(input []any) *volumes.VolumeBackupProperties {
	if len(input) == 0 || input == nil {
		return nil
	}

	backupPolicyObject := volumes.VolumeBackupProperties{}

	backupRaw := input[0].(map[string]any)

	if v, ok := backupRaw["backup_policy_id"]; ok {
		backupPolicyObject.BackupPolicyId = pointer.To(v.(string))
	}

	if v, ok := backupRaw["policy_enabled"]; ok {
		backupPolicyObject.PolicyEnforced = pointer.To(v.(bool))
	}

	if v, ok := backupRaw["backup_vault_id"]; ok {
		backupPolicyObject.BackupVaultId = pointer.To(v.(string))
	}

	return &backupPolicyObject
}

func expandNetAppVolumeDataProtectionBackupPolicyPatch(input []any) *volumes.VolumeBackupProperties {
	if len(input) == 0 || input == nil {
		return nil
	}

	backupPolicyObject := volumes.VolumeBackupProperties{}

	backupRaw := input[0].(map[string]any)

	if v, ok := backupRaw["backup_policy_id"]; ok {
		backupPolicyObject.BackupPolicyId = pointer.To(v.(string))
	}

	if v, ok := backupRaw["policy_enabled"]; ok {
		backupPolicyObject.PolicyEnforced = pointer.To(v.(bool))
	}

	if v, ok := backupRaw["backup_vault_id"]; ok {
		backupPolicyObject.BackupVaultId = pointer.To(v.(string))
	}

	return &backupPolicyObject
}

func expandNetAppVolumeDataProtectionAdvancedRansomwareProtection(input []any) *volumes.RansomwareProtectionSettings {
	if len(input) == 0 {
		return nil
	}

	arpObject := volumes.RansomwareProtectionSettings{}

	arpRaw := input[0].(map[string]any)

	if v, ok := arpRaw["protection_enabled"]; ok {
		desiredState := volumes.DesiredRansomwareProtectionStateDisabled
		if v.(bool) {
			desiredState = volumes.DesiredRansomwareProtectionStateEnabled
		}
		arpObject.DesiredRansomwareProtectionState = pointer.To(desiredState)
	}

	return &arpObject
}

func expandNetAppVolumeDataProtectionAdvancedRansomwareProtectionPatch(input []any) *volumes.RansomwareProtectionPatchSettings {
	if len(input) == 0 {
		return &volumes.RansomwareProtectionPatchSettings{
			DesiredRansomwareProtectionState: pointer.To(volumes.DesiredRansomwareProtectionStateDisabled),
		}
	}

	arpObject := volumes.RansomwareProtectionPatchSettings{}

	arpRaw := input[0].(map[string]any)

	if v, ok := arpRaw["protection_enabled"]; ok {
		desiredState := volumes.DesiredRansomwareProtectionStateDisabled
		if v.(bool) {
			desiredState = volumes.DesiredRansomwareProtectionStateEnabled
		}
		arpObject.DesiredRansomwareProtectionState = pointer.To(desiredState)
	}

	return &arpObject
}

func waitForBackupRelationshipStateForDeletion(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200"}, // 200 means not in the state we need for backup
		Target:                    []string{"204"}, // 204 means backup is in a state that need (! transitioning)
		Refresh:                   netappVolumeBackupRelationshipStateForDeletionRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to not be in the transferring state", id)
	}

	return nil
}

func netappVolumeBackupRelationshipStateForDeletionRefreshFunc(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.BackupsGetLatestStatus(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving backup relationship status information from %s: %s", id, err)
			}
		}

		response := 200
		if res.Model != nil && res.Model.RelationshipStatus != nil && *res.Model.RelationshipStatus != volumes.VolumeBackupRelationshipStatusTransferring {
			// return 204 if state matches desired state
			response = 204
		}

		return res, strconv.Itoa(response), nil
	}
}
