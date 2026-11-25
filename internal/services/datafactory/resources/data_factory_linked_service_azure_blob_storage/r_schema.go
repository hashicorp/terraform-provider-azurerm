// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_blob_storage

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryLinkedServiceAzureBlobStorageSchema() map[string]*pluginsdk.Schema {
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
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"connection_string", "connection_string_insecure", "sas_uri", "service_endpoint"},
		},

		"connection_string_insecure": {
			Description:  "The connection string sent insecurely.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"connection_string", "connection_string_insecure", "sas_uri", "service_endpoint"},
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				accountKeyRegex := regexp.MustCompile("AccountKey=[^;]+")

				maskedNew := accountKeyRegex.ReplaceAllString(new, "")
				return (new == d.Get(k).(string)) && (maskedNew == old)
			},
		},

		"storage_kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Storage",
				"StorageV2",
				"BlobStorage",
				"BlockBlobStorage",
			}, false),
		},

		"sas_uri": {
			Description:  "The SAS URI.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"connection_string", "connection_string_insecure", "sas_uri", "service_endpoint"},
		},

		"sas_token_linked_key_vault_key": {
			Description: "Use this argument to store SAS Token in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service. A `sas_uri` is required.",
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
						Description:  "The secret name in Azure Key Vault that stores the SAS token.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"service_principal_linked_key_vault_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Description:  "The name of an existing Key Vault Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Description:  "The secret name in Azure Key Vault that stores the Service Principal key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
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

		"use_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"service_principal_id",
			},
		},

		"service_endpoint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"connection_string", "connection_string_insecure", "sas_uri", "service_endpoint"},
		},

		"service_principal_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			ConflictsWith: []string{
				"use_managed_identity",
			},
		},

		"service_principal_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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
