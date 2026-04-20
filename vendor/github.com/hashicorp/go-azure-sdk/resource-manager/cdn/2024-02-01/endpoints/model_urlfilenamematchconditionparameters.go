package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = URLFileNameMatchConditionParameters{}

type URLFileNameMatchConditionParameters struct {
	MatchValues     *[]string           `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        URLFileNameOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s URLFileNameMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = URLFileNameMatchConditionParameters{}

func (s URLFileNameMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper URLFileNameMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLFileNameMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLFileNameMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleUrlFilenameConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLFileNameMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
