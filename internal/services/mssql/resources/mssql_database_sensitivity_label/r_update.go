// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_sensitivity_label

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.SensitivityLabelsClient
			id, err := sensitivitylabels.ParseSensitivityLabelSourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config MssqlDatabaseSensitivityLabelModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` or `properties` was nil", id)
			}

			properties := existing.Model.Properties
			if metadata.ResourceData.HasChange("information_type") {
				properties.InformationType = config.InformationType
			}
			if metadata.ResourceData.HasChange("information_type_id") {
				properties.InformationTypeId = config.InformationTypeId
			}
			if metadata.ResourceData.HasChange("label_id") {
				properties.LabelId = config.LabelId
			}
			if metadata.ResourceData.HasChange("label_name") {
				properties.LabelName = config.LabelName
			}
			if metadata.ResourceData.HasChange("rank") {
				properties.Rank = pointer.ToEnum[sensitivitylabels.SensitivityLabelRank](pointer.From(config.Rank))
			}

			properties.ClientClassificationSource = pointer.To(sensitivitylabels.ClientClassificationSourceNative)
			properties.ColumnName = nil
			properties.IsDisabled = nil
			properties.SchemaName = nil
			properties.TableName = nil

			columnId := sensitivitylabels.NewColumnID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName, id.SchemaName, id.TableName, id.ColumnName)
			if _, err := client.CreateOrUpdate(ctx, columnId, sensitivitylabels.SensitivityLabel{Properties: properties}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
