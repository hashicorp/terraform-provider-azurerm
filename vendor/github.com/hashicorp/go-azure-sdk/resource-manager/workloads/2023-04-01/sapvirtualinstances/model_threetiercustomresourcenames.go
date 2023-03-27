package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreeTierCustomResourceNames interface {
}

func unmarshalThreeTierCustomResourceNamesImplementation(input []byte) (ThreeTierCustomResourceNames, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ThreeTierCustomResourceNames into map[string]interface: %+v", err)
	}

	value, ok := temp["namingPatternType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "FullResourceName") {
		var out ThreeTierFullResourceNames
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreeTierFullResourceNames: %+v", err)
		}
		return out, nil
	}

	type RawThreeTierCustomResourceNamesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawThreeTierCustomResourceNamesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
