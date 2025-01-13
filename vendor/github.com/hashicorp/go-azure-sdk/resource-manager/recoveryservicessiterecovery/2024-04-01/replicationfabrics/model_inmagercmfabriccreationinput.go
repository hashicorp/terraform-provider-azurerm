package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificCreationInput = InMageRcmFabricCreationInput{}

type InMageRcmFabricCreationInput struct {
	PhysicalSiteId      string                `json:"physicalSiteId"`
	SourceAgentIdentity IdentityProviderInput `json:"sourceAgentIdentity"`
	VMwareSiteId        string                `json:"vmwareSiteId"`

	// Fields inherited from FabricSpecificCreationInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmFabricCreationInput) FabricSpecificCreationInput() BaseFabricSpecificCreationInputImpl {
	return BaseFabricSpecificCreationInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmFabricCreationInput{}

func (s InMageRcmFabricCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmFabricCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmFabricCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmFabricCreationInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmFabricCreationInput: %+v", err)
	}

	return encoded, nil
}
