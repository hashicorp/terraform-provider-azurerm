package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopyOption = CustomCopyOption{}

type CustomCopyOption struct {
	Duration *string `json:"duration,omitempty"`

	// Fields inherited from CopyOption
}

var _ json.Marshaler = CustomCopyOption{}

func (s CustomCopyOption) MarshalJSON() ([]byte, error) {
	type wrapper CustomCopyOption
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomCopyOption: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomCopyOption: %+v", err)
	}
	decoded["objectType"] = "CustomCopyOption"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomCopyOption: %+v", err)
	}

	return encoded, nil
}
