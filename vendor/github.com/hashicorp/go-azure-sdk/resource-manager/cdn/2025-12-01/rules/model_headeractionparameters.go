package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = HeaderActionParameters{}

type HeaderActionParameters struct {
	HeaderAction HeaderAction `json:"headerAction"`
	HeaderName   string       `json:"headerName"`
	Value        *string      `json:"value,omitempty"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s HeaderActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = HeaderActionParameters{}

func (s HeaderActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper HeaderActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HeaderActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HeaderActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleHeaderActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HeaderActionParameters: %+v", err)
	}

	return encoded, nil
}
