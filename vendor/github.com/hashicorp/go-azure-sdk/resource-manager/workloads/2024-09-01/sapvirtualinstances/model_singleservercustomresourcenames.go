package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SingleServerCustomResourceNames interface {
	SingleServerCustomResourceNames() BaseSingleServerCustomResourceNamesImpl
}

var _ SingleServerCustomResourceNames = BaseSingleServerCustomResourceNamesImpl{}

type BaseSingleServerCustomResourceNamesImpl struct {
	NamingPatternType NamingPatternType `json:"namingPatternType"`
}

func (s BaseSingleServerCustomResourceNamesImpl) SingleServerCustomResourceNames() BaseSingleServerCustomResourceNamesImpl {
	return s
}

var _ SingleServerCustomResourceNames = RawSingleServerCustomResourceNamesImpl{}

// RawSingleServerCustomResourceNamesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSingleServerCustomResourceNamesImpl struct {
	singleServerCustomResourceNames BaseSingleServerCustomResourceNamesImpl
	Type                            string
	Values                          map[string]interface{}
}

func (s RawSingleServerCustomResourceNamesImpl) SingleServerCustomResourceNames() BaseSingleServerCustomResourceNamesImpl {
	return s.singleServerCustomResourceNames
}

func (s RawSingleServerCustomResourceNamesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSingleServerCustomResourceNamesImplementation(input []byte) (SingleServerCustomResourceNames, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SingleServerCustomResourceNames into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["namingPatternType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "FullResourceName") {
		var out SingleServerFullResourceNames
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SingleServerFullResourceNames: %+v", err)
		}
		return out, nil
	}

	var parent BaseSingleServerCustomResourceNamesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSingleServerCustomResourceNamesImpl: %+v", err)
	}

	return RawSingleServerCustomResourceNamesImpl{
		singleServerCustomResourceNames: parent,
		Type:                            value,
		Values:                          temp,
	}, nil

}
