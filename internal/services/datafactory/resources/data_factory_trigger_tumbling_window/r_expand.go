// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_tumbling_window

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func expandDataFactoryTriggerTumblingWindowRetryPolicy(input []any) *datafactory.RetryPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]any)
	return &datafactory.RetryPolicy{
		Count:             pointer.To(int32(raw["count"].(int))),
		IntervalInSeconds: pointer.To(int32(raw["interval"].(int))),
	}
}

func expandDataFactoryTriggerSinglePipeline(input []any) *datafactory.TriggerPipelineReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &datafactory.TriggerPipelineReference{
		PipelineReference: &datafactory.PipelineReference{
			ReferenceName: pointer.To(raw["name"].(string)),
			Type:          pointer.To("PipelineReference"),
		},
		Parameters: raw["parameters"].(map[string]any),
	}
}

func expandDataFactoryTriggerDependency(input []any) *[]datafactory.BasicDependencyReference {
	if len(input) == 0 {
		return nil
	}

	result := make([]datafactory.BasicDependencyReference, 0, len(input))
	for _, item := range input {
		raw := item.(map[string]any)

		var trigger datafactory.BasicDependencyReference

		var offset, size *string
		if v := raw["offset"].(string); v != "" {
			offset = pointer.To(v)
		}
		if v := raw["size"].(string); v != "" {
			size = pointer.To(v)
		}

		if v := raw["trigger_name"].(string); v != "" {
			trigger = &datafactory.TumblingWindowTriggerDependencyReference{
				Offset: offset,
				Size:   size,
				ReferenceTrigger: &datafactory.TriggerReference{
					ReferenceName: pointer.To(v),
					Type:          pointer.To("TriggerReference"),
				},
			}
		} else {
			trigger = &datafactory.SelfDependencyTumblingWindowTriggerReference{
				Offset: offset,
				Size:   size,
			}
		}

		result = append(result, trigger)
	}
	return &result
}
