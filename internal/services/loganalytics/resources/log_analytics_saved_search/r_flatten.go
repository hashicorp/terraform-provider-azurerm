// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
)

func flattenSavedSearchTag(input *[]savedsearches.Tag) map[string]any {
	results := make(map[string]any)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results[item.Name] = item.Value
	}
	return results
}
