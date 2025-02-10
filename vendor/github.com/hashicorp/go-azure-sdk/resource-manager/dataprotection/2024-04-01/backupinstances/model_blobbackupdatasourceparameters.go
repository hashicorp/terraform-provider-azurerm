package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupDatasourceParameters = BlobBackupDatasourceParameters{}

type BlobBackupDatasourceParameters struct {
	ContainersList []string `json:"containersList"`

	// Fields inherited from BackupDatasourceParameters

	ObjectType string `json:"objectType"`
}

func (s BlobBackupDatasourceParameters) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return BaseBackupDatasourceParametersImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = BlobBackupDatasourceParameters{}

func (s BlobBackupDatasourceParameters) MarshalJSON() ([]byte, error) {
	type wrapper BlobBackupDatasourceParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobBackupDatasourceParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobBackupDatasourceParameters: %+v", err)
	}

	decoded["objectType"] = "BlobBackupDatasourceParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobBackupDatasourceParameters: %+v", err)
	}

	return encoded, nil
}
