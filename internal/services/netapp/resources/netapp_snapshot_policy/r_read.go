// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.SnapshotPoliciesClient

			id, err := snapshotpolicies.ParseSnapshotPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := NetAppSnapshotPolicyResourceModel{
				Name:              id.SnapshotPolicyName,
				ResourceGroupName: id.ResourceGroupName,
				AccountName:       id.NetAppAccountName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				props := model.Properties
				state.Enabled = pointer.From(props.Enabled)
				state.HourlySchedule = flattenNetAppSnapshotPolicyHourlySchedule(props.HourlySchedule)
				state.DailySchedule = flattenNetAppSnapshotPolicyDailySchedule(props.DailySchedule)
				state.WeeklySchedule = flattenNetAppSnapshotPolicyWeeklySchedule(props.WeeklySchedule)
				state.MonthlySchedule = flattenNetAppSnapshotPolicyMonthlySchedule(props.MonthlySchedule)

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func findVolumesUsingSnapshotPolicy(ctx context.Context, client *clients.Client, snapshotPolicyId snapshotpolicies.SnapshotPolicyId) ([]string, error) {
	volumeIds := make([]string, 0)
	poolClient := client.NetApp.PoolClient
	accountId := capacitypools.NewNetAppAccountID(snapshotPolicyId.SubscriptionId, snapshotPolicyId.ResourceGroupName, snapshotPolicyId.NetAppAccountName)

	poolsResult, err := poolClient.PoolsList(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("listing capacity pools in account %s: %+v", snapshotPolicyId.NetAppAccountName, err)
	}

	if model := poolsResult.Model; model != nil {
		volumeClient := client.NetApp.VolumeClient

		for _, pool := range *model {
			if pool.Name == nil {
				continue
			}

			poolNameParts := strings.Split(pointer.From(pool.Name), "/")
			poolName := poolNameParts[len(poolNameParts)-1]
			volumeId := volumes.NewCapacityPoolID(snapshotPolicyId.SubscriptionId, snapshotPolicyId.ResourceGroupName, snapshotPolicyId.NetAppAccountName, poolName)

			volumesResult, err := volumeClient.List(ctx, volumeId)
			if err != nil {
				return nil, fmt.Errorf("listing volumes in pool %s: %+v", poolName, err)
			}

			if volumesModel := volumesResult.Model; volumesModel != nil {
				for _, volume := range *volumesModel {
					if volume.Id == nil || volume.Properties.DataProtection == nil ||
						volume.Properties.DataProtection.Snapshot == nil || volume.Properties.DataProtection.Snapshot.SnapshotPolicyId == nil {
						continue
					}

					if strings.EqualFold(*volume.Properties.DataProtection.Snapshot.SnapshotPolicyId, snapshotPolicyId.ID()) {
						volumeIds = append(volumeIds, *volume.Id)
					}
				}
			}
		}
	}

	return volumeIds, nil
}

func waitForSnapshotPolicyCreation(ctx context.Context, client *snapshotpolicies.SnapshotPoliciesClient, id snapshotpolicies.SnapshotPolicyId, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappSnapshotPolicyStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to complete: %+v", id, err)
	}

	return nil
}

func waitForSnapshotPolicyDeletion(ctx context.Context, client *snapshotpolicies.SnapshotPoliciesClient, id snapshotpolicies.SnapshotPolicyId, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappSnapshotPolicyStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func netappSnapshotPolicyStateRefreshFunc(ctx context.Context, client *snapshotpolicies.SnapshotPoliciesClient, id snapshotpolicies.SnapshotPolicyId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
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
