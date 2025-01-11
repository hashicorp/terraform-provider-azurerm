package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = RequestMethodMatchConditionParameters{}

type RequestMethodMatchConditionParameters struct {
	MatchValues     *[]RequestMethodMatchValue `json:"matchValues,omitempty"`
	NegateCondition *bool                      `json:"negateCondition,omitempty"`
	Operator        RequestMethodOperator      `json:"operator"`
	Transforms      *[]Transform               `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s RequestMethodMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = RequestMethodMatchConditionParameters{}

func (s RequestMethodMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper RequestMethodMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RequestMethodMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RequestMethodMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleRequestMethodConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RequestMethodMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
