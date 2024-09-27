package workspaces

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutboundRule = ServiceTagOutboundRule{}

type ServiceTagOutboundRule struct {
	Destination *ServiceTagDestination `json:"destination,omitempty"`

	// Fields inherited from OutboundRule

	Category *RuleCategory `json:"category,omitempty"`
	Status   *RuleStatus   `json:"status,omitempty"`
	Type     RuleType      `json:"type"`
}

func (s ServiceTagOutboundRule) OutboundRule() BaseOutboundRuleImpl {
	return BaseOutboundRuleImpl{
		Category: s.Category,
		Status:   s.Status,
		Type:     s.Type,
	}
}

var _ json.Marshaler = ServiceTagOutboundRule{}

func (s ServiceTagOutboundRule) MarshalJSON() ([]byte, error) {
	type wrapper ServiceTagOutboundRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceTagOutboundRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceTagOutboundRule: %+v", err)
	}

	decoded["type"] = "ServiceTag"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceTagOutboundRule: %+v", err)
	}

	return encoded, nil
}
