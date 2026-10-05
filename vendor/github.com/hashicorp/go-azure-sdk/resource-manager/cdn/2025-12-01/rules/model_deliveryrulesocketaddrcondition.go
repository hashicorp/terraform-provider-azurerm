package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleCondition = DeliveryRuleSocketAddrCondition{}

type DeliveryRuleSocketAddrCondition struct {
	Parameters SocketAddrMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition

	Name MatchVariable `json:"name"`
}

func (s DeliveryRuleSocketAddrCondition) DeliveryRuleCondition() BaseDeliveryRuleConditionImpl {
	return BaseDeliveryRuleConditionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleSocketAddrCondition{}

func (s DeliveryRuleSocketAddrCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleSocketAddrCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleSocketAddrCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleSocketAddrCondition: %+v", err)
	}

	decoded["name"] = "SocketAddr"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleSocketAddrCondition: %+v", err)
	}

	return encoded, nil
}
