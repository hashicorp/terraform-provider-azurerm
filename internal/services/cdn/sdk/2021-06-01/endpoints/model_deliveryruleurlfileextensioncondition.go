package endpoints

import (
	"encoding/json"
	"fmt"
)

var _ DeliveryRuleCondition = DeliveryRuleUrlFileExtensionCondition{}

type DeliveryRuleUrlFileExtensionCondition struct {
	Parameters UrlFileExtensionMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition
}

var _ json.Marshaler = DeliveryRuleUrlFileExtensionCondition{}

func (s DeliveryRuleUrlFileExtensionCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleUrlFileExtensionCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleUrlFileExtensionCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleUrlFileExtensionCondition: %+v", err)
	}
	decoded["name"] = "UrlFileExtension"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleUrlFileExtensionCondition: %+v", err)
	}

	return encoded, nil
}
