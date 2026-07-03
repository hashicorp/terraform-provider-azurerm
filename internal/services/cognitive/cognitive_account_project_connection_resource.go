// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func flattenProjectConnectionMetadata(priorMetadata map[string]string, apiMetadata *map[string]string) map[string]string {
	// The API returns additional metadata fields beyond those configured (e.g. `ApiVersion`,
	// `DeploymentApiVersion`). When prior configuration is known (Read/Update) only the configured
	// keys are surfaced to avoid diffs; otherwise (import or list) all API metadata fields are returned.
	apiMetadataValues := pointer.From(apiMetadata)
	if len(priorMetadata) == 0 {
		return apiMetadataValues
	}

	filtered := make(map[string]string)
	for configKey := range priorMetadata {
		for apiKey, apiValue := range apiMetadataValues {
			if strings.EqualFold(configKey, apiKey) {
				filtered[configKey] = apiValue
				break
			}
		}
	}
	return filtered
}
