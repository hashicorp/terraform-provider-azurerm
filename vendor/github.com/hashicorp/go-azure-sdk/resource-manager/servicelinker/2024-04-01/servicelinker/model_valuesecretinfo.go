package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretInfoBase = ValueSecretInfo{}

type ValueSecretInfo struct {
	Value *string `json:"value,omitempty"`

	// Fields inherited from SecretInfoBase

	SecretType SecretType `json:"secretType"`
}

func (s ValueSecretInfo) SecretInfoBase() BaseSecretInfoBaseImpl {
	return BaseSecretInfoBaseImpl{
		SecretType: s.SecretType,
	}
}

var _ json.Marshaler = ValueSecretInfo{}

func (s ValueSecretInfo) MarshalJSON() ([]byte, error) {
	type wrapper ValueSecretInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ValueSecretInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ValueSecretInfo: %+v", err)
	}

	decoded["secretType"] = "rawValue"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ValueSecretInfo: %+v", err)
	}

	return encoded, nil
}
