// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppSnapshotPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := snapshotpolicies.ParseSnapshotPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] NetApp SnapshotPolicy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.SnapshotPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.NetAppAccountName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties
		d.Set("enabled", props.Enabled)
		if err = d.Set("hourly_schedule", flattenNetAppVolumeSnapshotPolicyHourlySchedule(props.HourlySchedule)); err != nil {
			return fmt.Errorf("setting `hourly_schedule`: %+v", err)
		}
		if err = d.Set("daily_schedule", flattenNetAppVolumeSnapshotPolicyDailySchedule(props.DailySchedule)); err != nil {
			return fmt.Errorf("setting `daily_schedule`: %+v", err)
		}
		if err = d.Set("weekly_schedule", flattenNetAppVolumeSnapshotPolicyWeeklySchedule(props.WeeklySchedule)); err != nil {
			return fmt.Errorf("setting `weekly_schedule`: %+v", err)
		}
		if err = d.Set("monthly_schedule", flattenNetAppVolumeSnapshotPolicyMonthlySchedule(props.MonthlySchedule)); err != nil {
			return fmt.Errorf("setting `monthly_schedule`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
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
