package datafactory

import (
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
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

func SchemaForDataFlowSourceTransformation() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
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
			},
		},
	}
}

func expandDataFactoryDataFlowSource(input []interface{}) *[]datafactory.DataFlowSource {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]datafactory.DataFlowSource, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, datafactory.DataFlowSource{
			Description:         utils.String(raw["description"].(string)),
			Name:                utils.String(raw["name"].(string)),
			Dataset:             expandDataFactoryDatasetReference(raw["dataset"].([]interface{})),
			LinkedService:       expandDataFactoryLinkedServiceReference(raw["linked_service"].([]interface{})),
			SchemaLinkedService: expandDataFactoryLinkedServiceReference(raw["schema_linked_service"].([]interface{})),
		})
	}
	return &result
}

func expandDataFactoryDataFlowSink(input []interface{}) *[]datafactory.DataFlowSink {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]datafactory.DataFlowSink, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, datafactory.DataFlowSink{
			Description:         utils.String(raw["description"].(string)),
			Name:                utils.String(raw["name"].(string)),
			Dataset:             expandDataFactoryDatasetReference(raw["dataset"].([]interface{})),
			LinkedService:       expandDataFactoryLinkedServiceReference(raw["linked_service"].([]interface{})),
			SchemaLinkedService: expandDataFactoryLinkedServiceReference(raw["schema_linked_service"].([]interface{})),
		})
	}
	return &result
}

func expandDataFactoryDataFlowTransformation(input []interface{}) *[]datafactory.Transformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]datafactory.Transformation, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, datafactory.Transformation{
			Description:   utils.String(raw["description"].(string)),
			Name:          utils.String(raw["name"].(string)),
			Dataset:       expandDataFactoryDatasetReference(raw["dataset"].([]interface{})),
			LinkedService: expandDataFactoryLinkedServiceReference(raw["linked_service"].([]interface{})),
		})
	}
	return &result
}

func expandDataFactoryDatasetReference(input []interface{}) *datafactory.DatasetReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &datafactory.DatasetReference{
		Type:          utils.String("DatasetReference"),
		ReferenceName: utils.String(raw["name"].(string)),
		Parameters:    raw["parameters"].(map[string]interface{}),
	}
}

func expandDataFactoryLinkedServiceReference(input []interface{}) *datafactory.LinkedServiceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &datafactory.LinkedServiceReference{
		Type:          utils.String("LinkedServiceReference"),
		ReferenceName: utils.String(raw["name"].(string)),
		Parameters:    raw["parameters"].(map[string]interface{}),
	}
}

func flattenDataFactoryDataFlowSource(input *[]datafactory.DataFlowSource) []interface{} {
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
			"dataset":               flattenDataFactoryDatasetReference(v.Dataset),
			"linked_service":        flattenDataFactoryLinkedServiceReference(v.LinkedService),
			"schema_linked_service": flattenDataFactoryLinkedServiceReference(v.SchemaLinkedService),
		})
	}
	return result
}

func flattenDataFactoryDataFlowSink(input *[]datafactory.DataFlowSink) []interface{} {
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
			"dataset":               flattenDataFactoryDatasetReference(v.Dataset),
			"linked_service":        flattenDataFactoryLinkedServiceReference(v.LinkedService),
			"schema_linked_service": flattenDataFactoryLinkedServiceReference(v.SchemaLinkedService),
		})
	}
	return result
}

func flattenDataFactoryDataFlowTransformation(input *[]datafactory.Transformation) []interface{} {
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
			"name":           name,
			"description":    description,
			"dataset":        flattenDataFactoryDatasetReference(v.Dataset),
			"linked_service": flattenDataFactoryLinkedServiceReference(v.LinkedService),
		})
	}
	return result
}

func flattenDataFactoryDatasetReference(input *datafactory.DatasetReference) []interface{} {
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

func flattenDataFactoryLinkedServiceReference(input *datafactory.LinkedServiceReference) []interface{} {
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
