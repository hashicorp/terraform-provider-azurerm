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

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForColumnTypeEnum(), false),
		},

		"type_hint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForColumnDataTypeHintEnum(), false),
		},

		"hidden": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"display_by_default": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func expandColumns(columns *[]Column) *[]tables.Column {
	result := make([]tables.Column, 0, len(*columns))
	for _, column := range *columns {
		result = append(result, tables.Column{
			Name:             pointer.To(column.Name),
			DisplayName:      pointer.To(column.DisplayName),
			Description:      pointer.To(column.Description),
			IsHidden:         pointer.To(column.IsHidden),
			IsDefaultDisplay: pointer.To(column.IsDefaultDisplay),
			Type:             pointer.To(tables.ColumnTypeEnum(column.Type)),
			DataTypeHint:     pointer.To(tables.ColumnDataTypeHintEnum(column.TypeHint)),
		})
	}
	return pointer.To(result)
}

func flattenColumns(columns *[]tables.Column) []Column {
	if columns == nil {
		return nil
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
