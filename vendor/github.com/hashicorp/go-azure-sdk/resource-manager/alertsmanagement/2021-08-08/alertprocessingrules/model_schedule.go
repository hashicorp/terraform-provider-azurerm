package alertprocessingrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Schedule struct {
	EffectiveFrom  *string       `json:"effectiveFrom,omitempty"`
	EffectiveUntil *string       `json:"effectiveUntil,omitempty"`
	Recurrences    *[]Recurrence `json:"recurrences,omitempty"`
	TimeZone       *string       `json:"timeZone,omitempty"`
}

var _ json.Unmarshaler = &Schedule{}

func (s *Schedule) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EffectiveFrom  *string `json:"effectiveFrom,omitempty"`
		EffectiveUntil *string `json:"effectiveUntil,omitempty"`
		TimeZone       *string `json:"timeZone,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EffectiveFrom = decoded.EffectiveFrom
	s.EffectiveUntil = decoded.EffectiveUntil
	s.TimeZone = decoded.TimeZone

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Schedule into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["recurrences"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Recurrences into list []json.RawMessage: %+v", err)
		}

		output := make([]Recurrence, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalRecurrenceImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Recurrences' for 'Schedule': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Recurrences = &output
	}

	return nil
}
