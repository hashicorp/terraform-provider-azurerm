package virtualnetworkappliancecapabilities

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VirtualNetworkApplianceCapability = PLGatewayCapability{}

type PLGatewayCapability struct {

	// Fields inherited from VirtualNetworkApplianceCapability

	Id         *string                                      `json:"id,omitempty"`
	Kind       VirtualNetworkApplianceCapabilityKind        `json:"kind"`
	Location   *string                                      `json:"location,omitempty"`
	Name       *string                                      `json:"name,omitempty"`
	Properties *VirtualNetworkApplianceCapabilityProperties `json:"properties,omitempty"`
	Tags       *map[string]string                           `json:"tags,omitempty"`
	Type       *string                                      `json:"type,omitempty"`
}

func (s PLGatewayCapability) VirtualNetworkApplianceCapability() BaseVirtualNetworkApplianceCapabilityImpl {
	return BaseVirtualNetworkApplianceCapabilityImpl{
		Id:         s.Id,
		Kind:       s.Kind,
		Location:   s.Location,
		Name:       s.Name,
		Properties: s.Properties,
		Tags:       s.Tags,
		Type:       s.Type,
	}
}

var _ json.Marshaler = PLGatewayCapability{}

func (s PLGatewayCapability) MarshalJSON() ([]byte, error) {
	type wrapper PLGatewayCapability
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PLGatewayCapability: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PLGatewayCapability: %+v", err)
	}

	decoded["kind"] = "PLGateway"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PLGatewayCapability: %+v", err)
	}

	return encoded, nil
}
