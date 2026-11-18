package virtualnetworkappliancecapabilities

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VirtualNetworkApplianceCapability = PLIPForwardersCapability{}

type PLIPForwardersCapability struct {

	// Fields inherited from VirtualNetworkApplianceCapability

	Id         *string                                      `json:"id,omitempty"`
	Kind       VirtualNetworkApplianceCapabilityKind        `json:"kind"`
	Location   *string                                      `json:"location,omitempty"`
	Name       *string                                      `json:"name,omitempty"`
	Properties *VirtualNetworkApplianceCapabilityProperties `json:"properties,omitempty"`
	Tags       *map[string]string                           `json:"tags,omitempty"`
	Type       *string                                      `json:"type,omitempty"`
}

func (s PLIPForwardersCapability) VirtualNetworkApplianceCapability() BaseVirtualNetworkApplianceCapabilityImpl {
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

var _ json.Marshaler = PLIPForwardersCapability{}

func (s PLIPForwardersCapability) MarshalJSON() ([]byte, error) {
	type wrapper PLIPForwardersCapability
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PLIPForwardersCapability: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PLIPForwardersCapability: %+v", err)
	}

	decoded["kind"] = "PLIPForwarders"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PLIPForwardersCapability: %+v", err)
	}

	return encoded, nil
}
