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
				return fmt.Errorf("reading Log Analytics Workspace Table %s: %v", id, err)
			}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					updateInput := tables.Table{
						Properties: &tables.TableProperties{
							Plan: props.Plan,
						},
					}

					if metadata.ResourceData.HasChange("plan") {
						updateInput.Properties.Plan = pointer.ToEnum[tables.TablePlanEnum](state.Plan)
					}

					if state.Plan == string(tables.TablePlanEnumAnalytics) {
						updateInput.Properties.RetentionInDays = existing.Model.Properties.RetentionInDays

						if metadata.ResourceData.HasChange("retention_in_days") {
							updateInput.Properties.RetentionInDays = pointer.To(state.RetentionInDays)
							// `0` is not a valid value for `retention_in_days`, and the service will return HTTP 400
							// to reset it to its default value, we need to pass `-1`
							if state.RetentionInDays == 0 {
								updateInput.Properties.RetentionInDays = pointer.To(int64(-1))
							}
						}
					}

					if metadata.ResourceData.HasChange("total_retention_in_days") {
						updateInput.Properties.TotalRetentionInDays = pointer.To(state.TotalRetentionInDays)
					}

					if err := client.UpdateThenPoll(ctx, *id, updateInput); err != nil {
						return fmt.Errorf("failed to update table: %s: %+v", id.TableName, err)
					}
				}
			}

			return nil
		},
	}
}
