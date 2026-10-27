// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_http

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
)

func expandLogicAppActionHttpHeaders(headersRaw map[string]any) (*map[string]string, error) {
	headers := make(map[string]string)

	for i, v := range headersRaw {
		value, err := tags.TagValueToString(v)
		if err != nil {
			return nil, err
		}

		headers[i] = value
	}

	return &headers, nil
}

func expandLogicAppActionHttpQueries(queriesRaw map[string]any) (*map[string]string, error) {
	// per the issue https://github.com/hashicorp/terraform-provider-azurerm/issues/28429
	// empty map and nil are different on the service side.
	if len(queriesRaw) == 0 {
		return nil, nil
	}

	queries := make(map[string]string)

	for i, v := range queriesRaw {
		value, err := tags.TagValueToString(v)
		if err != nil {
			return nil, err
		}

		queries[i] = value
	}

	return &queries, nil
}
