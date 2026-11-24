// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_custom_service

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryLinkedCustomServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Linked Service. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
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

		"type": {
			Description: "The type of data stores that will be connected to Data Factory. For full list of supported data stores, refer to [Azure Data Factory connector](https://learn.microsoft.com/azure/data-factory/connector-overview). Refer to [Microsoft.DataFactory factories/linkedservices](https://learn.microsoft.com/azure/templates/microsoft.datafactory/factories/linkedservices?pivots=deployment-language-terraform) for the type values syntax.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"type_properties_json": {
			Description:      "A JSON object that contains the properties of the Data Factory Linked Service. Refer to [datafactory/models.go](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/vendor/github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory/models.go) for the shape of the expected JSON. For example, the JSON object for `AzureBlobStorage`-typed Linked Service will be unmarshaled into `AzureBlobStorageLinkedServiceTypeProperties struct`.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			StateFunc:        pluginsdk.NormalizeJson,
			DiffSuppressFunc: helpers.SuppressJsonOrderingDifference,
		},

		"description": {
			Description:  "The description for the Data Factory Linked Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"integration_runtime": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The integration runtime reference to associate with the Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"parameters": {
						Description: "A map of parameters to associate with the integration runtime.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Linked Service.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Linked Service.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"additional_properties": {
			Description: "A map of additional properties to associate with the Data Factory Linked Service.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
