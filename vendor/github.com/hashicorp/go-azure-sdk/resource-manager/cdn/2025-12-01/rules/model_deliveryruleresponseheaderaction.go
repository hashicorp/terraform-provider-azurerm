package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = DeliveryRuleResponseHeaderAction{}

type DeliveryRuleResponseHeaderAction struct {
	Parameters HeaderActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s DeliveryRuleResponseHeaderAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleResponseHeaderAction{}

func (s DeliveryRuleResponseHeaderAction) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleResponseHeaderAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleResponseHeaderAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleResponseHeaderAction: %+v", err)
	}

	decoded["name"] = "ModifyResponseHeader"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleResponseHeaderAction: %+v", err)
	}

	return encoded, nil
}
