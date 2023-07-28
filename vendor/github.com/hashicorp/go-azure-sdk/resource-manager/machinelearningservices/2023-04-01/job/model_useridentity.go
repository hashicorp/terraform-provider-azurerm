package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ IdentityConfiguration = UserIdentity{}

type UserIdentity struct {

	// Fields inherited from IdentityConfiguration
}

var _ json.Marshaler = UserIdentity{}

func (s UserIdentity) MarshalJSON() ([]byte, error) {
	type wrapper UserIdentity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UserIdentity: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UserIdentity: %+v", err)
	}
	decoded["identityType"] = "UserIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UserIdentity: %+v", err)
	}

	return encoded, nil
}
