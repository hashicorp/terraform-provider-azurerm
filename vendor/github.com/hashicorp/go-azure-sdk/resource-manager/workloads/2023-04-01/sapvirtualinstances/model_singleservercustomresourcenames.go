package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SingleServerCustomResourceNames interface {
}

func unmarshalSingleServerCustomResourceNamesImplementation(input []byte) (SingleServerCustomResourceNames, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SingleServerCustomResourceNames into map[string]interface: %+v", err)
	}

	value, ok := temp["namingPatternType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "FullResourceName") {
		var out SingleServerFullResourceNames
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SingleServerFullResourceNames: %+v", err)
		}
		return out, nil
	}

	type RawSingleServerCustomResourceNamesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSingleServerCustomResourceNamesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
