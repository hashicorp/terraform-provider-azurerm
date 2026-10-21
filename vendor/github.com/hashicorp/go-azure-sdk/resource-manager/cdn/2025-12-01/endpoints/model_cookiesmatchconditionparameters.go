package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = CookiesMatchConditionParameters{}

type CookiesMatchConditionParameters struct {
	MatchValues     *[]string       `json:"matchValues,omitempty"`
	NegateCondition *bool           `json:"negateCondition,omitempty"`
	Operator        CookiesOperator `json:"operator"`
	Selector        *string         `json:"selector,omitempty"`
	Transforms      *[]Transform    `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s CookiesMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = CookiesMatchConditionParameters{}

func (s CookiesMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper CookiesMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CookiesMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CookiesMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleCookiesConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CookiesMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
