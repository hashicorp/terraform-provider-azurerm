package cosmosdb

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupPolicy interface {
	BackupPolicy() BaseBackupPolicyImpl
}

var _ BackupPolicy = BaseBackupPolicyImpl{}

type BaseBackupPolicyImpl struct {
	MigrationState *BackupPolicyMigrationState `json:"migrationState,omitempty"`
	Type           BackupPolicyType            `json:"type"`
}

func (s BaseBackupPolicyImpl) BackupPolicy() BaseBackupPolicyImpl {
	return s
}

var _ BackupPolicy = RawBackupPolicyImpl{}

// RawBackupPolicyImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBackupPolicyImpl struct {
	backupPolicy BaseBackupPolicyImpl
	Type         string
	Values       map[string]interface{}
}

func (s RawBackupPolicyImpl) BackupPolicy() BaseBackupPolicyImpl {
	return s.backupPolicy
}

func UnmarshalBackupPolicyImplementation(input []byte) (BackupPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupPolicy into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Continuous") {
		var out ContinuousModeBackupPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContinuousModeBackupPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Periodic") {
		var out PeriodicModeBackupPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PeriodicModeBackupPolicy: %+v", err)
		}
		return out, nil
	}

	var parent BaseBackupPolicyImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBackupPolicyImpl: %+v", err)
	}

	return RawBackupPolicyImpl{
		backupPolicy: parent,
		Type:         value,
		Values:       temp,
	}, nil

}
