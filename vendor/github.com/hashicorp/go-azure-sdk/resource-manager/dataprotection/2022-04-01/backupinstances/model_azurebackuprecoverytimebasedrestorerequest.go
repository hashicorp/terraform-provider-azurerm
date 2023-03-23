package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AzureBackupRestoreRequest = AzureBackupRecoveryTimeBasedRestoreRequest{}

type AzureBackupRecoveryTimeBasedRestoreRequest struct {
	RecoveryPointTime string `json:"recoveryPointTime"`

	// Fields inherited from AzureBackupRestoreRequest
	RestoreTargetInfo   RestoreTargetInfoBase `json:"restoreTargetInfo"`
	SourceDataStoreType SourceDataStoreType   `json:"sourceDataStoreType"`
	SourceResourceId    *string               `json:"sourceResourceId,omitempty"`
}

var _ json.Marshaler = AzureBackupRecoveryTimeBasedRestoreRequest{}

func (s AzureBackupRecoveryTimeBasedRestoreRequest) MarshalJSON() ([]byte, error) {
	type wrapper AzureBackupRecoveryTimeBasedRestoreRequest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBackupRecoveryTimeBasedRestoreRequest: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupRecoveryTimeBasedRestoreRequest: %+v", err)
	}
	decoded["objectType"] = "AzureBackupRecoveryTimeBasedRestoreRequest"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBackupRecoveryTimeBasedRestoreRequest: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureBackupRecoveryTimeBasedRestoreRequest{}

func (s *AzureBackupRecoveryTimeBasedRestoreRequest) UnmarshalJSON(bytes []byte) error {
	type alias AzureBackupRecoveryTimeBasedRestoreRequest
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AzureBackupRecoveryTimeBasedRestoreRequest: %+v", err)
	}

	s.RecoveryPointTime = decoded.RecoveryPointTime
	s.SourceDataStoreType = decoded.SourceDataStoreType
	s.SourceResourceId = decoded.SourceResourceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureBackupRecoveryTimeBasedRestoreRequest into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["restoreTargetInfo"]; ok {
		impl, err := unmarshalRestoreTargetInfoBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RestoreTargetInfo' for 'AzureBackupRecoveryTimeBasedRestoreRequest': %+v", err)
		}
		s.RestoreTargetInfo = impl
	}
	return nil
}
