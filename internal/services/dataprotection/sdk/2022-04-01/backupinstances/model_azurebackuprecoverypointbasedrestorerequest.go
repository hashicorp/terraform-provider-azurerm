package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AzureBackupRestoreRequest = AzureBackupRecoveryPointBasedRestoreRequest{}

type AzureBackupRecoveryPointBasedRestoreRequest struct {
	RecoveryPointId string `json:"recoveryPointId"`

	// Fields inherited from AzureBackupRestoreRequest
	RestoreTargetInfo   RestoreTargetInfoBase `json:"restoreTargetInfo"`
	SourceDataStoreType SourceDataStoreType   `json:"sourceDataStoreType"`
	SourceResourceId    *string               `json:"sourceResourceId,omitempty"`
}

var _ json.Marshaler = AzureBackupRecoveryPointBasedRestoreRequest{}

func (s AzureBackupRecoveryPointBasedRestoreRequest) MarshalJSON() ([]byte, error) {
	type wrapper AzureBackupRecoveryPointBasedRestoreRequest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBackupRecoveryPointBasedRestoreRequest: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupRecoveryPointBasedRestoreRequest: %+v", err)
	}
	decoded["objectType"] = "AzureBackupRecoveryPointBasedRestoreRequest"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBackupRecoveryPointBasedRestoreRequest: %+v", err)
	}

	return encoded, nil
}
