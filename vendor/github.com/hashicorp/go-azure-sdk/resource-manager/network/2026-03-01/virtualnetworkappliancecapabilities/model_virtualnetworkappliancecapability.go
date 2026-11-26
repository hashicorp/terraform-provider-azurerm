package virtualnetworkappliancecapabilities

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkApplianceCapability interface {
	VirtualNetworkApplianceCapability() BaseVirtualNetworkApplianceCapabilityImpl
}

var _ VirtualNetworkApplianceCapability = BaseVirtualNetworkApplianceCapabilityImpl{}

type BaseVirtualNetworkApplianceCapabilityImpl struct {
	Id         *string                                      `json:"id,omitempty"`
	Kind       VirtualNetworkApplianceCapabilityKind        `json:"kind"`
	Location   *string                                      `json:"location,omitempty"`
	Name       *string                                      `json:"name,omitempty"`
	Properties *VirtualNetworkApplianceCapabilityProperties `json:"properties,omitempty"`
	Tags       *map[string]string                           `json:"tags,omitempty"`
	Type       *string                                      `json:"type,omitempty"`
}

func (s BaseVirtualNetworkApplianceCapabilityImpl) VirtualNetworkApplianceCapability() BaseVirtualNetworkApplianceCapabilityImpl {
	return s
}

var _ VirtualNetworkApplianceCapability = RawVirtualNetworkApplianceCapabilityImpl{}

// RawVirtualNetworkApplianceCapabilityImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawVirtualNetworkApplianceCapabilityImpl struct {
	virtualNetworkApplianceCapability BaseVirtualNetworkApplianceCapabilityImpl
	Type                              string
	Values                            map[string]interface{}
}

func (s RawVirtualNetworkApplianceCapabilityImpl) VirtualNetworkApplianceCapability() BaseVirtualNetworkApplianceCapabilityImpl {
	return s.virtualNetworkApplianceCapability
}

func (s RawVirtualNetworkApplianceCapabilityImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalVirtualNetworkApplianceCapabilityImplementation(input []byte) (VirtualNetworkApplianceCapability, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling VirtualNetworkApplianceCapability into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "NAT64") {
		var out Nat64Capability
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Nat64Capability: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PLGateway") {
		var out PLGatewayCapability
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PLGatewayCapability: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PLGatewayFastpath") {
		var out PLGatewayFastpathCapability
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PLGatewayFastpathCapability: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PLIPForwarders") {
		var out PLIPForwardersCapability
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PLIPForwardersCapability: %+v", err)
		}
		return out, nil
	}

	var parent BaseVirtualNetworkApplianceCapabilityImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseVirtualNetworkApplianceCapabilityImpl: %+v", err)
	}

	return RawVirtualNetworkApplianceCapabilityImpl{
		virtualNetworkApplianceCapability: parent,
		Type:                              value,
		Values:                            temp,
	}, nil

}
