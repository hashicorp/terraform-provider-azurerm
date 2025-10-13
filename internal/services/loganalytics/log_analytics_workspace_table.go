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

type Column struct {
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

func expandColumns(columns []Column) *[]tables.Column {
	result := make([]tables.Column, 0, len(columns))
	for _, column := range columns {
		expandedColumn := tables.Column{
			Name:             pointer.To(column.Name),
			IsHidden:         pointer.To(column.IsHidden),
			IsDefaultDisplay: pointer.To(column.IsDefaultDisplay),
			Type:             pointer.To(tables.ColumnTypeEnum(column.Type)),
		}
		// NB: leaving this as empty strings will prevent the DCR from being created, seeing the following error:
		// Bad Request({"error":{"code":"InvalidPayload","message":"Data collection rule is invalid","details":[{"code":"InvalidTransform","target":"properties.dataFlows[0]"}]}})
		if column.DisplayName != "" {
			expandedColumn.DisplayName = pointer.To(column.DisplayName)
		}
		if column.Description != "" {
			expandedColumn.Description = pointer.To(column.Description)
		}
		if column.TypeHint != "" {
			expandedColumn.DataTypeHint = pointer.To(tables.ColumnDataTypeHintEnum(column.TypeHint))
		}
		result = append(result, expandedColumn)
	}
	return pointer.To(result)
}

func flattenColumns(columns *[]tables.Column) []Column {
	if columns == nil {
		return []Column{}
	}
	result := make([]Column, 0, len(*columns))
	for _, column := range *columns {
		result = append(result, Column{
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
