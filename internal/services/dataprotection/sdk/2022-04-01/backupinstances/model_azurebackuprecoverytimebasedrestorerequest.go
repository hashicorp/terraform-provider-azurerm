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
