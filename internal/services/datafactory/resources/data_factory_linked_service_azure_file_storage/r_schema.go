// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_file_storage

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryLinkedServiceAzureFileStorageSchema() map[string]*pluginsdk.Schema {
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

		"connection_string": {
			Description:  "The connection string.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"file_share": {
			Description:  "The name of the file share.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host": {
			Description:  "The Host name of the server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_id": {
			Description:  "The user ID to log in the server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"password": {
			Description:  "The password to log in the server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "The description for the Data Factory Linked Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"integration_runtime_name": {
			Description:  "The integration runtime reference to associate with the Data Factory Linked Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"key_vault_password": {
			Description: "Use this argument to store Azure File Storage password in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Description:  "The name of an existing Key Vault Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Description:  "The secret name in Azure Key Vault that stores Azure File Storage password.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
