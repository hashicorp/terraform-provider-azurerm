package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleUrlPathCondition{}

type DeliveryRuleUrlPathCondition struct {
	Parameters UrlPathMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleUrlPathCondition{}

func (s DeliveryRuleUrlPathCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleUrlPathCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleUrlPathCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleUrlPathCondition: %+v", err)
	}
	decoded["name"] = "UrlPath"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleUrlPathCondition: %+v", err)
	}

	return encoded, nil
}
