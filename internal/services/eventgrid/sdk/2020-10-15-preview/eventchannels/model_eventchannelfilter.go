package eventchannels

import (
	"encoding/json"
	"fmt"
)

type EventChannelFilter struct {
	AdvancedFilters                 *[]AdvancedFilter `json:"advancedFilters,omitempty"`
	EnableAdvancedFilteringOnArrays *bool             `json:"enableAdvancedFilteringOnArrays,omitempty"`
}

var _ json.Unmarshaler = &EventChannelFilter{}

func (s *EventChannelFilter) UnmarshalJSON(bytes []byte) error {
	type alias EventChannelFilter
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into EventChannelFilter: %+v", err)
	}

	s.EnableAdvancedFilteringOnArrays = decoded.EnableAdvancedFilteringOnArrays

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EventChannelFilter into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'AdvancedFilters' for 'EventChannelFilter': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.AdvancedFilters = &output
	}
	return nil
}
