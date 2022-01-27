package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleServerPortCondition{}

type DeliveryRuleServerPortCondition struct {
	Parameters ServerPortMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleServerPortCondition{}

func (s DeliveryRuleServerPortCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleServerPortCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleServerPortCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleServerPortCondition: %+v", err)
	}
	decoded["name"] = "ServerPort"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleServerPortCondition: %+v", err)
	}

	return encoded, nil
}
