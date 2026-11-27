// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_spark_pool

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
)

func flattenSynapseSparkPoolAutoPauseProperties(input *bigdatapools.AutoPauseProperties) []any {
	if input == nil {
		return make([]any, 0)
	}

	if !pointer.From(input.Enabled) {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"delay_in_minutes": pointer.From(input.DelayInMinutes),
		},
	}
}

func flattenSynapseSparkPoolAutoScaleProperties(input *bigdatapools.AutoScaleProperties) []any {
	if input == nil {
		return make([]any, 0)
	}

	if !pointer.From(input.Enabled) {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"max_node_count": pointer.From(input.MaxNodeCount),
			"min_node_count": pointer.From(input.MinNodeCount),
		},
	}
}

func flattenSynapseSparkPoolLibraryRequirements(input *bigdatapools.LibraryRequirements) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"content":  pointer.From(input.Content),
			"filename": pointer.From(input.Filename),
		},
	}
}

func flattenSynapseSparkPoolSparkConfig(input *bigdatapools.SparkConfigProperties) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"content":  pointer.From(input.Content),
			"filename": pointer.From(input.Filename),
		},
	}
}
