package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = RequestSchemeMatchConditionParameters{}

type RequestSchemeMatchConditionParameters struct {
	MatchValues     *[]RequestSchemeMatchValue                    `json:"matchValues,omitempty"`
	NegateCondition *bool                                         `json:"negateCondition,omitempty"`
	Operator        RequestSchemeMatchConditionParametersOperator `json:"operator"`
	Transforms      *[]Transform                                  `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s RequestSchemeMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = RequestSchemeMatchConditionParameters{}

func (s RequestSchemeMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper RequestSchemeMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RequestSchemeMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RequestSchemeMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleRequestSchemeConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RequestSchemeMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
