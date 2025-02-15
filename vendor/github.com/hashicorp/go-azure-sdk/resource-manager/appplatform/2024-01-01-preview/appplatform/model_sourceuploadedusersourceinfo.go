package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UserSourceInfo = SourceUploadedUserSourceInfo{}

type SourceUploadedUserSourceInfo struct {
	ArtifactSelector *string `json:"artifactSelector,omitempty"`
	RelativePath     *string `json:"relativePath,omitempty"`
	RuntimeVersion   *string `json:"runtimeVersion,omitempty"`

	// Fields inherited from UserSourceInfo

	Type    string  `json:"type"`
	Version *string `json:"version,omitempty"`
}

func (s SourceUploadedUserSourceInfo) UserSourceInfo() BaseUserSourceInfoImpl {
	return BaseUserSourceInfoImpl{
		Type:    s.Type,
		Version: s.Version,
	}
}

var _ json.Marshaler = SourceUploadedUserSourceInfo{}

func (s SourceUploadedUserSourceInfo) MarshalJSON() ([]byte, error) {
	type wrapper SourceUploadedUserSourceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SourceUploadedUserSourceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SourceUploadedUserSourceInfo: %+v", err)
	}

	decoded["type"] = "Source"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SourceUploadedUserSourceInfo: %+v", err)
	}

	return encoded, nil
}
