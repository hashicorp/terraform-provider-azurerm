package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleQueryStringCondition{}

type DeliveryRuleQueryStringCondition struct {
	Parameters QueryStringMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleQueryStringCondition{}

func (s DeliveryRuleQueryStringCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleQueryStringCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleQueryStringCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleQueryStringCondition: %+v", err)
	}
	decoded["name"] = "QueryString"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleQueryStringCondition: %+v", err)
	}

	return encoded, nil
}
