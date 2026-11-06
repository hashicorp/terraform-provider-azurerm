// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_key_vault

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryLinkedServiceKeyVaultSchema() map[string]*pluginsdk.Schema {
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

		"key_vault_base_url_dynamic_expression": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: func(val any, key string) ([]string, []error) {
				var warnings []string
				var errors []error

				w, e := validation.StringIsNotEmpty(val, key)
				warnings = append(warnings, w...)
				errors = append(errors, e...)

				if _, e := commonids.ParseKeyVaultID(val.(string)); e == nil {
					errors = append(errors, fmt.Errorf("`%s` must not be a Key Vault ID, use `key_vault_id` attribute instead", key))
				}

				if _, e := url.ParseRequestURI(val.(string)); e == nil {
					errors = append(errors, fmt.Errorf("`%s` must be a URL with dynamic expression, use `key_vault_id` attribute instead", key))
				}

				return warnings, errors
			},
			ExactlyOneOf: []string{"key_vault_id", "key_vault_base_url_dynamic_expression"},
		},

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			ExactlyOneOf: []string{"key_vault_id", "key_vault_base_url_dynamic_expression"},
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
