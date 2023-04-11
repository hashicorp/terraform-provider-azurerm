package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyOption struct {
	Configuration  ContentKeyPolicyConfiguration `json:"configuration"`
	Name           *string                       `json:"name,omitempty"`
	PolicyOptionId *string                       `json:"policyOptionId,omitempty"`
	Restriction    ContentKeyPolicyRestriction   `json:"restriction"`
}

var _ json.Unmarshaler = &ContentKeyPolicyOption{}

func (s *ContentKeyPolicyOption) UnmarshalJSON(bytes []byte) error {
	type alias ContentKeyPolicyOption
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ContentKeyPolicyOption: %+v", err)
	}

	s.Name = decoded.Name
	s.PolicyOptionId = decoded.PolicyOptionId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ContentKeyPolicyOption into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["configuration"]; ok {
		impl, err := unmarshalContentKeyPolicyConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Configuration' for 'ContentKeyPolicyOption': %+v", err)
		}
		s.Configuration = impl
	}

	if v, ok := temp["restriction"]; ok {
		impl, err := unmarshalContentKeyPolicyRestrictionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Restriction' for 'ContentKeyPolicyOption': %+v", err)
		}
		s.Restriction = impl
	}
	return nil
}
