// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_builder

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func expandBuildServiceBuilderBuildPacksGroupPropertiesArray(input []any) *[]appplatform.BuildpacksGroupProperties {
	results := make([]appplatform.BuildpacksGroupProperties, 0)
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, appplatform.BuildpacksGroupProperties{
			Name:       pointer.To(v["name"].(string)),
			Buildpacks: expandBuildServiceBuilderBuildPackPropertiesArray(v["build_pack_ids"].([]any)),
		})
	}
	return &results
}

func expandBuildServiceBuilderStackProperties(input []any) *appplatform.StackProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &appplatform.StackProperties{
		ID:      pointer.To(v["id"].(string)),
		Version: pointer.To(v["version"].(string)),
	}
}

func expandBuildServiceBuilderBuildPackPropertiesArray(input []any) *[]appplatform.BuildpackProperties {
	results := make([]appplatform.BuildpackProperties, 0)
	for _, item := range input {
		results = append(results, appplatform.BuildpackProperties{
			ID: pointer.To(item.(string)),
		})
	}
	return &results
}
