package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleRemoteAddressCondition{}

type DeliveryRuleRemoteAddressCondition struct {
	Parameters RemoteAddressMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleRemoteAddressCondition{}

func (s DeliveryRuleRemoteAddressCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleRemoteAddressCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleRemoteAddressCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleRemoteAddressCondition: %+v", err)
	}
	decoded["name"] = "RemoteAddress"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleRemoteAddressCondition: %+v", err)
	}

	return encoded, nil
}
