package videoanalyzer

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenKey = RsaTokenKey{}

type RsaTokenKey struct {
	Alg AccessPolicyRsaAlgo `json:"alg"`
	E   string              `json:"e"`
	N   string              `json:"n"`

	// Fields inherited from TokenKey
	Kid string `json:"kid"`
}

var _ json.Marshaler = RsaTokenKey{}

func (s RsaTokenKey) MarshalJSON() ([]byte, error) {
	type wrapper RsaTokenKey
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RsaTokenKey: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RsaTokenKey: %+v", err)
	}
	decoded["@type"] = "#Microsoft.VideoAnalyzer.RsaTokenKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RsaTokenKey: %+v", err)
	}

	return encoded, nil
}
