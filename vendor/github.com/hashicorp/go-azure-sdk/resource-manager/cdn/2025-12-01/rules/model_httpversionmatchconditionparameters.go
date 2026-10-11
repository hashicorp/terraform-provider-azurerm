package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = HTTPVersionMatchConditionParameters{}

type HTTPVersionMatchConditionParameters struct {
	MatchValues     *[]string           `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        HTTPVersionOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s HTTPVersionMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = HTTPVersionMatchConditionParameters{}

func (s HTTPVersionMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper HTTPVersionMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HTTPVersionMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HTTPVersionMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleHttpVersionConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HTTPVersionMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
