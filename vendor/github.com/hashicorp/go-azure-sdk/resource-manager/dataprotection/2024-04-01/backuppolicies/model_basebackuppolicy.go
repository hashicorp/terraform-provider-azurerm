package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseBackupPolicy interface {
	BaseBackupPolicy() BaseBaseBackupPolicyImpl
}

var _ BaseBackupPolicy = BaseBaseBackupPolicyImpl{}

type BaseBaseBackupPolicyImpl struct {
	DatasourceTypes []string `json:"datasourceTypes"`
	ObjectType      string   `json:"objectType"`
}

func (s BaseBaseBackupPolicyImpl) BaseBackupPolicy() BaseBaseBackupPolicyImpl {
	return s
}

var _ BaseBackupPolicy = RawBaseBackupPolicyImpl{}

// RawBaseBackupPolicyImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBaseBackupPolicyImpl struct {
	baseBackupPolicy BaseBaseBackupPolicyImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawBaseBackupPolicyImpl) BaseBackupPolicy() BaseBaseBackupPolicyImpl {
	return s.baseBackupPolicy
}

func UnmarshalBaseBackupPolicyImplementation(input []byte) (BaseBackupPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BaseBackupPolicy into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BackupPolicy") {
		var out BackupPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BackupPolicy: %+v", err)
		}
		return out, nil
	}

	var parent BaseBaseBackupPolicyImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBaseBackupPolicyImpl: %+v", err)
	}

	return RawBaseBackupPolicyImpl{
		baseBackupPolicy: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
