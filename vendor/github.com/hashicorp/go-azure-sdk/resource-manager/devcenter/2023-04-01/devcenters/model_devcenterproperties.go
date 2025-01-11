package devcenters

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevCenterProperties struct {
	DevCenterUri      *string            `json:"devCenterUri,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = DevCenterProperties{}

func (s DevCenterProperties) MarshalJSON() ([]byte, error) {
	type wrapper DevCenterProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DevCenterProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DevCenterProperties: %+v", err)
	}

	delete(decoded, "devCenterUri")

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DevCenterProperties: %+v", err)
	}

	return encoded, nil
}
