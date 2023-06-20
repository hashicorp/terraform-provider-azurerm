package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopyOption = CopyOnExpiryOption{}

type CopyOnExpiryOption struct {

	// Fields inherited from CopyOption
}

var _ json.Marshaler = CopyOnExpiryOption{}

func (s CopyOnExpiryOption) MarshalJSON() ([]byte, error) {
	type wrapper CopyOnExpiryOption
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CopyOnExpiryOption: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CopyOnExpiryOption: %+v", err)
	}
	decoded["objectType"] = "CopyOnExpiryOption"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CopyOnExpiryOption: %+v", err)
	}

	return encoded, nil
}
