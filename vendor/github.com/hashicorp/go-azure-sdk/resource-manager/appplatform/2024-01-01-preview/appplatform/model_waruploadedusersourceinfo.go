package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UserSourceInfo = WarUploadedUserSourceInfo{}

type WarUploadedUserSourceInfo struct {
	JVMOptions     *string `json:"jvmOptions,omitempty"`
	RelativePath   *string `json:"relativePath,omitempty"`
	RuntimeVersion *string `json:"runtimeVersion,omitempty"`
	ServerVersion  *string `json:"serverVersion,omitempty"`

	// Fields inherited from UserSourceInfo

	Type    string  `json:"type"`
	Version *string `json:"version,omitempty"`
}

func (s WarUploadedUserSourceInfo) UserSourceInfo() BaseUserSourceInfoImpl {
	return BaseUserSourceInfoImpl{
		Type:    s.Type,
		Version: s.Version,
	}
}

var _ json.Marshaler = WarUploadedUserSourceInfo{}

func (s WarUploadedUserSourceInfo) MarshalJSON() ([]byte, error) {
	type wrapper WarUploadedUserSourceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WarUploadedUserSourceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WarUploadedUserSourceInfo: %+v", err)
	}

	decoded["type"] = "War"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WarUploadedUserSourceInfo: %+v", err)
	}

	return encoded, nil
}
