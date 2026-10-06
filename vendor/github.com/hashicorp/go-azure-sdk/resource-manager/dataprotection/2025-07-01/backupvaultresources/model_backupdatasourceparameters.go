package backupvaultresources

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

// RawBackupDatasourceParametersImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawBackupDatasourceParametersImpl struct {
	backupDatasourceParameters BaseBackupDatasourceParametersImpl
	Type                       string
	Values                     map[string]interface{}
}

func (s RawBackupDatasourceParametersImpl) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return s.backupDatasourceParameters
}

func (s RawBackupDatasourceParametersImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
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

	if strings.EqualFold(value, "AdlsBlobBackupDatasourceParameters") {
		var out AdlsBlobBackupDatasourceParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AdlsBlobBackupDatasourceParameters: %+v", err)
		}
		return out, nil
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
