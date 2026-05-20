package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleConditionParameters = IsDeviceMatchConditionParameters{}

type IsDeviceMatchConditionParameters struct {
	MatchValues     *[]IsDeviceMatchValue `json:"matchValues,omitempty"`
	NegateCondition *bool                 `json:"negateCondition,omitempty"`
	Operator        IsDeviceOperator      `json:"operator"`
	Transforms      *[]Transform          `json:"transforms,omitempty"`

	// Fields inherited from DeliveryRuleConditionParameters

	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s IsDeviceMatchConditionParameters) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return BaseDeliveryRuleConditionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = IsDeviceMatchConditionParameters{}

func (s IsDeviceMatchConditionParameters) MarshalJSON() ([]byte, error) {
	type wrapper IsDeviceMatchConditionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IsDeviceMatchConditionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IsDeviceMatchConditionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleIsDeviceConditionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IsDeviceMatchConditionParameters: %+v", err)
	}

	return encoded, nil
}
