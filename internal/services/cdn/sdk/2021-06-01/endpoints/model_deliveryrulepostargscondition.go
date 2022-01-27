package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRulePostArgsCondition{}

type DeliveryRulePostArgsCondition struct {
	Parameters PostArgsMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRulePostArgsCondition{}

func (s DeliveryRulePostArgsCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRulePostArgsCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRulePostArgsCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRulePostArgsCondition: %+v", err)
	}
	decoded["name"] = "PostArgs"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRulePostArgsCondition: %+v", err)
	}

	return encoded, nil
}
