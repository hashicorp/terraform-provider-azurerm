package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplyRecoveryPointInputProperties struct {
	ProviderSpecificDetails ApplyRecoveryPointProviderSpecificInput `json:"providerSpecificDetails"`
	RecoveryPointId         *string                                 `json:"recoveryPointId,omitempty"`
}

var _ json.Unmarshaler = &ApplyRecoveryPointInputProperties{}

func (s *ApplyRecoveryPointInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		RecoveryPointId *string `json:"recoveryPointId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.RecoveryPointId = decoded.RecoveryPointId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ApplyRecoveryPointInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalApplyRecoveryPointProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'ApplyRecoveryPointInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
