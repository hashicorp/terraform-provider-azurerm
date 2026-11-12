package rulesets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = DeliveryRuleEdgeActionParameters{}

type DeliveryRuleEdgeActionParameters struct {
	EdgeActionReference ResourceReference `json:"edgeActionReference"`
	InvocationPoint     InvocationPoint   `json:"invocationPoint"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s DeliveryRuleEdgeActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = DeliveryRuleEdgeActionParameters{}

func (s DeliveryRuleEdgeActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleEdgeActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleEdgeActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleEdgeActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleEdgeActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleEdgeActionParameters: %+v", err)
	}

	return encoded, nil
}
