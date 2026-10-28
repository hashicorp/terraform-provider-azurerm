// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_custom_log

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient

			var config WorkspaceTableCustomLogResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			props := existing.Model.Properties
			if props.Schema == nil {
				props.Schema = &tables.Schema{}
			}

			// Create / Update requests MUST have a nil value for `StandardColumns` or they will get a 400 response
			props.Schema.StandardColumns = nil

			// If RetentionInDaysAsDefault is true, RetentionInDays still returns the actual value.
			// Since the `Update` reuses the payload from the GET request, that value is included.
			// If this value is then sent to Azure, `RetentionInDaysAsDefault` no longer returns `true`
			// causing a diff on the subsequent read where we then set this value into state.
			if pointer.From(props.RetentionInDaysAsDefault) {
				props.RetentionInDays = helpers.DefaultRetentionInDays
			}

			// The comment above applies to `TotalRetentionInDaysAsDefault` as well.
			if pointer.From(props.TotalRetentionInDaysAsDefault) {
				props.TotalRetentionInDays = helpers.DefaultRetentionInDays
			}

			if metadata.ResourceData.HasChange("column") {
				props.Schema.Columns = expandWorkspaceTableColumns(config.Columns)
			}

			if metadata.ResourceData.HasChange("display_name") {
				props.Schema.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Schema.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("plan") {
				props.Plan = pointer.ToEnum[tables.TablePlanEnum](config.Plan)
			}

			if metadata.ResourceData.HasChange("retention_in_days") {
				props.RetentionInDays = pointer.To(config.RetentionInDays)
				if config.RetentionInDays == 0 {
					props.RetentionInDays = helpers.DefaultRetentionInDays
				}
			}

			if metadata.ResourceData.HasChange("total_retention_in_days") {
				props.TotalRetentionInDays = pointer.To(config.TotalRetentionInDays)
				if config.TotalRetentionInDays == 0 {
					props.TotalRetentionInDays = helpers.DefaultRetentionInDays
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
