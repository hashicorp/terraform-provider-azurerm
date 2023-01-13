package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificCreationInput = AzureFabricCreationInput{}

type AzureFabricCreationInput struct {
	Location *string `json:"location,omitempty"`

	// Fields inherited from FabricSpecificCreationInput
}

var _ json.Marshaler = AzureFabricCreationInput{}

func (s AzureFabricCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper AzureFabricCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFabricCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFabricCreationInput: %+v", err)
	}
	decoded["instanceType"] = "Azure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFabricCreationInput: %+v", err)
	}

	return encoded, nil
}
