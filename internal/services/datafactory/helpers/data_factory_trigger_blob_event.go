// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func ExpandDataFactoryTriggerPipeline(input []any) *[]datafactory.TriggerPipelineReference {
	if len(input) == 0 {
		return nil
	}

	result := make([]datafactory.TriggerPipelineReference, 0)
	for _, item := range input {
		raw := item.(map[string]any)

		// issue https://github.com/hashicorp/terraform-plugin-sdk/issues/588
		// once it's resolved, we could remove the check empty logic
		name := raw["name"].(string)
		if name == "" {
			continue
		}

		result = append(result, datafactory.TriggerPipelineReference{
			PipelineReference: &datafactory.PipelineReference{
				ReferenceName: pointer.To(raw["name"].(string)),
				Type:          pointer.To("PipelineReference"),
			},
			Parameters: raw["parameters"].(map[string]any),
		})
	}
	return &result
}

func FlattenDataFactoryTriggerPipeline(input *[]datafactory.TriggerPipelineReference) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, item := range *input {
		name := ""
		if item.PipelineReference != nil && item.PipelineReference.ReferenceName != nil {
			name = *item.PipelineReference.ReferenceName
		}

		result = append(result, map[string]any{
			"name":       name,
			"parameters": item.Parameters,
		})
	}
	return result
}
