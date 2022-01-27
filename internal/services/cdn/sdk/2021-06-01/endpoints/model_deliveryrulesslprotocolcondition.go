package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleSslProtocolCondition{}

type DeliveryRuleSslProtocolCondition struct {
	Parameters SslProtocolMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleSslProtocolCondition{}

func (s DeliveryRuleSslProtocolCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleSslProtocolCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleSslProtocolCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleSslProtocolCondition: %+v", err)
	}
	decoded["name"] = "SslProtocol"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleSslProtocolCondition: %+v", err)
	}

	return encoded, nil
}
