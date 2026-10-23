// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppVolumeDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := volumes.ParseVolumeID(d.Id())
	if err != nil {
		return err
	}

	netApp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("fetching netapp error: %+v", err)
	}

	// Preventing unintentional severe data loss
	if meta.(*clients.Client).Features.NetApp.PreventVolumeDestruction {
		return fmt.Errorf("deleting NetApp Volume %q (Resource Group %q) is not allowed due to prevent_volume_destruction feature flag", id.VolumeName, id.ResourceGroupName)
	}

	// Handling DataProtection
	if netApp.Model != nil && netApp.Model.Properties.DataProtection != nil {
		// Handling Replication before volume deletion
		if netApp.Model.Properties.DataProtection.Replication != nil {
			dataProtectionReplication := netApp.Model.Properties.DataProtection
			replicaVolumeId, err := volumes.ParseVolumeID(id.ID())
			if err != nil {
				return err
			}
			if dataProtectionReplication.Replication != nil && dataProtectionReplication.Replication.EndpointType != nil && !strings.EqualFold(string(*dataProtectionReplication.Replication.EndpointType), "dst") {
				// This is the case where primary volume started the deletion, in this case, to be consistent we will remove replication from secondary
				replicaVolumeId, err = volumes.ParseVolumeID(pointer.From(dataProtectionReplication.Replication.RemoteVolumeResourceId))
				if err != nil {
					return err
				}
			}

			// Checking replication status before deletion, it needs to be broken before proceeding with deletion
			if res, err := client.ReplicationStatus(ctx, *replicaVolumeId); err == nil {
				// Wait for replication state = "mirrored"
				if model := res.Model; model != nil {
					if model.MirrorState != nil && strings.ToLower(string(*model.MirrorState)) == "uninitialized" {
						if err := helpers.WaitForReplMirrorState(ctx, client, *replicaVolumeId, "mirrored"); err != nil {
							return fmt.Errorf("waiting for replica %s to become 'mirrored': %+v", *replicaVolumeId, err)
						}
					}
				}

				// Breaking replication
				// Can't use VolumesBreakReplicationThenPoll because from time to time the LRO SDK fails,
				// please see Pandora's issue: https://github.com/hashicorp/pandora/issues/4571
				if _, err = client.BreakReplication(ctx, *replicaVolumeId, volumes.BreakReplicationRequest{
					ForceBreakReplication: pointer.To(true),
				}); err != nil {
					return fmt.Errorf("breaking replication for %s: %+v", *replicaVolumeId, err)
				}

				// Waiting for replication be in broken state
				log.Printf("[DEBUG] Waiting for the replication of %s to be in broken state", *replicaVolumeId)
				if err := helpers.WaitForReplMirrorState(ctx, client, *replicaVolumeId, "broken"); err != nil {
					return fmt.Errorf("waiting for the breaking of replication for %s: %+v", *replicaVolumeId, err)
				}
			}

			// Deleting replication and waiting for it to fully complete the operation
			if _, err = client.DeleteReplication(ctx, *replicaVolumeId); err != nil {
				return fmt.Errorf("deleting replicate %s: %+v", *replicaVolumeId, err)
			}

			if err := helpers.WaitForReplicationDeletion(ctx, client, *replicaVolumeId); err != nil {
				return fmt.Errorf("waiting for the replica %s to be deleted: %+v", *replicaVolumeId, err)
			}
		}

		// Handling Backup before volume deletion
		if netApp.Model.Properties.DataProtection.Backup != nil {
			dataProtectionBackup := netApp.Model.Properties.DataProtection

			if dataProtectionBackup.Backup != nil {
				// Checking if initial backup is in progress
				volumeIdFromBackupClient := volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)
				if err = waitForBackupRelationshipStateForDeletion(ctx, client, volumeIdFromBackupClient); err != nil {
					return fmt.Errorf("waiting for of %s: %+v", *id, err)
				}

				// Disabling backup policy first, PolicyEnforced and BackupPolicyId can't be sent together
				disableBackupPolicy := volumes.VolumePatch{
					Properties: &volumes.VolumePatchProperties{
						DataProtection: &volumes.VolumePatchPropertiesDataProtection{
							Backup: &volumes.VolumeBackupProperties{
								PolicyEnforced: pointer.To(false),
							},
						},
					},
				}

				if err = client.UpdateThenPoll(ctx, *id, disableBackupPolicy); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				// Wait for volume to complete update
				if err := helpers.WaitForVolumeCreateOrUpdate(ctx, client, *id); err != nil {
					return err
				}

				// Checking again if backup is in progress
				if err = waitForBackupRelationshipStateForDeletion(ctx, client, volumeIdFromBackupClient); err != nil {
					return fmt.Errorf("waiting for of %s: %+v", *id, err)
				}

				// Removing BackupPolicyId
				backupPolicyIdRemoval := volumes.VolumePatch{
					Properties: &volumes.VolumePatchProperties{
						DataProtection: &volumes.VolumePatchPropertiesDataProtection{
							Backup: &volumes.VolumeBackupProperties{
								BackupPolicyId: pointer.To(""),
							},
						},
					},
				}

				if err = client.UpdateThenPoll(ctx, *id, backupPolicyIdRemoval); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				// Wait for volume to complete update
				if err := helpers.WaitForVolumeCreateOrUpdate(ctx, client, *id); err != nil {
					return err
				}
			}
		}
	}

	// Deleting volume and waiting for it fo fully complete the operation
	if err = client.DeleteThenPoll(ctx, *id, volumes.DeleteOperationOptions{
		ForceDelete: pointer.To(true),
	}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = helpers.WaitForVolumeDeletion(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
