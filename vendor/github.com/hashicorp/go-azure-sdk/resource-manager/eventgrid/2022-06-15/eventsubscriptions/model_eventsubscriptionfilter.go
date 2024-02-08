package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSubscriptionFilter struct {
	AdvancedFilters                 *[]AdvancedFilter `json:"advancedFilters,omitempty"`
	EnableAdvancedFilteringOnArrays *bool             `json:"enableAdvancedFilteringOnArrays,omitempty"`
	IncludedEventTypes              *[]string         `json:"includedEventTypes,omitempty"`
	IsSubjectCaseSensitive          *bool             `json:"isSubjectCaseSensitive,omitempty"`
	SubjectBeginsWith               *string           `json:"subjectBeginsWith,omitempty"`
	SubjectEndsWith                 *string           `json:"subjectEndsWith,omitempty"`
}

var _ json.Unmarshaler = &EventSubscriptionFilter{}

func (s *EventSubscriptionFilter) UnmarshalJSON(bytes []byte) error {
	type alias EventSubscriptionFilter
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into EventSubscriptionFilter: %+v", err)
	}

	s.EnableAdvancedFilteringOnArrays = decoded.EnableAdvancedFilteringOnArrays
	s.IncludedEventTypes = decoded.IncludedEventTypes
	s.IsSubjectCaseSensitive = decoded.IsSubjectCaseSensitive
	s.SubjectBeginsWith = decoded.SubjectBeginsWith
	s.SubjectEndsWith = decoded.SubjectEndsWith

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EventSubscriptionFilter into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["advancedFilters"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling AdvancedFilters into list []json.RawMessage: %+v", err)
		}

		output := make([]AdvancedFilter, 0)
		for i, val := range listTemp {
			impl, err := unmarshalAdvancedFilterImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'AdvancedFilters' for 'EventSubscriptionFilter': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.AdvancedFilters = &output
	}
	return nil
}
