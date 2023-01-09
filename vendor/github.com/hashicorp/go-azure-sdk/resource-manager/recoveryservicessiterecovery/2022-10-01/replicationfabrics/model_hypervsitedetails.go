package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificDetails = HyperVSiteDetails{}

type HyperVSiteDetails struct {
	HyperVHosts *[]HyperVHostDetails `json:"hyperVHosts,omitempty"`

	// Fields inherited from FabricSpecificDetails
}

var _ json.Marshaler = HyperVSiteDetails{}

func (s HyperVSiteDetails) MarshalJSON() ([]byte, error) {
	type wrapper HyperVSiteDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVSiteDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVSiteDetails: %+v", err)
	}
	decoded["instanceType"] = "HyperVSite"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVSiteDetails: %+v", err)
	}

	return encoded, nil
}
