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

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			// We do not delete the resource here, just set the retention to workspace default value, which is
			// achieved by setting the value to `-1`
			retentionInDays := pointer.To(int64(-1))
			totalRetentionInDays := pointer.To(int64(-1))

			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					RetentionInDays:      retentionInDays,
					TotalRetentionInDays: totalRetentionInDays,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
				return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", id.TableName, id.WorkspaceName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}
