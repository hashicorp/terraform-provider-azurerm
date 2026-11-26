// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_data_lake_storage_gen2

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryLinkedServiceDataLakeStorageGen2Schema() map[string]*pluginsdk.Schema {
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

		"url": {
			Description:  "The endpoint for the Azure Data Lake Storage Gen2 service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"use_managed_identity": {
			Description:   "Whether to use the Data Factory's Managed Identity to authenticate against the Azure Data Lake Storage Gen2 account. Incompatible with `service_principal_id`, `service_principal_key`, `tenant` and `storage_account_key`.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"service_principal_key", "service_principal_id", "storage_account_key", "tenant"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"service_principal_id": {
			Description:   "The Service Principal ID with which to authenticate against the Azure Data Lake Storage Gen2 account. Incompatible with `storage_account_key` and `use_managed_identity`.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			RequiredWith:  []string{"service_principal_key", "tenant"},
			ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"service_principal_key": {
			Description:   "The Service Principal key with which to authenticate against the Azure Data Lake Storage Gen2 account.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"service_principal_id", "tenant"},
			ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"storage_account_key": {
			Description:   "The Storage Account Key with which to authenticate against the Azure Data Lake Storage Gen2 account. Incompatible with `service_principal_id`, `service_principal_key`, `tenant` and `use_managed_identity`.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"service_principal_id", "service_principal_key", "use_managed_identity", "tenant"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"tenant": {
			Description:   "The tenant ID or name in which the Service Principal exists to authenticate against the Azure Data Lake Storage Gen2 account.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"service_principal_id", "service_principal_key"},
			ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
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
