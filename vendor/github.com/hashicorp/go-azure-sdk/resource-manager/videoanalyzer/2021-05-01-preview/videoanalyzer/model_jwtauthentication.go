package videoanalyzer

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthenticationBase = JwtAuthentication{}

type JwtAuthentication struct {
	Audiences *[]string     `json:"audiences,omitempty"`
	Claims    *[]TokenClaim `json:"claims,omitempty"`
	Issuers   *[]string     `json:"issuers,omitempty"`
	Keys      *[]TokenKey   `json:"keys,omitempty"`

	// Fields inherited from AuthenticationBase
}

var _ json.Marshaler = JwtAuthentication{}

func (s JwtAuthentication) MarshalJSON() ([]byte, error) {
	type wrapper JwtAuthentication
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JwtAuthentication: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JwtAuthentication: %+v", err)
	}
	decoded["@type"] = "#Microsoft.VideoAnalyzer.JwtAuthentication"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JwtAuthentication: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &JwtAuthentication{}

func (s *JwtAuthentication) UnmarshalJSON(bytes []byte) error {
	type alias JwtAuthentication
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into JwtAuthentication: %+v", err)
	}

	s.Audiences = decoded.Audiences
	s.Claims = decoded.Claims
	s.Issuers = decoded.Issuers

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JwtAuthentication into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["keys"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Keys into list []json.RawMessage: %+v", err)
		}

		output := make([]TokenKey, 0)
		for i, val := range listTemp {
			impl, err := unmarshalTokenKeyImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Keys' for 'JwtAuthentication': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Keys = &output
	}
	return nil
}
