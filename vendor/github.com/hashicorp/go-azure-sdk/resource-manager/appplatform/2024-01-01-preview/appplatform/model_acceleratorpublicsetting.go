package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AcceleratorAuthSetting = AcceleratorPublicSetting{}

type AcceleratorPublicSetting struct {
	CaCertResourceId *string `json:"caCertResourceId,omitempty"`

	// Fields inherited from AcceleratorAuthSetting

	AuthType string `json:"authType"`
}

func (s AcceleratorPublicSetting) AcceleratorAuthSetting() BaseAcceleratorAuthSettingImpl {
	return BaseAcceleratorAuthSettingImpl{
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = AcceleratorPublicSetting{}

func (s AcceleratorPublicSetting) MarshalJSON() ([]byte, error) {
	type wrapper AcceleratorPublicSetting
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AcceleratorPublicSetting: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AcceleratorPublicSetting: %+v", err)
	}

	decoded["authType"] = "Public"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AcceleratorPublicSetting: %+v", err)
	}

	return encoded, nil
}
