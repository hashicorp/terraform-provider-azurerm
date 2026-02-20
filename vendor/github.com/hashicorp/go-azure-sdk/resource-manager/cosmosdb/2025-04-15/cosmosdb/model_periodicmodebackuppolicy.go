package cosmosdb

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupPolicy = PeriodicModeBackupPolicy{}

type PeriodicModeBackupPolicy struct {
	PeriodicModeProperties *PeriodicModeProperties `json:"periodicModeProperties,omitempty"`

	// Fields inherited from BackupPolicy

	MigrationState *BackupPolicyMigrationState `json:"migrationState,omitempty"`
	Type           BackupPolicyType            `json:"type"`
}

func (s PeriodicModeBackupPolicy) BackupPolicy() BaseBackupPolicyImpl {
	return BaseBackupPolicyImpl{
		MigrationState: s.MigrationState,
		Type:           s.Type,
	}
}

var _ json.Marshaler = PeriodicModeBackupPolicy{}

func (s PeriodicModeBackupPolicy) MarshalJSON() ([]byte, error) {
	type wrapper PeriodicModeBackupPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PeriodicModeBackupPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PeriodicModeBackupPolicy: %+v", err)
	}

	decoded["type"] = "Periodic"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PeriodicModeBackupPolicy: %+v", err)
	}

	return encoded, nil
}
