package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupDatasourceParameters interface {
	BackupDatasourceParameters() BaseBackupDatasourceParametersImpl
}

var _ BackupDatasourceParameters = BaseBackupDatasourceParametersImpl{}

type BaseBackupDatasourceParametersImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseBackupDatasourceParametersImpl) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return s
}

var _ BackupDatasourceParameters = RawBackupDatasourceParametersImpl{}

// RawBackupDatasourceParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBackupDatasourceParametersImpl struct {
	backupDatasourceParameters BaseBackupDatasourceParametersImpl
	Type                       string
	Values                     map[string]interface{}
}

func (s RawBackupDatasourceParametersImpl) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return s.backupDatasourceParameters
}

func UnmarshalBackupDatasourceParametersImplementation(input []byte) (BackupDatasourceParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupDatasourceParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BlobBackupDatasourceParameters") {
		var out BlobBackupDatasourceParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobBackupDatasourceParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KubernetesClusterBackupDatasourceParameters") {
		var out KubernetesClusterBackupDatasourceParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KubernetesClusterBackupDatasourceParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseBackupDatasourceParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBackupDatasourceParametersImpl: %+v", err)
	}

	return RawBackupDatasourceParametersImpl{
		backupDatasourceParameters: parent,
		Type:                       value,
		Values:                     temp,
	}, nil

}
