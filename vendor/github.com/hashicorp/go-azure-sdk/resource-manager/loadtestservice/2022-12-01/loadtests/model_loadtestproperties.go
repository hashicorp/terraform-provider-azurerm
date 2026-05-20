package loadtests

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadTestProperties struct {
	DataPlaneURI      *string               `json:"dataPlaneURI,omitempty"`
	Description       *string               `json:"description,omitempty"`
	Encryption        *EncryptionProperties `json:"encryption,omitempty"`
	ProvisioningState *ResourceState        `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = LoadTestProperties{}

func (s LoadTestProperties) MarshalJSON() ([]byte, error) {
	type wrapper LoadTestProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LoadTestProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LoadTestProperties: %+v", err)
	}

	delete(decoded, "dataPlaneURI")

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LoadTestProperties: %+v", err)
	}

	return encoded, nil
}
