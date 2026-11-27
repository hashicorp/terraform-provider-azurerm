// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_spark_pool

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
)

func expandSynapseSparkPoolAutoPauseProperties(input []any) *bigdatapools.AutoPauseProperties {
	if len(input) == 0 {
		return &bigdatapools.AutoPauseProperties{
			Enabled: pointer.To(false),
		}
	}
	v := input[0].(map[string]any)
	return &bigdatapools.AutoPauseProperties{
		DelayInMinutes: pointer.To(int64(v["delay_in_minutes"].(int))),
		Enabled:        pointer.To(true),
	}
}

func expandSynapseSparkPoolAutoScaleProperties(input []any) *bigdatapools.AutoScaleProperties {
	if len(input) == 0 || input[0] == nil {
		return &bigdatapools.AutoScaleProperties{
			Enabled: pointer.To(false),
		}
	}
	v := input[0].(map[string]any)
	return &bigdatapools.AutoScaleProperties{
		MinNodeCount: pointer.To(int64(v["min_node_count"].(int))),
		Enabled:      pointer.To(true),
		MaxNodeCount: pointer.To(int64(v["max_node_count"].(int))),
	}
}

func expandSynapseSparkPoolLibraryRequirements(input []any) *bigdatapools.LibraryRequirements {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &bigdatapools.LibraryRequirements{
		Content:  pointer.To(v["content"].(string)),
		Filename: pointer.To(v["filename"].(string)),
	}
}

func expandSynapseSparkPoolSparkConfig(input []any) *bigdatapools.SparkConfigProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]any)
	return &bigdatapools.SparkConfigProperties{
		Content:  pointer.To(value["content"].(string)),
		Filename: pointer.To(value["filename"].(string)),
	}
}
