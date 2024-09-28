package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupParameters interface {
	BackupParameters() BaseBackupParametersImpl
}

var _ BackupParameters = BaseBackupParametersImpl{}

type BaseBackupParametersImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseBackupParametersImpl) BackupParameters() BaseBackupParametersImpl {
	return s
}

var _ BackupParameters = RawBackupParametersImpl{}

// RawBackupParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBackupParametersImpl struct {
	backupParameters BaseBackupParametersImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawBackupParametersImpl) BackupParameters() BaseBackupParametersImpl {
	return s.backupParameters
}

func UnmarshalBackupParametersImplementation(input []byte) (BackupParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureBackupParams") {
		var out AzureBackupParams
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupParams: %+v", err)
		}
		return out, nil
	}

	var parent BaseBackupParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBackupParametersImpl: %+v", err)
	}

	return RawBackupParametersImpl{
		backupParameters: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
