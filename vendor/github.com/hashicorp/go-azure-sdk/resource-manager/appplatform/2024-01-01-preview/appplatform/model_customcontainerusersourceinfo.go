package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UserSourceInfo = CustomContainerUserSourceInfo{}

type CustomContainerUserSourceInfo struct {
	CustomContainer *CustomContainer `json:"customContainer,omitempty"`

	// Fields inherited from UserSourceInfo

	Type    string  `json:"type"`
	Version *string `json:"version,omitempty"`
}

func (s CustomContainerUserSourceInfo) UserSourceInfo() BaseUserSourceInfoImpl {
	return BaseUserSourceInfoImpl{
		Type:    s.Type,
		Version: s.Version,
	}
}

var _ json.Marshaler = CustomContainerUserSourceInfo{}

func (s CustomContainerUserSourceInfo) MarshalJSON() ([]byte, error) {
	type wrapper CustomContainerUserSourceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomContainerUserSourceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomContainerUserSourceInfo: %+v", err)
	}

	decoded["type"] = "Container"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomContainerUserSourceInfo: %+v", err)
	}

	return encoded, nil
}
