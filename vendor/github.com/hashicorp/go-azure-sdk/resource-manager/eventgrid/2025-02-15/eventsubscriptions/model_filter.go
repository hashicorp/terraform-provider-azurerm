package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Filter interface {
	Filter() BaseFilterImpl
}

var _ Filter = BaseFilterImpl{}

type BaseFilterImpl struct {
	Key          *string            `json:"key,omitempty"`
	OperatorType FilterOperatorType `json:"operatorType"`
}

func (s BaseFilterImpl) Filter() BaseFilterImpl {
	return s
}

var _ Filter = RawFilterImpl{}

// RawFilterImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFilterImpl struct {
	filter BaseFilterImpl
	Type   string
	Values map[string]interface{}
}

func (s RawFilterImpl) Filter() BaseFilterImpl {
	return s.filter
}

func UnmarshalFilterImplementation(input []byte) (Filter, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Filter into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["operatorType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BoolEquals") {
		var out BoolEqualsFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BoolEqualsFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IsNotNull") {
		var out IsNotNullFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IsNotNullFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IsNullOrUndefined") {
		var out IsNullOrUndefinedFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IsNullOrUndefinedFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberGreaterThan") {
		var out NumberGreaterThanFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberGreaterThanFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberGreaterThanOrEquals") {
		var out NumberGreaterThanOrEqualsFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberGreaterThanOrEqualsFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberIn") {
		var out NumberInFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberInFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberInRange") {
		var out NumberInRangeFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberInRangeFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberLessThan") {
		var out NumberLessThanFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberLessThanFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberLessThanOrEquals") {
		var out NumberLessThanOrEqualsFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberLessThanOrEqualsFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberNotIn") {
		var out NumberNotInFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberNotInFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NumberNotInRange") {
		var out NumberNotInRangeFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NumberNotInRangeFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringBeginsWith") {
		var out StringBeginsWithFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringBeginsWithFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringContains") {
		var out StringContainsFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringContainsFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringEndsWith") {
		var out StringEndsWithFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringEndsWithFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringIn") {
		var out StringInFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringInFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotBeginsWith") {
		var out StringNotBeginsWithFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotBeginsWithFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotContains") {
		var out StringNotContainsFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotContainsFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotEndsWith") {
		var out StringNotEndsWithFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotEndsWithFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StringNotIn") {
		var out StringNotInFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StringNotInFilter: %+v", err)
		}
		return out, nil
	}

	var parent BaseFilterImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFilterImpl: %+v", err)
	}

	return RawFilterImpl{
		filter: parent,
		Type:   value,
		Values: temp,
	}, nil

}
