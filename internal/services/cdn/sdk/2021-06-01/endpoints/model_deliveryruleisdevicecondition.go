package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleIsDeviceCondition{}

type DeliveryRuleIsDeviceCondition struct {
	Parameters IsDeviceMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleIsDeviceCondition{}

func (s DeliveryRuleIsDeviceCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleIsDeviceCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleIsDeviceCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleIsDeviceCondition: %+v", err)
	}
	decoded["name"] = "IsDevice"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleIsDeviceCondition: %+v", err)
	}

	return encoded, nil
}
