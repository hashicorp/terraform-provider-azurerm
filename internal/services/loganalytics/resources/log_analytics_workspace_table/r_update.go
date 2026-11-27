// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			props := existing.Model.Properties
			if props.Schema != nil {
				// Azure rejects StandardColumns in CreateOrUpdate requests.
				props.Schema.StandardColumns = nil
			}

			if pointer.From(props.RetentionInDaysAsDefault) {
				props.RetentionInDays = pointer.To(int64(-1))
			}
			if pointer.From(props.TotalRetentionInDaysAsDefault) {
				props.TotalRetentionInDays = pointer.To(int64(-1))
			}

			if metadata.ResourceData.HasChange("plan") {
				props.Plan = pointer.ToEnum[tables.TablePlanEnum](state.Plan)
			}

			if state.Plan == string(tables.TablePlanEnumAnalytics) {
				if metadata.ResourceData.HasChange("retention_in_days") {
					props.RetentionInDays = pointer.To(state.RetentionInDays)
					// Azure uses -1 to reset retention to the workspace default.
					if state.RetentionInDays == 0 {
						props.RetentionInDays = pointer.To(int64(-1))
					}
				}
			} else {
				// Retention is read-only for Basic tables.
				props.RetentionInDays = nil
			}

			if metadata.ResourceData.HasChange("total_retention_in_days") {
				props.TotalRetentionInDays = pointer.To(state.TotalRetentionInDays)
				// Azure uses -1 to reset total retention to the table retention.
				if state.TotalRetentionInDays == 0 {
					props.TotalRetentionInDays = pointer.To(int64(-1))
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
