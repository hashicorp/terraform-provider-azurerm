// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_tumbling_window

import (
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func flattenDataFactoryTriggerRetryPolicy(input *datafactory.RetryPolicy) []any {
	if input == nil {
		return []any{}
	}

	count := 0
	// a little wield: after tested, it's of type float64
	if v, ok := input.Count.(float64); ok {
		count = int(v)
	}

	interval := 0
	if input.IntervalInSeconds != nil {
		interval = int(*input.IntervalInSeconds)
	}

	return []any{
		map[string]any{
			"count":    count,
			"interval": interval,
		},
	}
}

func flattenDataFactoryTriggerSinglePipeline(input *datafactory.TriggerPipelineReference) []any {
	if input == nil {
		return []any{}
	}

	name := ""
	if input.PipelineReference != nil && input.PipelineReference.ReferenceName != nil {
		name = *input.PipelineReference.ReferenceName
	}

	return []any{
		map[string]any{
			"name":       name,
			"parameters": input.Parameters,
		},
	}
}

func flattenDataFactoryTriggerDependency(input *[]datafactory.BasicDependencyReference) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, item := range *input {
		var offset, size, triggerName string

		switch item := item.(type) {
		case datafactory.TumblingWindowTriggerDependencyReference:
			if item.Size != nil {
				size = *item.Size
			}
			if item.Offset != nil {
				offset = *item.Offset
			}
			if item.ReferenceTrigger != nil && item.ReferenceTrigger.ReferenceName != nil {
				triggerName = *item.ReferenceTrigger.ReferenceName
			}
		case datafactory.SelfDependencyTumblingWindowTriggerReference:
			if item.Size != nil {
				size = *item.Size
			}
			if item.Offset != nil {
				offset = *item.Offset
			}
		}

		result = append(result, map[string]any{
			"trigger_name": triggerName,
			"offset":       offset,
			"size":         size,
		})
	}
	return result
}
