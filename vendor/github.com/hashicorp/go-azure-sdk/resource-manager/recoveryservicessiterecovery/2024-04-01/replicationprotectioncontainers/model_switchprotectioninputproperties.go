package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchProtectionInputProperties struct {
	ProviderSpecificDetails      SwitchProtectionProviderSpecificInput `json:"providerSpecificDetails"`
	ReplicationProtectedItemName *string                               `json:"replicationProtectedItemName,omitempty"`
}

var _ json.Unmarshaler = &SwitchProtectionInputProperties{}

func (s *SwitchProtectionInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ReplicationProtectedItemName *string `json:"replicationProtectedItemName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ReplicationProtectedItemName = decoded.ReplicationProtectedItemName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SwitchProtectionInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalSwitchProtectionProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'SwitchProtectionInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
