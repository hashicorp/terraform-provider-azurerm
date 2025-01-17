package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = URLSigningActionParameters{}

type URLSigningActionParameters struct {
	Algorithm             *Algorithm                   `json:"algorithm,omitempty"`
	ParameterNameOverride *[]URLSigningParamIdentifier `json:"parameterNameOverride,omitempty"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s URLSigningActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = URLSigningActionParameters{}

func (s URLSigningActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper URLSigningActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLSigningActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLSigningActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleUrlSigningActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLSigningActionParameters: %+v", err)
	}

	return encoded, nil
}
