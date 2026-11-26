// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_definition

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
)

func expandParameterDefinitionsValueFromStringTrack1(jsonString string) (map[string]*policy.ParameterDefinitionsValue, error) {
	var result map[string]*policy.ParameterDefinitionsValue

	err := json.Unmarshal([]byte(jsonString), &result)

	return result, err
}
