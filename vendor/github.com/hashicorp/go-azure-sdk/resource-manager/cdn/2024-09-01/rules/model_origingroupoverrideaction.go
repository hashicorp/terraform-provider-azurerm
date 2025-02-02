package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = OriginGroupOverrideAction{}

type OriginGroupOverrideAction struct {
	Parameters OriginGroupOverrideActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s OriginGroupOverrideAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = OriginGroupOverrideAction{}

func (s OriginGroupOverrideAction) MarshalJSON() ([]byte, error) {
	type wrapper OriginGroupOverrideAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OriginGroupOverrideAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OriginGroupOverrideAction: %+v", err)
	}

	decoded["name"] = "OriginGroupOverride"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OriginGroupOverrideAction: %+v", err)
	}

	return encoded, nil
}
