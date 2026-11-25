// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sql_managed_instance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Linked Service SQL Managed Instance. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},

		"data_factory_id": commonschema.ResourceIDReferenceRequiredForceNew(&factories.FactoryId{}, "Data Factory"),

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Linked Service SQL Managed Instance.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"connection_string": {
			Description:      "The connection string in which to authenticate with the SQL Managed Instance.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ExactlyOneOf:     []string{"connection_string", "key_vault_connection_string"},
			DiffSuppressFunc: helpers.AzureRmDataFactoryLinkedServiceConnectionStringDiff,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "The description for the Data Factory Linked Service SQL Managed Instance.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"integration_runtime_name": {
			Description:  "The integration runtime reference to associate with the Data Factory Linked Service SQL Managed Instance.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"key_vault_connection_string": {
			Description:  "Use this argument to store SQL Managed Instance connection string in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.",
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ExactlyOneOf: []string{"connection_string", "key_vault_connection_string"},
			MaxItems:     1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Description:  "The name of an existing Key Vault Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Description:  "The secret name in Azure Key Vault that stores SQL Managed Instance connection string.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"key_vault_password": {
			Description: "Use this argument to store SQL Managed Instance password in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.",
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
						Description:  "The secret name in Azure Key Vault that stores SQL Managed Instance password.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Linked Service SQL Managed Instance.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"service_principal_id": {
			Description:  "The Service Principal ID in which to authenticate against the Azure SQL Managed Instance.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"service_principal_key", "tenant"},
		},

		"service_principal_key": {
			Description:  "The Service Principal key in which to authenticate against the Azure SQL Managed Instance.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"service_principal_id", "tenant"},
		},

		"tenant": {
			Description:  "The tenant ID or name in which to authenticate against the Azure SQL Managed Instance.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"service_principal_id", "service_principal_key"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
