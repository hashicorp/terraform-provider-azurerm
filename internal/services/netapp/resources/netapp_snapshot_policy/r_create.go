// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netappClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := netappClient.From(metadata.Client).SnapshotPoliciesClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var model NetAppSnapshotPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := snapshotpolicies.NewSnapshotPolicyID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := snapshotpolicies.SnapshotPolicy{
				Location: location.Normalize(model.Location),
				Name:     pointer.To(id.SnapshotPolicyName),
				Properties: snapshotpolicies.SnapshotPolicyProperties{
					HourlySchedule:  expandNetAppSnapshotPolicyHourlySchedule(model.HourlySchedule),
					DailySchedule:   expandNetAppSnapshotPolicyDailySchedule(model.DailySchedule),
					WeeklySchedule:  expandNetAppSnapshotPolicyWeeklySchedule(model.WeeklySchedule),
					MonthlySchedule: expandNetAppSnapshotPolicyMonthlySchedule(model.MonthlySchedule),
					Enabled:         pointer.To(model.Enabled),
				},
				Tags: pointer.To(model.Tags),
			}

			if _, err := client.Create(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// Waiting for snapshot policy be completely provisioned
			metadata.Logger.Infof("[DEBUG] Waiting for %s to complete", id)
			if err := waitForSnapshotPolicyCreation(ctx, client, id, metadata.ResourceData.Timeout(pluginsdk.TimeoutCreate)); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}
