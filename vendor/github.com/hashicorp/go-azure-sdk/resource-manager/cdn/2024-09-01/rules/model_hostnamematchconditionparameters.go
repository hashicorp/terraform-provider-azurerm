package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = HostNameMatchConditionParameters{}

type HostNameMatchConditionParameters struct {
	MatchValues     *[]string        `json:"matchValues,omitempty"`
	NegateCondition *bool            `json:"negateCondition,omitempty"`
	Operator        HostNameOperator `json:"operator"`
	Transforms      *[]Transform     `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s HostNameMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = HostNameMatchConditionParameters{}

func (s HostNameMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper HostNameMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HostNameMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HostNameMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleHostNameConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HostNameMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
