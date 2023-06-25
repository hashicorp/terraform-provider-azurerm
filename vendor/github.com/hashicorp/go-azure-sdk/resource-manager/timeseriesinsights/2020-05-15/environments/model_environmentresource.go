package environments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentResource interface {
}

func unmarshalEnvironmentResourceImplementation(input []byte) (EnvironmentResource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EnvironmentResource into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Gen1") {
		var out Gen1EnvironmentResource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Gen1EnvironmentResource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Gen2") {
		var out Gen2EnvironmentResource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Gen2EnvironmentResource: %+v", err)
		}
		return out, nil
	}

	type RawEnvironmentResourceImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEnvironmentResourceImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
