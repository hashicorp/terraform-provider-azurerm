package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CustomPersistentDiskProperties = AzureFileVolume{}

type AzureFileVolume struct {
	ShareName *string `json:"shareName,omitempty"`

	// Fields inherited from CustomPersistentDiskProperties

	EnableSubPath *bool     `json:"enableSubPath,omitempty"`
	MountOptions  *[]string `json:"mountOptions,omitempty"`
	MountPath     string    `json:"mountPath"`
	ReadOnly      *bool     `json:"readOnly,omitempty"`
	Type          Type      `json:"type"`
}

func (s AzureFileVolume) CustomPersistentDiskProperties() BaseCustomPersistentDiskPropertiesImpl {
	return BaseCustomPersistentDiskPropertiesImpl{
		EnableSubPath: s.EnableSubPath,
		MountOptions:  s.MountOptions,
		MountPath:     s.MountPath,
		ReadOnly:      s.ReadOnly,
		Type:          s.Type,
	}
}

var _ json.Marshaler = AzureFileVolume{}

func (s AzureFileVolume) MarshalJSON() ([]byte, error) {
	type wrapper AzureFileVolume
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFileVolume: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFileVolume: %+v", err)
	}

	decoded["type"] = "AzureFileVolume"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFileVolume: %+v", err)
	}

	return encoded, nil
}
