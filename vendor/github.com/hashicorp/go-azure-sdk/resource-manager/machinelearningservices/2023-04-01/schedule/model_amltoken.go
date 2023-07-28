package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ IdentityConfiguration = AmlToken{}

type AmlToken struct {

	// Fields inherited from IdentityConfiguration
}

var _ json.Marshaler = AmlToken{}

func (s AmlToken) MarshalJSON() ([]byte, error) {
	type wrapper AmlToken
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AmlToken: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AmlToken: %+v", err)
	}
	decoded["identityType"] = "AMLToken"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AmlToken: %+v", err)
	}

	return encoded, nil
}
