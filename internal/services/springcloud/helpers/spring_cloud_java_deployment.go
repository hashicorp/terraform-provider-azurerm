// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func ExpandSpringCloudDeploymentEnvironmentVariables(envMap map[string]any) map[string]*string {
	output := make(map[string]*string, len(envMap))

	for k, v := range envMap {
		output[k] = pointer.To(v.(string))
	}

	return output
}

func FlattenSpringCloudDeploymentEnvironmentVariables(envMap map[string]*string) map[string]any {
	output := make(map[string]any, len(envMap))
	for i, v := range envMap {
		if v == nil {
			continue
		}
		output[i] = *v
	}
	return output
}

func FlattenSpringCloudDeploymentResourceRequests(input *appplatform.ResourceRequests) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"cpu":    pointer.From(input.CPU),
			"memory": pointer.From(input.Memory),
		},
	}
}
