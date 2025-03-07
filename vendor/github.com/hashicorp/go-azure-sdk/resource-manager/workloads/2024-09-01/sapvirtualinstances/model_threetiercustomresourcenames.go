package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreeTierCustomResourceNames interface {
	ThreeTierCustomResourceNames() BaseThreeTierCustomResourceNamesImpl
}

var _ ThreeTierCustomResourceNames = BaseThreeTierCustomResourceNamesImpl{}

type BaseThreeTierCustomResourceNamesImpl struct {
	NamingPatternType NamingPatternType `json:"namingPatternType"`
}

func (s BaseThreeTierCustomResourceNamesImpl) ThreeTierCustomResourceNames() BaseThreeTierCustomResourceNamesImpl {
	return s
}

var _ ThreeTierCustomResourceNames = RawThreeTierCustomResourceNamesImpl{}

// RawThreeTierCustomResourceNamesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawThreeTierCustomResourceNamesImpl struct {
	threeTierCustomResourceNames BaseThreeTierCustomResourceNamesImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawThreeTierCustomResourceNamesImpl) ThreeTierCustomResourceNames() BaseThreeTierCustomResourceNamesImpl {
	return s.threeTierCustomResourceNames
}

func UnmarshalThreeTierCustomResourceNamesImplementation(input []byte) (ThreeTierCustomResourceNames, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ThreeTierCustomResourceNames into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["namingPatternType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "FullResourceName") {
		var out ThreeTierFullResourceNames
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreeTierFullResourceNames: %+v", err)
		}
		return out, nil
	}

	var parent BaseThreeTierCustomResourceNamesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseThreeTierCustomResourceNamesImpl: %+v", err)
	}

	return RawThreeTierCustomResourceNamesImpl{
		threeTierCustomResourceNames: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
