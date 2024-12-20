package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = URLPathMatchConditionParameters{}

type URLPathMatchConditionParameters struct {
	MatchValues     *[]string       `json:"matchValues,omitempty"`
	NegateCondition *bool           `json:"negateCondition,omitempty"`
	Operator        URLPathOperator `json:"operator"`
	Transforms      *[]Transform    `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s URLPathMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = URLPathMatchConditionParameters{}

func (s URLPathMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper URLPathMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLPathMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLPathMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleUrlPathMatchConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLPathMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
