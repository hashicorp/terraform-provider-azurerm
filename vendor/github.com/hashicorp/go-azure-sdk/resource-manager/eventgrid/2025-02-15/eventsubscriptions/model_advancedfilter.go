package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvancedFilter interface {
	AdvancedFilter() BaseAdvancedFilterImpl
}

var _ AdvancedFilter = BaseAdvancedFilterImpl{}

type BaseAdvancedFilterImpl struct {
	Key          *string                    `json:"key,omitempty"`
	OperatorType AdvancedFilterOperatorType `json:"operatorType"`
}

func (s BaseAdvancedFilterImpl) AdvancedFilter() BaseAdvancedFilterImpl {
	return s
}

var _ AdvancedFilter = RawAdvancedFilterImpl{}

// RawAdvancedFilterImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAdvancedFilterImpl struct {
	advancedFilter BaseAdvancedFilterImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawAdvancedFilterImpl) AdvancedFilter() BaseAdvancedFilterImpl {
	return s.advancedFilter
}

func UnmarshalAdvancedFilterImplementation(input []byte) (AdvancedFilter, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AdvancedFilter into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["operatorType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BoolEquals") {
		var out BoolEqualsAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BoolEqualsAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IsNotNull") {
		var out IsNotNullAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IsNotNullAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IsNullOrUndefined") {
		var out IsNullOrUndefinedAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IsNullOrUndefinedAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberGreaterThan") {
		var out NumberGreaterThanAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberGreaterThanAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberGreaterThanOrEquals") {
		var out NumberGreaterThanOrEqualsAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberGreaterThanOrEqualsAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberIn") {
		var out NumberInAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberInAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberInRange") {
		var out NumberInRangeAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberInRangeAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberLessThan") {
		var out NumberLessThanAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberLessThanAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberLessThanOrEquals") {
		var out NumberLessThanOrEqualsAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberLessThanOrEqualsAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberNotIn") {
		var out NumberNotInAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberNotInAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberNotInRange") {
		var out NumberNotInRangeAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberNotInRangeAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringBeginsWith") {
		var out StringBeginsWithAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringBeginsWithAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringContains") {
		var out StringContainsAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringContainsAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringEndsWith") {
		var out StringEndsWithAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringEndsWithAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringIn") {
		var out StringInAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringInAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotBeginsWith") {
		var out StringNotBeginsWithAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotBeginsWithAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotContains") {
		var out StringNotContainsAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotContainsAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotEndsWith") {
		var out StringNotEndsWithAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotEndsWithAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotIn") {
		var out StringNotInAdvancedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotInAdvancedFilter: %+v", err)
		}
		return out, nil
	}

	var parent BaseAdvancedFilterImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAdvancedFilterImpl: %+v", err)
	}

	return RawAdvancedFilterImpl{
		advancedFilter: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
