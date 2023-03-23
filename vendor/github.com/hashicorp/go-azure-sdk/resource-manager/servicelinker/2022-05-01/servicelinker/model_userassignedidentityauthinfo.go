package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthInfoBase = UserAssignedIdentityAuthInfo{}

type UserAssignedIdentityAuthInfo struct {
	ClientId       *string `json:"clientId,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`

	// Fields inherited from AuthInfoBase
}

var _ json.Marshaler = UserAssignedIdentityAuthInfo{}

func (s UserAssignedIdentityAuthInfo) MarshalJSON() ([]byte, error) {
	type wrapper UserAssignedIdentityAuthInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UserAssignedIdentityAuthInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UserAssignedIdentityAuthInfo: %+v", err)
	}
	decoded["authType"] = "userAssignedIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UserAssignedIdentityAuthInfo: %+v", err)
	}

	return encoded, nil
}
