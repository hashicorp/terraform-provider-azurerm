// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

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
		Elem:     ResourceSourceAndSink(),
	}
}

func SchemaForDataFlowletSourceAndSink() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem:     ResourceSourceAndSink(),
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

func ExpandDataFactoryDataFlowSource(input []any) *[]dataflows.DataFlowSource {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]dataflows.DataFlowSource, 0)
	for _, v := range input {
		raw := v.(map[string]any)
		result = append(result, dataflows.DataFlowSource{
			Description:         pointer.To(raw["description"].(string)),
			Name:                raw["name"].(string),
			Dataset:             ExpandDataFactoryDatasetReference(raw["dataset"].([]any)),
			LinkedService:       ExpandDataFactoryLinkedServiceReference(raw["linked_service"].([]any)),
			SchemaLinkedService: ExpandDataFactoryLinkedServiceReference(raw["schema_linked_service"].([]any)),
			Flowlet:             ExpandDataFactoryDataFlowReference(raw["flowlet"].([]any)),
		})
	}
	return &result
}

func ExpandDataFactoryDataFlowSink(input []any) *[]dataflows.DataFlowSink {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]dataflows.DataFlowSink, 0)
	for _, v := range input {
		raw := v.(map[string]any)
		result = append(result, dataflows.DataFlowSink{
			Description:               pointer.To(raw["description"].(string)),
			Name:                      raw["name"].(string),
			Dataset:                   ExpandDataFactoryDatasetReference(raw["dataset"].([]any)),
			LinkedService:             ExpandDataFactoryLinkedServiceReference(raw["linked_service"].([]any)),
			SchemaLinkedService:       ExpandDataFactoryLinkedServiceReference(raw["schema_linked_service"].([]any)),
			RejectedDataLinkedService: ExpandDataFactoryLinkedServiceReference(raw["rejected_linked_service"].([]any)),
			Flowlet:                   ExpandDataFactoryDataFlowReference(raw["flowlet"].([]any)),
		})
	}
	return &result
}

func ExpandDataFactoryDataFlowTransformation(input []any) *[]dataflows.Transformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]dataflows.Transformation, 0)
	for _, v := range input {
		raw := v.(map[string]any)
		result = append(result, dataflows.Transformation{
			Description:   pointer.To(raw["description"].(string)),
			Name:          raw["name"].(string),
			Dataset:       ExpandDataFactoryDatasetReference(raw["dataset"].([]any)),
			LinkedService: ExpandDataFactoryLinkedServiceReference(raw["linked_service"].([]any)),
			Flowlet:       ExpandDataFactoryDataFlowReference(raw["flowlet"].([]any)),
		})
	}
	return &result
}

func ExpandDataFactoryDatasetReference(input []any) *dataflows.DatasetReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &dataflows.DatasetReference{
		Type:          dataflows.DatasetReferenceTypeDatasetReference,
		ReferenceName: raw["name"].(string),
		Parameters:    pointer.To(raw["parameters"].(map[string]any)),
	}
}

func ExpandDataFactoryLinkedServiceReference(input []any) *dataflows.LinkedServiceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &dataflows.LinkedServiceReference{
		Type:          dataflows.TypeLinkedServiceReference,
		ReferenceName: raw["name"].(string),
		Parameters:    pointer.To(raw["parameters"].(map[string]any)),
	}
}

func ExpandDataFactoryDataFlowReference(input []any) *dataflows.DataFlowReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &dataflows.DataFlowReference{
		Type:              dataflows.DataFlowReferenceTypeDataFlowReference,
		ReferenceName:     raw["name"].(string),
		Parameters:        pointer.To(raw["parameters"].(map[string]any)),
		DatasetParameters: pointer.To(raw["dataset_parameters"]),
	}
}

func FlattenDataFactoryDataFlowSource(input *[]dataflows.DataFlowSource) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, v := range *input {
		result = append(result, map[string]any{
			"name":                  v.Name,
			"description":           pointer.From(v.Description),
			"dataset":               FlattenDataFactoryDatasetReference(v.Dataset),
			"linked_service":        FlattenDataFactoryLinkedServiceReference(v.LinkedService),
			"schema_linked_service": FlattenDataFactoryLinkedServiceReference(v.SchemaLinkedService),
			"flowlet":               FlattenDataFactoryDataFlowReference(v.Flowlet),
		})
	}
	return result
}

func FlattenDataFactoryDataFlowSink(input *[]dataflows.DataFlowSink) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, v := range *input {
		result = append(result, map[string]any{
			"name":                    v.Name,
			"description":             pointer.From(v.Description),
			"dataset":                 FlattenDataFactoryDatasetReference(v.Dataset),
			"linked_service":          FlattenDataFactoryLinkedServiceReference(v.LinkedService),
			"rejected_linked_service": FlattenDataFactoryLinkedServiceReference(v.RejectedDataLinkedService),
			"schema_linked_service":   FlattenDataFactoryLinkedServiceReference(v.SchemaLinkedService),
			"flowlet":                 FlattenDataFactoryDataFlowReference(v.Flowlet),
		})
	}
	return result
}

func FlattenDataFactoryDataFlowTransformation(input *[]dataflows.Transformation) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, v := range *input {
		result = append(result, map[string]any{
			"name":           v.Name,
			"description":    pointer.From(v.Description),
			"dataset":        FlattenDataFactoryDatasetReference(v.Dataset),
			"linked_service": FlattenDataFactoryLinkedServiceReference(v.LinkedService),
			"flowlet":        FlattenDataFactoryDataFlowReference(v.Flowlet),
		})
	}
	return result
}

func FlattenDataFactoryDatasetReference(input *dataflows.DatasetReference) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":       input.ReferenceName,
			"parameters": pointer.From(input.Parameters),
		},
	}
}

func FlattenDataFactoryLinkedServiceReference(input *dataflows.LinkedServiceReference) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":       input.ReferenceName,
			"parameters": pointer.From(input.Parameters),
		},
	}
}

func FlattenDataFactoryDataFlowReference(input *dataflows.DataFlowReference) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":               input.ReferenceName,
			"parameters":         pointer.From(input.Parameters),
			"dataset_parameters": pointer.From(input.DatasetParameters),
		},
	}
}

func ResourceSourceAndSink() *pluginsdk.Resource {
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
