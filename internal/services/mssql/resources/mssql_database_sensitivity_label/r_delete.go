// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_sensitivity_label

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.SensitivityLabelsClient
			id, err := sensitivitylabels.ParseSensitivityLabelSourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			columnId := sensitivitylabels.NewColumnID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName, id.SchemaName, id.TableName, id.ColumnName)
			if resp, err := client.Delete(ctx, columnId); err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
