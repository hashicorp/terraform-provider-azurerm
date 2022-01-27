package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleClientPortCondition{}

type DeliveryRuleClientPortCondition struct {
	Parameters ClientPortMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleClientPortCondition{}

func (s DeliveryRuleClientPortCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleClientPortCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleClientPortCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleClientPortCondition: %+v", err)
	}
	decoded["name"] = "ClientPort"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleClientPortCondition: %+v", err)
	}

	return encoded, nil
}
