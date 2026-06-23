// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func importDataFactoryLinkedService(expectType datafactory.TypeBasicLinkedService) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.LinkedServiceID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).DataFactory.LinkedServiceClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			return nil, fmt.Errorf("retrieving Data Factory %s: %+v", *id, err)
		}

		byteArr, err := json.Marshal(resp.Properties)
		if err != nil {
			return nil, err
		}

		var m map[string]*json.RawMessage
		if err = json.Unmarshal(byteArr, &m); err != nil {
			return nil, err
		}

		t := ""
		if v, ok := m["type"]; ok && v != nil {
			if err := json.Unmarshal(*v, &t); err != nil {
				return nil, err
			}
			delete(m, "type")
		}

		if datafactory.TypeBasicLinkedService(t) != expectType {
			return nil, fmt.Errorf("data factory linked service has mismatched type, expected: %q, got %q", expectType, t)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}

func expandTypedLinkedServiceParameters(input map[string]interface{}) *map[string]linkedservices.ParameterSpecification {
	if len(input) == 0 {
		return nil
	}

	parameterSpec := make(map[string]linkedservices.ParameterSpecification)
	for key, value := range input {
		parameterSpec[key] = linkedservices.ParameterSpecification{
			Type:         linkedservices.ParameterTypeString,
			DefaultValue: pointer.To(value),
		}
	}

	return &parameterSpec
}

func expandTypedLinkedServiceAnnotations(input []string) *[]interface{} {
	if len(input) == 0 {
		return nil
	}

	annotations := make([]interface{}, len(input))
	for i, v := range input {
		annotations[i] = v
	}

	return &annotations
}

func expandTypedLinkedServiceIntegrationRuntimeName(input string) *linkedservices.IntegrationRuntimeReference {
	if input == "" {
		return nil
	}

	return &linkedservices.IntegrationRuntimeReference{
		Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
		ReferenceName: input,
	}
}

func expandTypedLinkedServiceKeyVaultPassword(input []KeyVaultPasswordConfig) *linkedservices.AzureKeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &linkedservices.AzureKeyVaultSecretReference{
		SecretName: config.SecretName,
		Store: linkedservices.LinkedServiceReference{
			Type:          linkedservices.TypeLinkedServiceReference,
			ReferenceName: config.LinkedServiceName,
		},
	}
}

func flattenTypedLinkedServiceParameters(input *map[string]linkedservices.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})
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

func flattenTypedLinkedServiceAnnotations(input *[]interface{}) []string {
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

func flattenTypedLinkedServiceKeyVaultPassword(input *linkedservices.AzureKeyVaultSecretReference) []KeyVaultPasswordConfig {
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

func flattenTypedLinkedServiceIntegrationRuntimeName(input *linkedservices.IntegrationRuntimeReference) string {
	if input == nil {
		return ""
	}

	return input.ReferenceName
}
