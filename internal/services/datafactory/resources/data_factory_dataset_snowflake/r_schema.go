// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_snowflake

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryDatasetSnowflakeSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Dataset Snowflake. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},

		"data_factory_id": {
			Description:  "The Data Factory ID in which to associate the Linked Service with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"linked_service_name": {
			Description:  "The Data Factory Linked Service name in which to associate the Dataset with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"table_name": {
			Description:  "The table name of the Data Factory Dataset Snowflake.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"schema_name": {
			Description:  "The schema name of the Data Factory Dataset Snowflake.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Dataset Snowflake.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Description:  "The description for the Data Factory Dataset Snowflake.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Dataset Snowflake.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"folder": {
			Description:  "The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_properties": {
			Description: "A map of additional properties to associate with the Data Factory Dataset Snowflake.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"schema_column": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the column.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Description: "Type of the column. Please note these values are case sensitive.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"NUMBER",
							"DECIMAL",
							"NUMERIC",
							"INT",
							"INTEGER",
							"BIGINT",
							"SMALLINT",
							"FLOAT",
							"FLOAT4",
							"FLOAT8",
							"DOUBLE",
							"DOUBLE PRECISION",
							"REAL",
							"VARCHAR",
							"CHAR",
							"CHARACTER",
							"STRING",
							"TEXT",
							"BINARY",
							"VARBINARY",
							"BOOLEAN",
							"DATE",
							"DATETIME",
							"TIME",
							"TIMESTAMP",
							"TIMESTAMP_LTZ",
							"TIMESTAMP_NTZ",
							"TIMESTAMP_TZ",
							"VARIANT",
							"OBJECT",
							"ARRAY",
							"GEOGRAPHY",
						}, false),
					},
					"precision": {
						Description:  "The total number of digits allowed.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"scale": {
						Description:  "The number of digits allowed to the right of the decimal point.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},
	}
}
