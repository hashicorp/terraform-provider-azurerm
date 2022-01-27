package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleUrlFileNameCondition{}

type DeliveryRuleUrlFileNameCondition struct {
	Parameters UrlFileNameMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleUrlFileNameCondition{}

func (s DeliveryRuleUrlFileNameCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleUrlFileNameCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleUrlFileNameCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleUrlFileNameCondition: %+v", err)
	}
	decoded["name"] = "UrlFileName"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleUrlFileNameCondition: %+v", err)
	}

	return encoded, nil
}
