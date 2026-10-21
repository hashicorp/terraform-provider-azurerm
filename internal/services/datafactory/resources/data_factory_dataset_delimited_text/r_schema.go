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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"linked_service_name": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_path_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"filename": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_container_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"dynamic_path_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"filename": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_file_system_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_path_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"filename": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_filename_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		// Delimited Text Specific Field
		"column_delimiter": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  ",",
		},

		// Delimited Text Specific Field
		"row_delimiter": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		// Delimited Text Specific Field
		"quote_character": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  `"`,
		},

		// Delimited Text Specific Field
		"escape_character": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  `\`,
		},

		// Delimited Text Specific Field
		"encoding": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Delimited Text Specific Field
		"first_row_as_header": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// Delimited Text Specific Field
		"null_value": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "",
		},

		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"folder": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"compression_codec": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Optimal",
				"Fastest",
			}, false),
		},
	}
}
