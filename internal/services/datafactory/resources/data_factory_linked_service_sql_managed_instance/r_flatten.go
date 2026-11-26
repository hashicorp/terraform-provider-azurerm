// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sql_managed_instance

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
)

func flattenLinkedServiceSqlManagedInstanceKeyVaultConnectionString(input any) []KeyVaultConnectionStringConfig {
	if input == nil {
		return []KeyVaultConnectionStringConfig{}
	}

	flattened := helpers.FlattenAzureKeyVaultConnectionString(input.(map[string]any))
	if len(flattened) == 0 {
		return []KeyVaultConnectionStringConfig{}
	}

	configMap := flattened[0].(map[string]any)
	return []KeyVaultConnectionStringConfig{{
		LinkedServiceName: configMap["linked_service_name"].(string),
		SecretName:        configMap["secret_name"].(string),
	}}
}

func flattenLinkedServiceSqlManagedInstanceKeyVaultPassword(input *linkedservices.AzureKeyVaultSecretReference) []KeyVaultPasswordConfig {
	if input == nil {
		return []KeyVaultPasswordConfig{}
	}

	config := KeyVaultPasswordConfig{}

	if secretName, ok := input.SecretName.(string); ok {
		config.SecretName = secretName
	}

	config.LinkedServiceName = input.Store.ReferenceName

	return []KeyVaultPasswordConfig{config}
}

func flattenLinkedServiceSqlManagedInstanceParameters(input *map[string]linkedservices.ParameterSpecification) map[string]any {
	output := make(map[string]any)
	if input == nil {
		return output
	}

	for key, param := range *input {
		if param.DefaultValue != nil {
			if str, ok := pointer.From(param.DefaultValue).(string); ok {
				output[key] = str
			}
		}
	}
	return output
}

func flattenLinkedServiceSqlManagedInstanceAnnotations(input *[]any) []string {
	annotations := make([]string, 0)
	if input == nil {
		return annotations
	}

	for _, annotation := range *input {
		if str, ok := annotation.(string); ok {
			annotations = append(annotations, str)
		}
	}
	return annotations
}
