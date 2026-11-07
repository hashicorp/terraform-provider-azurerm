// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_custom_log

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config WorkspaceTableCustomLogResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient

			workspaceId, err := workspaces.ParseWorkspaceID(config.WorkspaceId)
			if err != nil {
				return err
			}

			id := tables.NewTableID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := tables.Table{
				Properties: &tables.TableProperties{
					Plan:                 pointer.ToEnum[tables.TablePlanEnum](config.Plan),
					RetentionInDays:      helpers.DefaultRetentionInDays,
					TotalRetentionInDays: helpers.DefaultRetentionInDays,
					Schema: &tables.Schema{
						Columns:      expandWorkspaceTableColumns(config.Columns),
						DisplayName:  pointer.To(config.DisplayName),
						Description:  pointer.To(config.Description),
						Name:         pointer.To(config.Name),
						TableSubType: pointer.To(tables.TableSubTypeEnumDataCollectionRuleBased),
						TableType:    pointer.To(tables.TableTypeEnumCustomLog),
					},
				},
			}

			if config.RetentionInDays > 0 {
				param.Properties.RetentionInDays = pointer.To(config.RetentionInDays)
			}

			if config.TotalRetentionInDays > 0 {
				param.Properties.TotalRetentionInDays = pointer.To(config.TotalRetentionInDays)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
