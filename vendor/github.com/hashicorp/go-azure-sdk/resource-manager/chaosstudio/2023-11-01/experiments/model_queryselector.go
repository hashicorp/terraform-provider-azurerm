package experiments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Selector = QuerySelector{}

type QuerySelector struct {
	QueryString     string   `json:"queryString"`
	SubscriptionIds []string `json:"subscriptionIds"`

	// Fields inherited from Selector

	Filter Filter       `json:"filter"`
	Id     string       `json:"id"`
	Type   SelectorType `json:"type"`
}

func (s QuerySelector) Selector() BaseSelectorImpl {
	return BaseSelectorImpl{
		Filter: s.Filter,
		Id:     s.Id,
		Type:   s.Type,
	}
}

var _ json.Marshaler = QuerySelector{}

func (s QuerySelector) MarshalJSON() ([]byte, error) {
	type wrapper QuerySelector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling QuerySelector: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling QuerySelector: %+v", err)
	}

	decoded["type"] = "Query"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling QuerySelector: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &QuerySelector{}

func (s *QuerySelector) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		QueryString     string       `json:"queryString"`
		SubscriptionIds []string     `json:"subscriptionIds"`
		Id              string       `json:"id"`
		Type            SelectorType `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.QueryString = decoded.QueryString
	s.SubscriptionIds = decoded.SubscriptionIds
	s.Id = decoded.Id
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling QuerySelector into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["filter"]; ok {
		impl, err := UnmarshalFilterImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Filter' for 'QuerySelector': %+v", err)
		}
		s.Filter = impl
	}

	return nil
}
