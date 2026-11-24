// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_binary

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryDatasetBinarySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Binary Dataset. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
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
			Description:  "The Data Factory Linked Service name in which to associate the Binary Dataset with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Binary Dataset Specific Field
		"http_server_location": {
			Type:          pluginsdk.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"azure_blob_storage_location", "sftp_server_location"},
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
				},
			},
		},

		"sftp_server_location": {
			Type:          pluginsdk.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"azure_blob_storage_location", "http_server_location"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"path": {
						Description:  "The folder path to the file on the SFTP server.",
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
						Description: "The filename of the file on the SFTP server.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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

		// Binary Dataset Specific Field
		"azure_blob_storage_location": {
			Type:          pluginsdk.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"http_server_location", "sftp_server_location"},
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
					"path": {
						Description:  "The folder path to the file in the blob container.",
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
						Description:  "The filename of the file in the blob container.",
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

		"parameters": {
			Description: "A list of parameters to associate with the Data Factory Binary Dataset.",
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
			Description: "List of tags that can be used for describing the Data Factory Binary Dataset.",
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
			Description: "A map of additional properties to associate with the Data Factory Binary Dataset.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"compression": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// TarGZip, GZip, ZipDeflate
					"level": {
						Description: "The level of compression.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Optimal",
							"Fastest",
						}, false),
					},
					// SFTP Specific field
					"type": {
						Description: "The type of compression used during transport.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							helpers.TypeBasicDatasetCompressionTypeBZip2,
							helpers.TypeBasicDatasetCompressionTypeDeflate,
							helpers.TypeBasicDatasetCompressionTypeGZip,
							helpers.TypeBasicDatasetCompressionTypeTar,
							helpers.TypeBasicDatasetCompressionTypeTarGZip,
							helpers.TypeBasicDatasetCompressionTypeZipDeflate,
						}, false),
					},
				},
			},
		},
	}
}
