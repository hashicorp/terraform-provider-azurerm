package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = URLRewriteActionParameters{}

type URLRewriteActionParameters struct {
	Destination           string `json:"destination"`
	PreserveUnmatchedPath *bool  `json:"preserveUnmatchedPath,omitempty"`
	SourcePattern         string `json:"sourcePattern"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s URLRewriteActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = URLRewriteActionParameters{}

func (s URLRewriteActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper URLRewriteActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLRewriteActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLRewriteActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleUrlRewriteActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLRewriteActionParameters: %+v", err)
	}

	return encoded, nil
}
