package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificCreationInput = VMwareV2FabricCreationInput{}

type VMwareV2FabricCreationInput struct {
	MigrationSolutionId string  `json:"migrationSolutionId"`
	PhysicalSiteId      *string `json:"physicalSiteId,omitempty"`
	VMwareSiteId        *string `json:"vmwareSiteId,omitempty"`

	// Fields inherited from FabricSpecificCreationInput
}

var _ json.Marshaler = VMwareV2FabricCreationInput{}

func (s VMwareV2FabricCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper VMwareV2FabricCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareV2FabricCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareV2FabricCreationInput: %+v", err)
	}
	decoded["instanceType"] = "VMwareV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareV2FabricCreationInput: %+v", err)
	}

	return encoded, nil
}
