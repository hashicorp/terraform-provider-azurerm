package workspaces

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedNetworkSettings struct {
	FirewallPublicIPAddress *string                        `json:"firewallPublicIpAddress,omitempty"`
	FirewallSku             *FirewallSku                   `json:"firewallSku,omitempty"`
	IsolationMode           *IsolationMode                 `json:"isolationMode,omitempty"`
	ManagedNetworkKind      *ManagedNetworkKind            `json:"managedNetworkKind,omitempty"`
	NetworkId               *string                        `json:"networkId,omitempty"`
	OutboundRules           *map[string]OutboundRule       `json:"outboundRules,omitempty"`
	Status                  *ManagedNetworkProvisionStatus `json:"status,omitempty"`
}

var _ json.Unmarshaler = &ManagedNetworkSettings{}

func (s *ManagedNetworkSettings) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FirewallPublicIPAddress *string                        `json:"firewallPublicIpAddress,omitempty"`
		FirewallSku             *FirewallSku                   `json:"firewallSku,omitempty"`
		IsolationMode           *IsolationMode                 `json:"isolationMode,omitempty"`
		ManagedNetworkKind      *ManagedNetworkKind            `json:"managedNetworkKind,omitempty"`
		NetworkId               *string                        `json:"networkId,omitempty"`
		Status                  *ManagedNetworkProvisionStatus `json:"status,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FirewallPublicIPAddress = decoded.FirewallPublicIPAddress
	s.FirewallSku = decoded.FirewallSku
	s.IsolationMode = decoded.IsolationMode
	s.ManagedNetworkKind = decoded.ManagedNetworkKind
	s.NetworkId = decoded.NetworkId
	s.Status = decoded.Status

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ManagedNetworkSettings into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["outboundRules"]; ok {
		var dictionaryTemp map[string]json.RawMessage
		if err := json.Unmarshal(v, &dictionaryTemp); err != nil {
			return fmt.Errorf("unmarshaling OutboundRules into dictionary map[string]json.RawMessage: %+v", err)
		}

		output := make(map[string]OutboundRule)
		for key, val := range dictionaryTemp {
			impl, err := UnmarshalOutboundRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling key %q field 'OutboundRules' for 'ManagedNetworkSettings': %+v", key, err)
			}
			output[key] = impl
		}
		s.OutboundRules = &output
	}

	return nil
}
