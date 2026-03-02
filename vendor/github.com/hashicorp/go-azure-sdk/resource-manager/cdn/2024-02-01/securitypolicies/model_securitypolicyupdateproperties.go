package securitypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityPolicyUpdateProperties struct {
	Parameters SecurityPolicyPropertiesParameters `json:"parameters"`
}

var _ json.Unmarshaler = &SecurityPolicyUpdateProperties{}

func (s *SecurityPolicyUpdateProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SecurityPolicyUpdateProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["parameters"]; ok {
		impl, err := UnmarshalSecurityPolicyPropertiesParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Parameters' for 'SecurityPolicyUpdateProperties': %+v", err)
		}
		s.Parameters = impl
	}

	return nil
}
