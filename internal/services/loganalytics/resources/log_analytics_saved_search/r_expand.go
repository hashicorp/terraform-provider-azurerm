// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
)

func expandSavedSearchFunctionParameters(input []any) *string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		}
	}

	if len(result) == 0 {
		return nil
	}

	return pointer.To(strings.Join(result, ", "))
}

func expandSavedSearchTag(input map[string]any) *[]savedsearches.Tag {
	results := make([]savedsearches.Tag, 0)
	for key, value := range input {
		result := savedsearches.Tag{
			Name:  key,
			Value: value.(string),
		}
		results = append(results, result)
	}
	return &results
}
