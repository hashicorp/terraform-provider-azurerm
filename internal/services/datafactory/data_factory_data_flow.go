// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaForDataFlowSourceAndSink() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem:     resourceSourceAndSink(),
	}
}

func SchemaForDataFlowletSourceAndSink() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem:     resourceSourceAndSink(),
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

				"flowlet": {
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

							"dataset_parameters": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
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

func expandDataFactoryDataFlowSource(input []interface{}) *[]dataflows.DataFlowSource {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]dataflows.DataFlowSource, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, dataflows.DataFlowSource{
			Description:         pointer.To(raw["description"].(string)),
			Name:                raw["name"].(string),
			Dataset:             expandDataFactoryDatasetReference(raw["dataset"].([]interface{})),
			LinkedService:       expandDataFactoryLinkedServiceReference(raw["linked_service"].([]interface{})),
			SchemaLinkedService: expandDataFactoryLinkedServiceReference(raw["schema_linked_service"].([]interface{})),
			Flowlet:             expandDataFactoryDataFlowReference(raw["flowlet"].([]interface{})),
		})
	}
	return &result
}

func expandDataFactoryDataFlowSink(input []interface{}) *[]dataflows.DataFlowSink {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]dataflows.DataFlowSink, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, dataflows.DataFlowSink{
			Description:               pointer.To(raw["description"].(string)),
			Name:                      raw["name"].(string),
			Dataset:                   expandDataFactoryDatasetReference(raw["dataset"].([]interface{})),
			LinkedService:             expandDataFactoryLinkedServiceReference(raw["linked_service"].([]interface{})),
			SchemaLinkedService:       expandDataFactoryLinkedServiceReference(raw["schema_linked_service"].([]interface{})),
			RejectedDataLinkedService: expandDataFactoryLinkedServiceReference(raw["rejected_linked_service"].([]interface{})),
			Flowlet:                   expandDataFactoryDataFlowReference(raw["flowlet"].([]interface{})),
		})
	}
	return &result
}

func expandDataFactoryDataFlowTransformation(input []interface{}) *[]dataflows.Transformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]dataflows.Transformation, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, dataflows.Transformation{
			Description:   pointer.To(raw["description"].(string)),
			Name:          raw["name"].(string),
			Dataset:       expandDataFactoryDatasetReference(raw["dataset"].([]interface{})),
			LinkedService: expandDataFactoryLinkedServiceReference(raw["linked_service"].([]interface{})),
			Flowlet:       expandDataFactoryDataFlowReference(raw["flowlet"].([]interface{})),
		})
	}
	return &result
}

func expandDataFactoryDatasetReference(input []interface{}) *dataflows.DatasetReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &dataflows.DatasetReference{
		Type:          dataflows.DatasetReferenceTypeDatasetReference,
		ReferenceName: raw["name"].(string),
		Parameters:    pointer.To(raw["parameters"].(map[string]interface{})),
	}
}

func expandDataFactoryLinkedServiceReference(input []interface{}) *dataflows.LinkedServiceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &dataflows.LinkedServiceReference{
		Type:          dataflows.TypeLinkedServiceReference,
		ReferenceName: raw["name"].(string),
		Parameters:    pointer.To(raw["parameters"].(map[string]interface{})),
	}
}

func expandDataFactoryDataFlowReference(input []interface{}) *dataflows.DataFlowReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &dataflows.DataFlowReference{
		Type:              dataflows.DataFlowReferenceTypeDataFlowReference,
		ReferenceName:     raw["name"].(string),
		Parameters:        pointer.To(raw["parameters"].(map[string]interface{})),
		DatasetParameters: pointer.To(raw["dataset_parameters"]),
	}
}

func flattenDataFactoryDataFlowSource(input *[]dataflows.DataFlowSource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		result = append(result, map[string]interface{}{
			"name":                  v.Name,
			"description":           pointer.From(v.Description),
			"dataset":               flattenDataFactoryDatasetReference(v.Dataset),
			"linked_service":        flattenDataFactoryLinkedServiceReference(v.LinkedService),
			"schema_linked_service": flattenDataFactoryLinkedServiceReference(v.SchemaLinkedService),
			"flowlet":               flattenDataFactoryDataFlowReference(v.Flowlet),
		})
	}
	return result
}

func flattenDataFactoryDataFlowSink(input *[]dataflows.DataFlowSink) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		result = append(result, map[string]interface{}{
			"name":                    v.Name,
			"description":             pointer.From(v.Description),
			"dataset":                 flattenDataFactoryDatasetReference(v.Dataset),
			"linked_service":          flattenDataFactoryLinkedServiceReference(v.LinkedService),
			"rejected_linked_service": flattenDataFactoryLinkedServiceReference(v.RejectedDataLinkedService),
			"schema_linked_service":   flattenDataFactoryLinkedServiceReference(v.SchemaLinkedService),
			"flowlet":                 flattenDataFactoryDataFlowReference(v.Flowlet),
		})
	}
	return result
}

func flattenDataFactoryDataFlowTransformation(input *[]dataflows.Transformation) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		result = append(result, map[string]interface{}{
			"name":           v.Name,
			"description":    pointer.From(v.Description),
			"dataset":        flattenDataFactoryDatasetReference(v.Dataset),
			"linked_service": flattenDataFactoryLinkedServiceReference(v.LinkedService),
			"flowlet":        flattenDataFactoryDataFlowReference(v.Flowlet),
		})
	}
	return result
}

func flattenDataFactoryDatasetReference(input *dataflows.DatasetReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"name":       input.ReferenceName,
			"parameters": pointer.From(input.Parameters),
		},
	}
}

func flattenDataFactoryLinkedServiceReference(input *dataflows.LinkedServiceReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"name":       input.ReferenceName,
			"parameters": pointer.From(input.Parameters),
		},
	}
}

func flattenDataFactoryDataFlowReference(input *dataflows.DataFlowReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"name":               input.ReferenceName,
			"parameters":         pointer.From(input.Parameters),
			"dataset_parameters": pointer.From(input.DatasetParameters),
		},
	}
}

func resourceSourceAndSink() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

			"flowlet": {
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

						"dataset_parameters": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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

			"rejected_linked_service": {
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
	}
}
