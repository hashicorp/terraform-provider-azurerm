// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var defaultRetentionInDays = pointer.To(int64(-1))

type workspaceTableColumn struct {
	Name        string `tfschema:"name"`
	DisplayName string `tfschema:"display_name"`
	Description string `tfschema:"description"`
	Type        string `tfschema:"type"`
}

func workspaceTableColumnSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForColumnTypeEnum(), false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func workspaceTableColumnSchemaComputed() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func expandWorkspaceTableColumns(columns []workspaceTableColumn) *[]tables.Column {
	result := make([]tables.Column, 0, len(columns))
	for _, column := range columns {
		columnToAdd := tables.Column{
			Name: pointer.To(column.Name),
			Type: pointer.ToEnum[tables.ColumnTypeEnum](column.Type),
		}
		if column.DisplayName != "" {
			columnToAdd.DisplayName = pointer.To(column.DisplayName)
		}
		if column.Description != "" {
			columnToAdd.Description = pointer.To(column.Description)
		}
		result = append(result, columnToAdd)
	}
	return pointer.To(result)
}

func flattenWorkspaceTableColumns(columns *[]tables.Column) []workspaceTableColumn {
	if columns == nil {
		return nil
	}
	result := make([]workspaceTableColumn, 0, len(*columns))
	for _, column := range *columns {
		result = append(result, workspaceTableColumn{
			Name:        pointer.From(column.Name),
			DisplayName: pointer.From(column.DisplayName),
			Description: pointer.From(column.Description),
			Type:        pointer.FromEnum(column.Type),
		})
	}
	return result
}
