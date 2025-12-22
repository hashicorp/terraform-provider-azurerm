package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupDatasourceParameters = AdlsBlobBackupDatasourceParameters{}

type AdlsBlobBackupDatasourceParameters struct {
	ContainersList []string `json:"containersList"`

	// Fields inherited from BackupDatasourceParameters

	ObjectType string `json:"objectType"`
}

func (s AdlsBlobBackupDatasourceParameters) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return BaseBackupDatasourceParametersImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = AdlsBlobBackupDatasourceParameters{}

func (s AdlsBlobBackupDatasourceParameters) MarshalJSON() ([]byte, error) {
	type wrapper AdlsBlobBackupDatasourceParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AdlsBlobBackupDatasourceParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AdlsBlobBackupDatasourceParameters: %+v", err)
	}

	decoded["objectType"] = "AdlsBlobBackupDatasourceParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AdlsBlobBackupDatasourceParameters: %+v", err)
	}

	return encoded, nil
}
