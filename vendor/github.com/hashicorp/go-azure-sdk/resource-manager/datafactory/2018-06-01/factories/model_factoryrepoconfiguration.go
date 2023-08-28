package factories

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FactoryRepoConfiguration interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFactoryRepoConfigurationImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalFactoryRepoConfigurationImplementation(input []byte) (FactoryRepoConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FactoryRepoConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "FactoryGitHubConfiguration") {
		var out FactoryGitHubConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FactoryGitHubConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FactoryVSTSConfiguration") {
		var out FactoryVSTSConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FactoryVSTSConfiguration: %+v", err)
		}
		return out, nil
	}

	out := RawFactoryRepoConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
