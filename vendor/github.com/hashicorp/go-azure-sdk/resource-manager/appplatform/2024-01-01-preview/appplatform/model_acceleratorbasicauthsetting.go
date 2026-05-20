package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AcceleratorAuthSetting = AcceleratorBasicAuthSetting{}

type AcceleratorBasicAuthSetting struct {
	CaCertResourceId *string `json:"caCertResourceId,omitempty"`
	Password         *string `json:"password,omitempty"`
	Username         string  `json:"username"`

	// Fields inherited from AcceleratorAuthSetting

	AuthType string `json:"authType"`
}

func (s AcceleratorBasicAuthSetting) AcceleratorAuthSetting() BaseAcceleratorAuthSettingImpl {
	return BaseAcceleratorAuthSettingImpl{
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = AcceleratorBasicAuthSetting{}

func (s AcceleratorBasicAuthSetting) MarshalJSON() ([]byte, error) {
	type wrapper AcceleratorBasicAuthSetting
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AcceleratorBasicAuthSetting: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AcceleratorBasicAuthSetting: %+v", err)
	}

	decoded["authType"] = "BasicAuth"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AcceleratorBasicAuthSetting: %+v", err)
	}

	return encoded, nil
}
