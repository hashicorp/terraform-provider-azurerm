package videoanalyzer

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyProperties struct {
	Authentication AuthenticationBase `json:"authentication"`
	Role           *AccessPolicyRole  `json:"role,omitempty"`
}

var _ json.Unmarshaler = &AccessPolicyProperties{}

func (s *AccessPolicyProperties) UnmarshalJSON(bytes []byte) error {
	type alias AccessPolicyProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AccessPolicyProperties: %+v", err)
	}

	s.Role = decoded.Role

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AccessPolicyProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authentication"]; ok {
		impl, err := unmarshalAuthenticationBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Authentication' for 'AccessPolicyProperties': %+v", err)
		}
		s.Authentication = impl
	}
	return nil
}
