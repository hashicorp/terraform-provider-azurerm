package projects

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectProperties struct {
	Description        *string            `json:"description,omitempty"`
	DevCenterId        string             `json:"devCenterId"`
	DevCenterUri       *string            `json:"devCenterUri,omitempty"`
	MaxDevBoxesPerUser *int64             `json:"maxDevBoxesPerUser,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = ProjectProperties{}

func (s ProjectProperties) MarshalJSON() ([]byte, error) {
	type wrapper ProjectProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ProjectProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ProjectProperties: %+v", err)
	}

	delete(decoded, "devCenterUri")

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ProjectProperties: %+v", err)
	}

	return encoded, nil
}
