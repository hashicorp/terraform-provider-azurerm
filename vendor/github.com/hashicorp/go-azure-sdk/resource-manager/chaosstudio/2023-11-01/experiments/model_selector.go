package experiments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Selector interface {
	Selector() BaseSelectorImpl
}

var _ Selector = BaseSelectorImpl{}

type BaseSelectorImpl struct {
	Filter Filter       `json:"filter"`
	Id     string       `json:"id"`
	Type   SelectorType `json:"type"`
}

func (s BaseSelectorImpl) Selector() BaseSelectorImpl {
	return s
}

var _ Selector = RawSelectorImpl{}

// RawSelectorImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSelectorImpl struct {
	selector BaseSelectorImpl
	Type     string
	Values   map[string]interface{}
}

func (s RawSelectorImpl) Selector() BaseSelectorImpl {
	return s.selector
}

var _ json.Unmarshaler = &BaseSelectorImpl{}

func (s *BaseSelectorImpl) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Id   string       `json:"id"`
		Type SelectorType `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Id = decoded.Id
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BaseSelectorImpl into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["filter"]; ok {
		impl, err := UnmarshalFilterImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Filter' for 'BaseSelectorImpl': %+v", err)
		}
		s.Filter = impl
	}

	return nil
}

func UnmarshalSelectorImplementation(input []byte) (Selector, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Selector into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "List") {
		var out ListSelector
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ListSelector: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Query") {
		var out QuerySelector
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into QuerySelector: %+v", err)
		}
		return out, nil
	}

	var parent BaseSelectorImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSelectorImpl: %+v", err)
	}

	return RawSelectorImpl{
		selector: parent,
		Type:     value,
		Values:   temp,
	}, nil

}
