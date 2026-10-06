package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificDetails = VMwareV2FabricSpecificDetails{}

type VMwareV2FabricSpecificDetails struct {
	MigrationSolutionId *string                 `json:"migrationSolutionId,omitempty"`
	PhysicalSiteId      *string                 `json:"physicalSiteId,omitempty"`
	ProcessServers      *[]ProcessServerDetails `json:"processServers,omitempty"`
	ServiceContainerId  *string                 `json:"serviceContainerId,omitempty"`
	ServiceEndpoint     *string                 `json:"serviceEndpoint,omitempty"`
	ServiceResourceId   *string                 `json:"serviceResourceId,omitempty"`
	VMwareSiteId        *string                 `json:"vmwareSiteId,omitempty"`

	// Fields inherited from FabricSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s VMwareV2FabricSpecificDetails) FabricSpecificDetails() BaseFabricSpecificDetailsImpl {
	return BaseFabricSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = VMwareV2FabricSpecificDetails{}

func (s VMwareV2FabricSpecificDetails) MarshalJSON() ([]byte, error) {
	type wrapper VMwareV2FabricSpecificDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareV2FabricSpecificDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareV2FabricSpecificDetails: %+v", err)
	}

	decoded["instanceType"] = "VMwareV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareV2FabricSpecificDetails: %+v", err)
	}

	return encoded, nil
}
