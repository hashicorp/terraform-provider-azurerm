// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_sensitivity_label

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *sensitivitylabels.SensitivityLabelSourceId, model *sensitivitylabels.SensitivityLabel) error {
	state := MssqlDatabaseSensitivityLabelModel{
		ColumnName: id.ColumnName,
		DatabaseId: commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName).ID(),
		SchemaName: id.SchemaName,
		TableName:  id.TableName,
	}

	if properties := model.Properties; properties != nil {
		state.InformationType = pointer.From(properties.InformationType)
		state.InformationTypeId = pointer.From(properties.InformationTypeId)
		state.LabelId = pointer.From(properties.LabelId)
		state.LabelName = pointer.From(properties.LabelName)
		if properties.Rank != nil && *properties.Rank != sensitivitylabels.SensitivityLabelRankNone {
			state.Rank = string(*properties.Rank)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
