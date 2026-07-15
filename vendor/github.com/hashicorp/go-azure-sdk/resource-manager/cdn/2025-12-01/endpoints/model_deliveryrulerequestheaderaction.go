package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = DeliveryRuleRequestHeaderAction{}

type DeliveryRuleRequestHeaderAction struct {
	Parameters HeaderActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s DeliveryRuleRequestHeaderAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleRequestHeaderAction{}

func (s DeliveryRuleRequestHeaderAction) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleRequestHeaderAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleRequestHeaderAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleRequestHeaderAction: %+v", err)
	}

	decoded["name"] = "ModifyRequestHeader"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleRequestHeaderAction: %+v", err)
	}

	return encoded, nil
}
