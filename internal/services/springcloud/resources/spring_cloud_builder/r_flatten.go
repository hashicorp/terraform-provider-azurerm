// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_builder

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func flattenBuildServiceBuilderBuildPacksGroupPropertiesArray(input *[]appplatform.BuildpacksGroupProperties) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"name":           pointer.From(item.Name),
			"build_pack_ids": flattenBuildServiceBuilderBuildPackPropertiesArray(item.Buildpacks),
		})
	}
	return results
}

func flattenBuildServiceBuilderStackProperties(input *appplatform.StackProperties) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"id":      pointer.From(input.ID),
			"version": pointer.From(input.Version),
		},
	}
}

func flattenBuildServiceBuilderBuildPackPropertiesArray(input *[]appplatform.BuildpackProperties) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, pointer.From(item.ID))
	}
	return results
}
