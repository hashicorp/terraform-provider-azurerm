package cosmosdb

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupPolicy = ContinuousModeBackupPolicy{}

type ContinuousModeBackupPolicy struct {

	// Fields inherited from BackupPolicy

	MigrationState *BackupPolicyMigrationState `json:"migrationState,omitempty"`
	Type           BackupPolicyType            `json:"type"`
}

func (s ContinuousModeBackupPolicy) BackupPolicy() BaseBackupPolicyImpl {
	return BaseBackupPolicyImpl{
		MigrationState: s.MigrationState,
		Type:           s.Type,
	}
}

var _ json.Marshaler = ContinuousModeBackupPolicy{}

func (s ContinuousModeBackupPolicy) MarshalJSON() ([]byte, error) {
	type wrapper ContinuousModeBackupPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContinuousModeBackupPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContinuousModeBackupPolicy: %+v", err)
	}

	decoded["type"] = "Continuous"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContinuousModeBackupPolicy: %+v", err)
	}

	return encoded, nil
}
