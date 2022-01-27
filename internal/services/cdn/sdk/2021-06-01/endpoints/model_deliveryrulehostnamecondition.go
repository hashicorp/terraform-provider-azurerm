package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleHostNameCondition{}

type DeliveryRuleHostNameCondition struct {
	Parameters HostNameMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleHostNameCondition{}

func (s DeliveryRuleHostNameCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleHostNameCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleHostNameCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleHostNameCondition: %+v", err)
	}
	decoded["name"] = "HostName"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleHostNameCondition: %+v", err)
	}

	return encoded, nil
}
