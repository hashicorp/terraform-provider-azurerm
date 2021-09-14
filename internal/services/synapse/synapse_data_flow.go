package synapse

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/sdk/2021-06-01-preview/artifacts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func SchemaForDataFlowSourceAndSink() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"description": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"dataset": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"parameters": {
								Type:     pluginsdk.TypeMap,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
				},

				"linked_service": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"parameters": {
								Type:     pluginsdk.TypeMap,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
				},

				"schema_linked_service": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"parameters": {
								Type:     pluginsdk.TypeMap,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
				},
			},
		},
	}
}

func expandSynapseDataFlowSource(input []interface{}) *[]artifacts.DataFlowSource {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]artifacts.DataFlowSource, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, artifacts.DataFlowSource{
			Description:         utils.String(raw["description"].(string)),
			Name:                utils.String(raw["name"].(string)),
			Dataset:             expandSynapseDatasetReference(raw["dataset"].([]interface{})),
			LinkedService:       expandSynapseLinkedServiceReference(raw["linked_service"].([]interface{})),
			SchemaLinkedService: expandSynapseLinkedServiceReference(raw["schema_linked_service"].([]interface{})),
		})
	}
	return &result
}

func expandSynapseDataFlowSink(input []interface{}) *[]artifacts.DataFlowSink {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]artifacts.DataFlowSink, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, artifacts.DataFlowSink{
			Description:         utils.String(raw["description"].(string)),
			Name:                utils.String(raw["name"].(string)),
			Dataset:             expandSynapseDatasetReference(raw["dataset"].([]interface{})),
			LinkedService:       expandSynapseLinkedServiceReference(raw["linked_service"].([]interface{})),
			SchemaLinkedService: expandSynapseLinkedServiceReference(raw["schema_linked_service"].([]interface{})),
		})
	}
	return &result
}

func expandSynapseDatasetReference(input []interface{}) *artifacts.DatasetReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &artifacts.DatasetReference{
		Type:          utils.String("DatasetReference"),
		ReferenceName: utils.String(raw["name"].(string)),
		Parameters:    raw["parameters"].(map[string]interface{}),
	}
}

func expandSynapseLinkedServiceReference(input []interface{}) *artifacts.LinkedServiceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &artifacts.LinkedServiceReference{
		Type:          utils.String("LinkedServiceReference"),
		ReferenceName: utils.String(raw["name"].(string)),
		Parameters:    raw["parameters"].(map[string]interface{}),
	}
}

func flattenSynapseDataFlowSource(input *[]artifacts.DataFlowSource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		name := ""
		description := ""
		if v.Name != nil {
			name = *v.Name
		}
		if v.Description != nil {
			description = *v.Description
		}
		result = append(result, map[string]interface{}{
			"name":                  name,
			"description":           description,
			"dataset":               flattenSynapseDatasetReference(v.Dataset),
			"linked_service":        flattenSynapseLinkedServiceReference(v.LinkedService),
			"schema_linked_service": flattenSynapseLinkedServiceReference(v.SchemaLinkedService),
		})
	}
	return result
}

func flattenSynapseDataFlowSink(input *[]artifacts.DataFlowSink) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		name := ""
		description := ""
		if v.Name != nil {
			name = *v.Name
		}
		if v.Description != nil {
			description = *v.Description
		}
		result = append(result, map[string]interface{}{
			"name":                  name,
			"description":           description,
			"dataset":               flattenSynapseDatasetReference(v.Dataset),
			"linked_service":        flattenSynapseLinkedServiceReference(v.LinkedService),
			"schema_linked_service": flattenSynapseLinkedServiceReference(v.SchemaLinkedService),
		})
	}
	return result
}

func flattenSynapseDatasetReference(input *artifacts.DatasetReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.ReferenceName != nil {
		name = *input.ReferenceName
	}

	return []interface{}{
		map[string]interface{}{
			"name":       name,
			"parameters": input.Parameters,
		},
	}
}

func flattenSynapseLinkedServiceReference(input *artifacts.LinkedServiceReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.ReferenceName != nil {
		name = *input.ReferenceName
	}

	return []interface{}{
		map[string]interface{}{
			"name":       name,
			"parameters": input.Parameters,
		},
	}
}
