// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type TemplateDeploymentDebugLevel string

const (
	DebugLevelNone                          TemplateDeploymentDebugLevel = "none"
	DebugLevelRequestContent                TemplateDeploymentDebugLevel = "requestContent"
	DebugLevelResponseContent               TemplateDeploymentDebugLevel = "responseContent"
	DebugLevelRequestContentResponseContent TemplateDeploymentDebugLevel = "requestContent, responseContent"
)

var TemplateDeploymentDebugLevels = []string{
	string(DebugLevelNone),
	string(DebugLevelRequestContent),
	string(DebugLevelResponseContent),
	string(DebugLevelRequestContentResponseContent),
}

func ExpandTemplateDeploymentDebugSetting(debugLevel string) *resources.DebugSetting {
	if debugLevel == "" {
		return &resources.DebugSetting{
			DetailLevel: pointer.To(string(DebugLevelNone)),
		}
	}

	return &resources.DebugSetting{
		DetailLevel: pointer.To(debugLevel),
	}
}

func FlattenTemplateDeploymentDebugSetting(input *resources.DebugSetting) string {
	if input == nil || input.DetailLevel == nil || strings.EqualFold(*input.DetailLevel, string(DebugLevelNone)) {
		return ""
	}

	for _, key := range TemplateDeploymentDebugLevels {
		if strings.EqualFold(key, *input.DetailLevel) {
			return key
		}
	}

	return ""
}

func ExpandTemplateDeploymentBody(input string) (*map[string]any, error) {
	var output map[string]any

	if err := json.Unmarshal([]byte(input), &output); err != nil {
		return nil, err
	}

	return &output, nil
}

func FlattenTemplateDeploymentBody(input any) (*string, error) {
	output := "{}" // since this should be json

	if input == nil {
		return &output, nil
	}

	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshalling json: %+v", err)
	}

	output = string(bytes)
	return &output, nil
}

func FilterOutTemplateDeploymentParameters(input any) any {
	if input == nil {
		return nil
	}

	items, ok := input.(map[string]any)
	if !ok {
		// this is best-effort
		return input
	}

	/*
			Example:

		{
			"someParam": {
				"type": "String",
				"value": "ohhai"
			},
			"dnsLabelPrefix": {
				"reference": {
					"keyvault": {
						"id": "/some/id/that/does-not/matter/right/now"
					},
					"secretName": "some-name"
				}
			}
		}
	*/

	output := make(map[string]any)
	for topLevelKey, topLevelValue := range items {
		if topLevelValue == nil {
			continue
		}

		// give us the original
		output[topLevelKey] = topLevelValue

		// then filter it if necessary
		if innerVals, ok := topLevelValue.(map[string]any); ok {
			outputVals := make(map[string]any)
			for innerKey, innerValue := range innerVals {
				if strings.EqualFold("type", innerKey) {
					continue
				}

				outputVals[innerKey] = innerValue
			}
			output[topLevelKey] = outputVals
		}
	}

	return output
}
