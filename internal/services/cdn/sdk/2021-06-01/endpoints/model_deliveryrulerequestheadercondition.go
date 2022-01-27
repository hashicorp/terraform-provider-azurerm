package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleRequestHeaderCondition{}

type DeliveryRuleRequestHeaderCondition struct {
	Parameters RequestHeaderMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleRequestHeaderCondition{}

func (s DeliveryRuleRequestHeaderCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleRequestHeaderCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleRequestHeaderCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleRequestHeaderCondition: %+v", err)
	}
	decoded["name"] = "RequestHeader"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleRequestHeaderCondition: %+v", err)
	}

	return encoded, nil
}
