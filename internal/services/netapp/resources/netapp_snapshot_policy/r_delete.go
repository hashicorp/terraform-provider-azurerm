// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.SnapshotPoliciesClient
			volumeClient := metadata.Client.NetApp.VolumeClient

			id, err := snapshotpolicies.ParseSnapshotPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Try to delete the snapshot policy using DeleteThenPoll
			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				// Check if error is about snapshot policy being in use
				if strings.Contains(err.Error(), "SnapshotPolicy is used") {
					// Get all volumes in the account that might be using this snapshot policy
					volumeIds, err := findVolumesUsingSnapshotPolicy(ctx, metadata.Client, *id)
					if err != nil {
						return fmt.Errorf("finding volumes using snapshot policy %s: %+v", *id, err)
					}

					// Disassociate snapshot policy from each volume
					for _, volumeId := range volumeIds {
						volId, err := volumes.ParseVolumeID(volumeId)
						if err != nil {
							return err
						}

						// Update volume to remove snapshot policy
						update := volumes.VolumePatch{
							Properties: &volumes.VolumePatchProperties{
								DataProtection: &volumes.VolumePatchPropertiesDataProtection{
									Snapshot: &volumes.VolumeSnapshotProperties{
										SnapshotPolicyId: pointer.To(""),
									},
								},
							},
						}

						locks.ByID(volumeId)

						if err = volumeClient.UpdateThenPoll(ctx, *volId, update); err != nil {
							locks.UnlockByID(volumeId)
							return fmt.Errorf("removing snapshot policy from volume %s: %+v", *volId, err)
						}

						// Wait for the update to complete
						if err := helpers.WaitForVolumeCreateOrUpdate(ctx, volumeClient, *volId); err != nil {
							locks.UnlockByID(volumeId)
							return fmt.Errorf("waiting for snapshot policy removal from volume %s: %+v", *volId, err)
						}

						locks.UnlockByID(volumeId)
					}

					// Try deleting the snapshot policy again
					if err = client.DeleteThenPoll(ctx, *id); err != nil {
						return fmt.Errorf("deleting %s after volume disassociation: %+v", *id, err)
					}
				} else {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}

			if err := waitForSnapshotPolicyDeletion(ctx, client, *id, metadata.ResourceData.Timeout(pluginsdk.TimeoutDelete)); err != nil {
				return err
			}

			return nil
		},
	}
}
