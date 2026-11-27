// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_sensitivity_label

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.SensitivityLabelsClient

			var config MssqlDatabaseSensitivityLabelModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			databaseId, err := commonids.ParseSqlDatabaseID(config.DatabaseId)
			if err != nil {
				return err
			}
			// `SensitivityLabelSourceRecommended` is recommended by the service instead of a "real" resource can be managed by Terraform,
			// Thus hard-code the source to `SensitivityLabelSourceCurrent` for the resource ID.
			id := sensitivitylabels.NewSensitivityLabelSourceID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.ServerName, databaseId.DatabaseName, config.SchemaName, config.TableName, config.ColumnName, sensitivitylabels.SensitivityLabelSourceCurrent)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			rank := sensitivitylabels.SensitivityLabelRankNone
			if config.Rank != "" {
				rank = sensitivitylabels.SensitivityLabelRank(config.Rank)
			}

			columnId := sensitivitylabels.NewColumnID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName, id.SchemaName, id.TableName, id.ColumnName)
			payload := sensitivitylabels.SensitivityLabel{
				Properties: &sensitivitylabels.SensitivityLabelProperties{
					ClientClassificationSource: pointer.To(sensitivitylabels.ClientClassificationSourceNative),
					InformationType:            pointer.ToOrNil(config.InformationType),
					InformationTypeId:          pointer.ToOrNil(config.InformationTypeId),
					LabelId:                    pointer.ToOrNil(config.LabelId),
					LabelName:                  pointer.ToOrNil(config.LabelName),
					Rank:                       pointer.To(rank),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, columnId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
