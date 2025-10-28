// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var defaultRetentionInDays = pointer.To(int64(-1))

type WorkspaceTableColumn struct {
	Name             string `tfschema:"name"`
	DisplayName      string `tfschema:"display_name"`
	Description      string `tfschema:"description"`
	IsHidden         bool   `tfschema:"hidden"`
	IsDefaultDisplay bool   `tfschema:"display_by_default"`
	Type             string `tfschema:"type"`
	TypeHint         string `tfschema:"type_hint"`
}

func columnSchema() map[string]*pluginsdk.Schema {
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

		"display_by_default": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hidden": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"type_hint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForColumnDataTypeHintEnum(), false),
		},
	}
}

func expandColumns(columns []WorkspaceTableColumn) *[]tables.Column {
	result := make([]tables.Column, 0, len(columns))
	for _, column := range columns {
		columnToAdd := tables.Column{
			Name:             pointer.To(column.Name),
			IsHidden:         pointer.To(column.IsHidden),
			IsDefaultDisplay: pointer.To(column.IsDefaultDisplay),
			Type:             pointer.ToEnum[tables.ColumnTypeEnum](column.Type),
		}
		if column.DisplayName != "" {
			columnToAdd.DisplayName = pointer.To(column.DisplayName)
		}
		if column.Description != "" {
			columnToAdd.Description = pointer.To(column.Description)
		}
		if column.TypeHint != "" {
			columnToAdd.DataTypeHint = pointer.ToEnum[tables.ColumnDataTypeHintEnum](column.TypeHint)
		}
		result = append(result, columnToAdd)
	}
	return pointer.To(result)
}

func flattenColumns(columns *[]tables.Column) []WorkspaceTableColumn {
	if columns == nil {
		return nil
	}
	result := make([]WorkspaceTableColumn, 0, len(*columns))
	for _, column := range *columns {
		result = append(result, WorkspaceTableColumn{
			Name:             pointer.From(column.Name),
			DisplayName:      pointer.From(column.DisplayName),
			Description:      pointer.From(column.Description),
			IsHidden:         pointer.From(column.IsHidden),
			IsDefaultDisplay: pointer.From(column.IsDefaultDisplay),
			Type:             pointer.FromEnum(column.Type),
			TypeHint:         pointer.FromEnum(column.DataTypeHint),
		})
	}
	return result
}
