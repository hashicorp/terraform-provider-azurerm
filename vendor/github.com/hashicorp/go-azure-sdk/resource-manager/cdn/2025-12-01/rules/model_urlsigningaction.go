package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = URLSigningAction{}

type URLSigningAction struct {
	Parameters URLSigningActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s URLSigningAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = URLSigningAction{}

func (s URLSigningAction) MarshalJSON() ([]byte, error) {
	type wrapper URLSigningAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLSigningAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLSigningAction: %+v", err)
	}

	decoded["name"] = "UrlSigning"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLSigningAction: %+v", err)
	}

	return encoded, nil
}
