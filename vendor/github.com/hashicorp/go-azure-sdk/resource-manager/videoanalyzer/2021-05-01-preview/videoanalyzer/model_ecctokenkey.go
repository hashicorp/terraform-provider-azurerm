package videoanalyzer

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenKey = EccTokenKey{}

type EccTokenKey struct {
	Alg AccessPolicyEccAlgo `json:"alg"`
	X   string              `json:"x"`
	Y   string              `json:"y"`

	// Fields inherited from TokenKey
	Kid string `json:"kid"`
}

var _ json.Marshaler = EccTokenKey{}

func (s EccTokenKey) MarshalJSON() ([]byte, error) {
	type wrapper EccTokenKey
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EccTokenKey: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EccTokenKey: %+v", err)
	}
	decoded["@type"] = "#Microsoft.VideoAnalyzer.EccTokenKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EccTokenKey: %+v", err)
	}

	return encoded, nil
}
