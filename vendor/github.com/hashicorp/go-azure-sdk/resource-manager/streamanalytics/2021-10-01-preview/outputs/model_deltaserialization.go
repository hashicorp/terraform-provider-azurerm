package outputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Serialization = DeltaSerialization{}

type DeltaSerialization struct {
	Properties *DeltaSerializationProperties `json:"properties,omitempty"`

	// Fields inherited from Serialization
}

var _ json.Marshaler = DeltaSerialization{}

func (s DeltaSerialization) MarshalJSON() ([]byte, error) {
	type wrapper DeltaSerialization
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeltaSerialization: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeltaSerialization: %+v", err)
	}
	decoded["type"] = "Delta"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeltaSerialization: %+v", err)
	}

	return encoded, nil
}
