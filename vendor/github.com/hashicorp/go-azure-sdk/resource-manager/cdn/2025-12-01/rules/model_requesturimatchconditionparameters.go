package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = RequestUriMatchConditionParameters{}

type RequestUriMatchConditionParameters struct {
	MatchValues     *[]string          `json:"matchValues,omitempty"`
	NegateCondition *bool              `json:"negateCondition,omitempty"`
	Operator        RequestUriOperator `json:"operator"`
	Transforms      *[]Transform       `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s RequestUriMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = RequestUriMatchConditionParameters{}

func (s RequestUriMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper RequestUriMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RequestUriMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RequestUriMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleRequestUriConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RequestUriMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
