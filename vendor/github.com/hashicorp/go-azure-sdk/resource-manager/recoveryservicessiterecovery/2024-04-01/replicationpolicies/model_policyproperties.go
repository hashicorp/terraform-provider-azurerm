package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyProperties struct {
	FriendlyName            *string                       `json:"friendlyName,omitempty"`
	ProviderSpecificDetails PolicyProviderSpecificDetails `json:"providerSpecificDetails"`
}

var _ json.Unmarshaler = &PolicyProperties{}

func (s *PolicyProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FriendlyName *string `json:"friendlyName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FriendlyName = decoded.FriendlyName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PolicyProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalPolicyProviderSpecificDetailsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'PolicyProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
