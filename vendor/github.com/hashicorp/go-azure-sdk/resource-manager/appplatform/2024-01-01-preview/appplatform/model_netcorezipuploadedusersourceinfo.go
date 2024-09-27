package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UserSourceInfo = NetCoreZipUploadedUserSourceInfo{}

type NetCoreZipUploadedUserSourceInfo struct {
	NetCoreMainEntryPath *string `json:"netCoreMainEntryPath,omitempty"`
	RelativePath         *string `json:"relativePath,omitempty"`
	RuntimeVersion       *string `json:"runtimeVersion,omitempty"`

	// Fields inherited from UserSourceInfo

	Type    string  `json:"type"`
	Version *string `json:"version,omitempty"`
}

func (s NetCoreZipUploadedUserSourceInfo) UserSourceInfo() BaseUserSourceInfoImpl {
	return BaseUserSourceInfoImpl{
		Type:    s.Type,
		Version: s.Version,
	}
}

var _ json.Marshaler = NetCoreZipUploadedUserSourceInfo{}

func (s NetCoreZipUploadedUserSourceInfo) MarshalJSON() ([]byte, error) {
	type wrapper NetCoreZipUploadedUserSourceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NetCoreZipUploadedUserSourceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NetCoreZipUploadedUserSourceInfo: %+v", err)
	}

	decoded["type"] = "NetCoreZip"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NetCoreZipUploadedUserSourceInfo: %+v", err)
	}

	return encoded, nil
}
