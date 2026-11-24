// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_delimited_text

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryDatasetDelimitedTextSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Dataset. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
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

		// Delimited Text Specific Field, one option for 'location'
		"http_server_location": {
			Type:         pluginsdk.TypeList,
			MaxItems:     1,
			Optional:     true,
			ExactlyOneOf: []string{"http_server_location", "azure_blob_storage_location", "azure_blob_fs_location"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"relative_url": {
						Description:  "The base URL to the web server hosting the file.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"path": {
						Description:  "The folder path to the file on the web server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_path_enabled": {
						Description: "Is the `path` using dynamic expression, function or system variables?",
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
						Description: "Is the `filename` using dynamic expression, function or system variables?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		// Delimited Text Specific Field, one option for 'location'
		"azure_blob_storage_location": {
			Type:         pluginsdk.TypeList,
			MaxItems:     1,
			Optional:     true,
			ExactlyOneOf: []string{"http_server_location", "azure_blob_storage_location", "azure_blob_fs_location"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container": {
						Description:  "The container on the Azure Blob Storage Account hosting the file.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_container_enabled": {
						Description: "Is the `container` using dynamic expression, function or system variables?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"path": {
						Description: "The folder path to the file. This can be an empty string.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"dynamic_path_enabled": {
						Description: "Is the `path` using dynamic expression, function or system variables?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"filename": {
						Description:  "The filename of the file.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Description: "Is the `filename` using dynamic expression, function or system variables?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"azure_blob_fs_location": {
			Type:         pluginsdk.TypeList,
			MaxItems:     1,
			Optional:     true,
			ExactlyOneOf: []string{"http_server_location", "azure_blob_storage_location", "azure_blob_fs_location"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"file_system": {
						Description:  "The storage Data Lake gen2 file system on the Azure Blob Storage Account hosting the file.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_file_system_enabled": {
						Description: "Is the `file_system` using dynamic expression, function or system variables?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"path": {
						Description:  "The folder path to the file.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_path_enabled": {
						Description: "Is the `path` using dynamic expression, function or system variables?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"filename": {
						Description:  "The filename of the file.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Description: "Is the `filename` using dynamic expression, function or system variables?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		// Delimited Text Specific Field
		"column_delimiter": {
			Description: "The column delimiter.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     ",",
		},

		// Delimited Text Specific Field
		"row_delimiter": {
			Description: "The row delimiter.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		// Delimited Text Specific Field
		"quote_character": {
			Description: "The quote character.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     `"`,
		},

		// Delimited Text Specific Field
		"escape_character": {
			Description: "The escape character.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     `\`,
		},

		// Delimited Text Specific Field
		"encoding": {
			Description:  "The encoding format for the file.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Delimited Text Specific Field
		"first_row_as_header": {
			Description: "When used as input, treat the first row of data as headers. When used as output, write the headers into the output as the first row of data.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		// Delimited Text Specific Field
		"null_value": {
			Description: "The null value string.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "",
		},

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Dataset.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Description:  "The description for the Data Factory Dataset.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Dataset.",
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
			Description: "A map of additional properties to associate with the Data Factory Dataset.",
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

		"compression_codec": {
			Description: "The compression codec used to read/write text files. Please note these values are case sensitive.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"None",
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
			Description: "The compression ratio for the Data Factory Dataset. Please note these values are case sensitive.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Optimal",
				"Fastest",
			}, false),
		},
	}
}
