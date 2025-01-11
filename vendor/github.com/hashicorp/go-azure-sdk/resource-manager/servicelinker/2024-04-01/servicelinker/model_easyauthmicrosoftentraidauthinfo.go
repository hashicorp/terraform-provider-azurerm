package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthInfoBase = EasyAuthMicrosoftEntraIDAuthInfo{}

type EasyAuthMicrosoftEntraIDAuthInfo struct {
	ClientId               *string                 `json:"clientId,omitempty"`
	DeleteOrUpdateBehavior *DeleteOrUpdateBehavior `json:"deleteOrUpdateBehavior,omitempty"`
	Secret                 *string                 `json:"secret,omitempty"`

	// Fields inherited from AuthInfoBase

	AuthMode *AuthMode `json:"authMode,omitempty"`
	AuthType AuthType  `json:"authType"`
}

func (s EasyAuthMicrosoftEntraIDAuthInfo) AuthInfoBase() BaseAuthInfoBaseImpl {
	return BaseAuthInfoBaseImpl{
		AuthMode: s.AuthMode,
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = EasyAuthMicrosoftEntraIDAuthInfo{}

func (s EasyAuthMicrosoftEntraIDAuthInfo) MarshalJSON() ([]byte, error) {
	type wrapper EasyAuthMicrosoftEntraIDAuthInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EasyAuthMicrosoftEntraIDAuthInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EasyAuthMicrosoftEntraIDAuthInfo: %+v", err)
	}

	decoded["authType"] = "easyAuthMicrosoftEntraID"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EasyAuthMicrosoftEntraIDAuthInfo: %+v", err)
	}

	return encoded, nil
}
