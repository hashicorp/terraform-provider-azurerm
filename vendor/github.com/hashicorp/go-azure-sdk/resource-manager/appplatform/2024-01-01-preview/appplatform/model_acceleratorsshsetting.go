package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AcceleratorAuthSetting = AcceleratorSshSetting{}

type AcceleratorSshSetting struct {
	HostKey          *string `json:"hostKey,omitempty"`
	HostKeyAlgorithm *string `json:"hostKeyAlgorithm,omitempty"`
	PrivateKey       *string `json:"privateKey,omitempty"`

	// Fields inherited from AcceleratorAuthSetting

	AuthType string `json:"authType"`
}

func (s AcceleratorSshSetting) AcceleratorAuthSetting() BaseAcceleratorAuthSettingImpl {
	return BaseAcceleratorAuthSettingImpl{
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = AcceleratorSshSetting{}

func (s AcceleratorSshSetting) MarshalJSON() ([]byte, error) {
	type wrapper AcceleratorSshSetting
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AcceleratorSshSetting: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AcceleratorSshSetting: %+v", err)
	}

	decoded["authType"] = "SSH"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AcceleratorSshSetting: %+v", err)
	}

	return encoded, nil
}
