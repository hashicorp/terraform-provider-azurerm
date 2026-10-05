package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = OriginGroupOverrideActionParameters{}

type OriginGroupOverrideActionParameters struct {
	OriginGroup ResourceReference `json:"originGroup"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s OriginGroupOverrideActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = OriginGroupOverrideActionParameters{}

func (s OriginGroupOverrideActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper OriginGroupOverrideActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OriginGroupOverrideActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OriginGroupOverrideActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleOriginGroupOverrideActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OriginGroupOverrideActionParameters: %+v", err)
	}

	return encoded, nil
}
