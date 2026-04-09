package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthInfoBase = AccessKeyInfoBase{}

type AccessKeyInfoBase struct {
	Permissions *[]AccessKeyPermissions `json:"permissions,omitempty"`

	// Fields inherited from AuthInfoBase

	AuthMode *AuthMode `json:"authMode,omitempty"`
	AuthType AuthType  `json:"authType"`
}

func (s AccessKeyInfoBase) AuthInfoBase() BaseAuthInfoBaseImpl {
	return BaseAuthInfoBaseImpl{
		AuthMode: s.AuthMode,
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = AccessKeyInfoBase{}

func (s AccessKeyInfoBase) MarshalJSON() ([]byte, error) {
	type wrapper AccessKeyInfoBase
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AccessKeyInfoBase: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AccessKeyInfoBase: %+v", err)
	}

	decoded["authType"] = "accessKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AccessKeyInfoBase: %+v", err)
	}

	return encoded, nil
}
