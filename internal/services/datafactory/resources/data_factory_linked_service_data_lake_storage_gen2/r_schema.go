// Copyright IBM Corp. 2014, 2025
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

		"url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"use_managed_identity": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"service_principal_key", "service_principal_id", "storage_account_key", "tenant"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"service_principal_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			RequiredWith:  []string{"service_principal_key", "tenant"},
			ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"service_principal_key": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"service_principal_id", "tenant"},
			ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"storage_account_key": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"service_principal_id", "service_principal_key", "use_managed_identity", "tenant"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"tenant": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"service_principal_id", "service_principal_key"},
			ConflictsWith: []string{"storage_account_key", "use_managed_identity"},
			AtLeastOneOf:  []string{"service_principal_key", "service_principal_id", "tenant", "storage_account_key", "use_managed_identity"},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"integration_runtime_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
