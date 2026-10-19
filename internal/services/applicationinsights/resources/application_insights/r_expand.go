// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"encoding/json"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
)

// Portal JS source-map blob storage is not an ARM property; it is stored as a
// special hidden-link tag. Community-confirmed shape (#13255):
//
//	"hidden-link:Insights.Sourcemap.Storage" = jsonencode({ Uri = "<blob container url>" })
const appInsightsSourceMapStorageTag = "hidden-link:Insights.Sourcemap.Storage"

func expandApplicationInsightsTagsWithSourceMap(uri string, input map[string]any) *map[string]string {
	// Start from user tags, then overlay the source-map hidden-link (attribute wins).
	outPtr := tags.Expand(input)
	var out map[string]string
	if outPtr == nil {
		out = make(map[string]string)
	} else {
		out = *outPtr
	}
	if uri == "" {
		delete(out, appInsightsSourceMapStorageTag)
		return &out
	}
	payload, err := json.Marshal(map[string]string{"Uri": uri})
	if err != nil {
		// Unreachable for map[string]string; keep compile-safe.
		return &out
	}
	out[appInsightsSourceMapStorageTag] = string(payload)
	return &out
}
