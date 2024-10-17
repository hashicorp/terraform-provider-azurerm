package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetServiceBase = SelfHostedServer{}

type SelfHostedServer struct {
	Endpoint *string `json:"endpoint,omitempty"`

	// Fields inherited from TargetServiceBase

	Type TargetServiceType `json:"type"`
}

func (s SelfHostedServer) TargetServiceBase() BaseTargetServiceBaseImpl {
	return BaseTargetServiceBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = SelfHostedServer{}

func (s SelfHostedServer) MarshalJSON() ([]byte, error) {
	type wrapper SelfHostedServer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SelfHostedServer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SelfHostedServer: %+v", err)
	}

	decoded["type"] = "SelfHostedServer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SelfHostedServer: %+v", err)
	}

	return encoded, nil
}
