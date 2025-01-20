package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UserSourceInfo = JarUploadedUserSourceInfo{}

type JarUploadedUserSourceInfo struct {
	JVMOptions     *string `json:"jvmOptions,omitempty"`
	RelativePath   *string `json:"relativePath,omitempty"`
	RuntimeVersion *string `json:"runtimeVersion,omitempty"`

	// Fields inherited from UserSourceInfo

	Type    string  `json:"type"`
	Version *string `json:"version,omitempty"`
}

func (s JarUploadedUserSourceInfo) UserSourceInfo() BaseUserSourceInfoImpl {
	return BaseUserSourceInfoImpl{
		Type:    s.Type,
		Version: s.Version,
	}
}

var _ json.Marshaler = JarUploadedUserSourceInfo{}

func (s JarUploadedUserSourceInfo) MarshalJSON() ([]byte, error) {
	type wrapper JarUploadedUserSourceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JarUploadedUserSourceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JarUploadedUserSourceInfo: %+v", err)
	}

	decoded["type"] = "Jar"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JarUploadedUserSourceInfo: %+v", err)
	}

	return encoded, nil
}
