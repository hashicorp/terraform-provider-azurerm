package endpoints

import (
	"encoding/json"
	"fmt"
)

type DeliveryRule struct {
	Actions    []DeliveryRuleAction     `json:"actions"`
	Conditions *[]DeliveryRuleCondition `json:"conditions,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Order      int64                    `json:"order"`
}

var _ json.Unmarshaler = &DeliveryRule{}

func (s *DeliveryRule) UnmarshalJSON(bytes []byte) error {
	type alias DeliveryRule
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into DeliveryRule: %+v", err)
	}

	s.Actions = decoded.Actions
	s.Name = decoded.Name
	s.Order = decoded.Order

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DeliveryRule into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["conditions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Conditions into list []json.RawMessage: %+v", err)
		}

		output := make([]DeliveryRuleCondition, 0)
		for i, val := range listTemp {
			impl, err := unmarshalDeliveryRuleConditionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Conditions' for 'DeliveryRule': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Conditions = &output
	}
	return nil
}
