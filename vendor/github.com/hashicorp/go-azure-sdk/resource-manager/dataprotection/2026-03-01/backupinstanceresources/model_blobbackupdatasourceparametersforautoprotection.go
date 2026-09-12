package backupinstanceresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupDatasourceParameters = BlobBackupDatasourceParametersForAutoProtection{}

type BlobBackupDatasourceParametersForAutoProtection struct {
	AutoProtectionSettings BlobBackupRuleBasedAutoProtectionSettings `json:"autoProtectionSettings"`

	// Fields inherited from BackupDatasourceParameters

	ObjectType string `json:"objectType"`
}

func (s BlobBackupDatasourceParametersForAutoProtection) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return BaseBackupDatasourceParametersImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = BlobBackupDatasourceParametersForAutoProtection{}

func (s BlobBackupDatasourceParametersForAutoProtection) MarshalJSON() ([]byte, error) {
	type wrapper BlobBackupDatasourceParametersForAutoProtection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	decoded["objectType"] = "BlobBackupDatasourceParametersForAutoProtection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	return encoded, nil
}
