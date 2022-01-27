package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleRequestSchemeCondition{}

type DeliveryRuleRequestSchemeCondition struct {
	Parameters RequestSchemeMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleRequestSchemeCondition{}

func (s DeliveryRuleRequestSchemeCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleRequestSchemeCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleRequestSchemeCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleRequestSchemeCondition: %+v", err)
	}
	decoded["name"] = "RequestScheme"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleRequestSchemeCondition: %+v", err)
	}

	return encoded, nil
}
