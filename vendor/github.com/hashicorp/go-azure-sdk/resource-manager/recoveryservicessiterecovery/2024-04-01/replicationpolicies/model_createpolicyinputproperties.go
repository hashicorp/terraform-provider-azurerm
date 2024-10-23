package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreatePolicyInputProperties struct {
	ProviderSpecificInput PolicyProviderSpecificInput `json:"providerSpecificInput"`
}

var _ json.Unmarshaler = &CreatePolicyInputProperties{}

func (s *CreatePolicyInputProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CreatePolicyInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificInput"]; ok {
		impl, err := UnmarshalPolicyProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificInput' for 'CreatePolicyInputProperties': %+v", err)
		}
		s.ProviderSpecificInput = impl
	}

	return nil
}
