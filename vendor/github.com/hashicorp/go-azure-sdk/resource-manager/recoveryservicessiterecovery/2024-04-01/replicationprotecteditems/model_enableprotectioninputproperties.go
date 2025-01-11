package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnableProtectionInputProperties struct {
	PolicyId                *string                               `json:"policyId,omitempty"`
	ProtectableItemId       *string                               `json:"protectableItemId,omitempty"`
	ProviderSpecificDetails EnableProtectionProviderSpecificInput `json:"providerSpecificDetails"`
}

var _ json.Unmarshaler = &EnableProtectionInputProperties{}

func (s *EnableProtectionInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		PolicyId          *string `json:"policyId,omitempty"`
		ProtectableItemId *string `json:"protectableItemId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.PolicyId = decoded.PolicyId
	s.ProtectableItemId = decoded.ProtectableItemId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EnableProtectionInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalEnableProtectionProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'EnableProtectionInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
