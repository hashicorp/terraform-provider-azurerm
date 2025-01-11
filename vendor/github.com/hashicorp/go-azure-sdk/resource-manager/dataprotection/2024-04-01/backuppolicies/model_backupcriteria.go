package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupCriteria interface {
	BackupCriteria() BaseBackupCriteriaImpl
}

var _ BackupCriteria = BaseBackupCriteriaImpl{}

type BaseBackupCriteriaImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseBackupCriteriaImpl) BackupCriteria() BaseBackupCriteriaImpl {
	return s
}

var _ BackupCriteria = RawBackupCriteriaImpl{}

// RawBackupCriteriaImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBackupCriteriaImpl struct {
	backupCriteria BaseBackupCriteriaImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawBackupCriteriaImpl) BackupCriteria() BaseBackupCriteriaImpl {
	return s.backupCriteria
}

func UnmarshalBackupCriteriaImplementation(input []byte) (BackupCriteria, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupCriteria into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ScheduleBasedBackupCriteria") {
		var out ScheduleBasedBackupCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScheduleBasedBackupCriteria: %+v", err)
		}
		return out, nil
	}

	var parent BaseBackupCriteriaImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBackupCriteriaImpl: %+v", err)
	}

	return RawBackupCriteriaImpl{
		backupCriteria: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
