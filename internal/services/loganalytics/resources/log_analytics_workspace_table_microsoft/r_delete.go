// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_microsoft

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/helpers"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceTableMicrosoftResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// We can't delete Microsoft tables, so we'll just set the retention to workspace default
			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					RetentionInDays:      helpers.DefaultRetentionInDays,
					TotalRetentionInDays: helpers.DefaultRetentionInDays,
					Schema: &tables.Schema{
						Name:    pointer.To(id.TableName),
						Columns: &[]tables.Column{},
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
