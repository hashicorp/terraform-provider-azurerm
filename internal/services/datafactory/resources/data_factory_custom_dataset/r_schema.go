// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_custom_dataset

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryCustomDatasetSchema() map[string]*pluginsdk.Schema {
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

		"linked_service": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"parameters": {
						Description: "A map of parameters to associate with the Data Factory Linked Service.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"type": {
			Description: "The type of dataset that will be associated with Data Factory.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"type_properties_json": {
			Description:      "A JSON object that contains the properties of the Data Factory Dataset. Refer to [datafactory/models.go](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/vendor/github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory/models.go) for the shape of the expected JSON. For example, the JSON object for `AzureBlob`-typed Dataset will be unmarshaled into `AzureBlobDatasetTypeProperties struct`.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			StateFunc:        pluginsdk.NormalizeJson,
			DiffSuppressFunc: helpers.SuppressJsonOrderingDifference,
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

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Dataset.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"schema_json": {
			Description:      "A JSON object that contains the schema of the Data Factory Dataset.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			StateFunc:        pluginsdk.NormalizeJson,
			DiffSuppressFunc: helpers.SuppressJsonOrderingDifference,
		},
	}
}
