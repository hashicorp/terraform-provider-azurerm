package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = SslProtocolMatchConditionParameters{}

type SslProtocolMatchConditionParameters struct {
	MatchValues     *[]SslProtocol      `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        SslProtocolOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s SslProtocolMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = SslProtocolMatchConditionParameters{}

func (s SslProtocolMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper SslProtocolMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SslProtocolMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SslProtocolMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleSslProtocolConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SslProtocolMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
