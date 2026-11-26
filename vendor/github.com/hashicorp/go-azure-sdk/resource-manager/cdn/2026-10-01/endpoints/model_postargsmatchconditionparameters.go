package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = PostArgsMatchConditionParameters{}

type PostArgsMatchConditionParameters struct {
	MatchValues     *[]string        `json:"matchValues,omitempty"`
	NegateCondition *bool            `json:"negateCondition,omitempty"`
	Operator        PostArgsOperator `json:"operator"`
	Selector        *string          `json:"selector,omitempty"`
	Transforms      *[]Transform     `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s PostArgsMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = PostArgsMatchConditionParameters{}

func (s PostArgsMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper PostArgsMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PostArgsMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PostArgsMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRulePostArgsConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PostArgsMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
