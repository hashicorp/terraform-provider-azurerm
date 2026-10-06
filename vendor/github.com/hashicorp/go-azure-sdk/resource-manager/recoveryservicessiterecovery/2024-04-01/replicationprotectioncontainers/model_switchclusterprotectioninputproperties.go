package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchClusterProtectionInputProperties struct {
	ProviderSpecificDetails          SwitchClusterProtectionProviderSpecificInput `json:"providerSpecificDetails"`
	ReplicationProtectionClusterName *string                                      `json:"replicationProtectionClusterName,omitempty"`
}

var _ json.Unmarshaler = &SwitchClusterProtectionInputProperties{}

func (s *SwitchClusterProtectionInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ReplicationProtectionClusterName *string `json:"replicationProtectionClusterName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ReplicationProtectionClusterName = decoded.ReplicationProtectionClusterName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SwitchClusterProtectionInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalSwitchClusterProtectionProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'SwitchClusterProtectionInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
