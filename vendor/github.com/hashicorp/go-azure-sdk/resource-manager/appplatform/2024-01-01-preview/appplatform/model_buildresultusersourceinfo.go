package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UserSourceInfo = BuildResultUserSourceInfo{}

type BuildResultUserSourceInfo struct {
	BuildResultId *string `json:"buildResultId,omitempty"`

	// Fields inherited from UserSourceInfo

	Type    string  `json:"type"`
	Version *string `json:"version,omitempty"`
}

func (s BuildResultUserSourceInfo) UserSourceInfo() BaseUserSourceInfoImpl {
	return BaseUserSourceInfoImpl{
		Type:    s.Type,
		Version: s.Version,
	}
}

var _ json.Marshaler = BuildResultUserSourceInfo{}

func (s BuildResultUserSourceInfo) MarshalJSON() ([]byte, error) {
	type wrapper BuildResultUserSourceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BuildResultUserSourceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BuildResultUserSourceInfo: %+v", err)
	}

	decoded["type"] = "BuildResult"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BuildResultUserSourceInfo: %+v", err)
	}

	return encoded, nil
}
