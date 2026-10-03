package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScoringFunction interface {
	ScoringFunction() BaseScoringFunctionImpl
}

var _ ScoringFunction = BaseScoringFunctionImpl{}

type BaseScoringFunctionImpl struct {
	Boost         float64                       `json:"boost"`
	FieldName     string                        `json:"fieldName"`
	Interpolation *ScoringFunctionInterpolation `json:"interpolation,omitempty"`
	Type          string                        `json:"type"`
}

func (s BaseScoringFunctionImpl) ScoringFunction() BaseScoringFunctionImpl {
	return s
}

var _ ScoringFunction = RawScoringFunctionImpl{}

// RawScoringFunctionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawScoringFunctionImpl struct {
	scoringFunction BaseScoringFunctionImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawScoringFunctionImpl) ScoringFunction() BaseScoringFunctionImpl {
	return s.scoringFunction
}

func UnmarshalScoringFunctionImplementation(input []byte) (ScoringFunction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ScoringFunction into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "distance") {
		var out DistanceScoringFunction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DistanceScoringFunction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "freshness") {
		var out FreshnessScoringFunction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FreshnessScoringFunction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "magnitude") {
		var out MagnitudeScoringFunction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MagnitudeScoringFunction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "tag") {
		var out TagScoringFunction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TagScoringFunction: %+v", err)
		}
		return out, nil
	}

	var parent BaseScoringFunctionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseScoringFunctionImpl: %+v", err)
	}

	return RawScoringFunctionImpl{
		scoringFunction: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
