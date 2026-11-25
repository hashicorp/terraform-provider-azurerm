package backupinstanceresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupDatasourceParameters = AdlsBlobBackupDatasourceParametersForAutoProtection{}

type AdlsBlobBackupDatasourceParametersForAutoProtection struct {
	AutoProtectionSettings BlobBackupRuleBasedAutoProtectionSettings `json:"autoProtectionSettings"`

	// Fields inherited from BackupDatasourceParameters

	ObjectType string `json:"objectType"`
}

func (s AdlsBlobBackupDatasourceParametersForAutoProtection) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return BaseBackupDatasourceParametersImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = AdlsBlobBackupDatasourceParametersForAutoProtection{}

func (s AdlsBlobBackupDatasourceParametersForAutoProtection) MarshalJSON() ([]byte, error) {
	type wrapper AdlsBlobBackupDatasourceParametersForAutoProtection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AdlsBlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AdlsBlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	decoded["objectType"] = "AdlsBlobBackupDatasourceParametersForAutoProtection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AdlsBlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	return encoded, nil
}
