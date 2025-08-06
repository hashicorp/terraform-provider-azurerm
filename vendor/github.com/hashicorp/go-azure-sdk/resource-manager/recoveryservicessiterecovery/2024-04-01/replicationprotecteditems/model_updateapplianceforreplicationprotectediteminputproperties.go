package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateApplianceForReplicationProtectedItemInputProperties struct {
	ProviderSpecificDetails UpdateApplianceForReplicationProtectedItemProviderSpecificInput `json:"providerSpecificDetails"`
	TargetApplianceId       string                                                          `json:"targetApplianceId"`
}

var _ json.Unmarshaler = &UpdateApplianceForReplicationProtectedItemInputProperties{}

func (s *UpdateApplianceForReplicationProtectedItemInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		TargetApplianceId string `json:"targetApplianceId"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.TargetApplianceId = decoded.TargetApplianceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UpdateApplianceForReplicationProtectedItemInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalUpdateApplianceForReplicationProtectedItemProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'UpdateApplianceForReplicationProtectedItemInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
