package links

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetServiceBase = ConfluentBootstrapServer{}

type ConfluentBootstrapServer struct {
	Endpoint *string `json:"endpoint,omitempty"`

	// Fields inherited from TargetServiceBase

	Type TargetServiceType `json:"type"`
}

func (s ConfluentBootstrapServer) TargetServiceBase() BaseTargetServiceBaseImpl {
	return BaseTargetServiceBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ConfluentBootstrapServer{}

func (s ConfluentBootstrapServer) MarshalJSON() ([]byte, error) {
	type wrapper ConfluentBootstrapServer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ConfluentBootstrapServer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ConfluentBootstrapServer: %+v", err)
	}

	decoded["type"] = "ConfluentBootstrapServer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ConfluentBootstrapServer: %+v", err)
	}

	return encoded, nil
}
