package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleRequestBodyCondition{}

type DeliveryRuleRequestBodyCondition struct {
	Parameters RequestBodyMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleRequestBodyCondition{}

func (s DeliveryRuleRequestBodyCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleRequestBodyCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleRequestBodyCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleRequestBodyCondition: %+v", err)
	}
	decoded["name"] = "RequestBody"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleRequestBodyCondition: %+v", err)
	}

	return encoded, nil
}
