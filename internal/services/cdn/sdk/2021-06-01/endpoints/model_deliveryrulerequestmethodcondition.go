package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleRequestMethodCondition{}

type DeliveryRuleRequestMethodCondition struct {
	Parameters RequestMethodMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleRequestMethodCondition{}

func (s DeliveryRuleRequestMethodCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleRequestMethodCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleRequestMethodCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleRequestMethodCondition: %+v", err)
	}
	decoded["name"] = "RequestMethod"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleRequestMethodCondition: %+v", err)
	}

	return encoded, nil
}
