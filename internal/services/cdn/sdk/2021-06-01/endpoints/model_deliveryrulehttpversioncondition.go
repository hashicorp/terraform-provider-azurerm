package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleHttpVersionCondition{}

type DeliveryRuleHttpVersionCondition struct {
	Parameters HttpVersionMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleHttpVersionCondition{}

func (s DeliveryRuleHttpVersionCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleHttpVersionCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleHttpVersionCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleHttpVersionCondition: %+v", err)
	}
	decoded["name"] = "HttpVersion"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleHttpVersionCondition: %+v", err)
	}

	return encoded, nil
}
