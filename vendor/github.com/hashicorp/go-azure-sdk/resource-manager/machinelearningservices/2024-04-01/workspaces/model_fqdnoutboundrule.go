package workspaces

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutboundRule = FqdnOutboundRule{}

type FqdnOutboundRule struct {
	Destination *string `json:"destination,omitempty"`

	// Fields inherited from OutboundRule

	Category *RuleCategory `json:"category,omitempty"`
	Status   *RuleStatus   `json:"status,omitempty"`
	Type     RuleType      `json:"type"`
}

func (s FqdnOutboundRule) OutboundRule() BaseOutboundRuleImpl {
	return BaseOutboundRuleImpl{
		Category: s.Category,
		Status:   s.Status,
		Type:     s.Type,
	}
}

var _ json.Marshaler = FqdnOutboundRule{}

func (s FqdnOutboundRule) MarshalJSON() ([]byte, error) {
	type wrapper FqdnOutboundRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FqdnOutboundRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FqdnOutboundRule: %+v", err)
	}

	decoded["type"] = "FQDN"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FqdnOutboundRule: %+v", err)
	}

	return encoded, nil
}
