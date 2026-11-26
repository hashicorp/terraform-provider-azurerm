// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_parquet

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryDatasetParquetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Dataset. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},

		"data_factory_id": {
			Description:  "The Data Factory ID in which to associate the Dataset with.",
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

		"additional_properties": {
			Description: "A map of additional properties to associate with the Data Factory Dataset.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Dataset.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		// Parquet Specific Field, one option for 'location'
		"azure_blob_fs_location": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"azure_blob_fs_location", "azure_blob_storage_location", "http_server_location"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"file_system": {
						Description:  "The container on the Azure Data Lake Storage Account hosting the file.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_file_system_enabled": {
						Description: "Whether the `file_system` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"path": {
						Description:  "The folder path to the file on the Azure Data Lake Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_path_enabled": {
						Description: "Whether the `path` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"filename": {
						Description:  "The filename of the file on the Azure Data Lake Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Description: "Whether the `filename` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		// Parquet Specific Field, one option for 'location'
		"azure_blob_storage_location": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"azure_blob_fs_location", "azure_blob_storage_location", "http_server_location"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container": {
						Description:  "The container on the Azure Blob Storage Account hosting the file.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_container_enabled": {
						Description: "Whether the `container` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"dynamic_path_enabled": {
						Description: "Whether the `path` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"filename": {
						Description:  "The filename of the file on the Azure Blob Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Description: "Whether the `filename` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"path": {
						Description:  "The folder path to the file on the Azure Blob Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"compression_codec": {
			Description: "The compression codec used to read/write text files. These values are case-sensitive.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"bzip2",
				"gzip",
				"deflate",
				"ZipDeflate",
				"TarGzip",
				"Tar",
				"snappy",
				"lz4",
			}, false),
		},

		"compression_level": {
			Description: "The compression level.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Optimal",
				"Fastest",
			}, false),
		},

		"description": {
			Description:  "The description for the Data Factory Dataset.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"folder": {
			Description:  "The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Parquet Specific Field, one option for 'location'
		"http_server_location": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"azure_blob_fs_location", "azure_blob_storage_location", "http_server_location"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"relative_url": {
						Description:  "The base URL to the web server hosting the file.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_path_enabled": {
						Description: "Whether the `path` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"filename": {
						Description:  "The filename of the file on the web server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Description: "Whether the `filename` is using dynamic expression, function or system variables.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"path": {
						Description:  "The folder path to the file on the web server.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Dataset.",
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
						Description: "Type of the column. These values are case sensitive.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Byte",
							"Byte[]",
							"Boolean",
							"Date",
							"DateTime",
							"DateTimeOffset",
							"Decimal",
							"Double",
							"Guid",
							"Int16",
							"Int32",
							"Int64",
							"Single",
							"String",
							"TimeSpan",
						}, false),
					},
					"description": {
						Description:  "The description of the column.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
