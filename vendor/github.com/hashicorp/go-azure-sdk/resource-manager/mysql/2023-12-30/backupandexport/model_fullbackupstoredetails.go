package backupandexport

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupStoreDetails = FullBackupStoreDetails{}

type FullBackupStoreDetails struct {
	SasUriList []string `json:"sasUriList"`

	// Fields inherited from BackupStoreDetails

	ObjectType string `json:"objectType"`
}

func (s FullBackupStoreDetails) BackupStoreDetails() BaseBackupStoreDetailsImpl {
	return BaseBackupStoreDetailsImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = FullBackupStoreDetails{}

func (s FullBackupStoreDetails) MarshalJSON() ([]byte, error) {
	type wrapper FullBackupStoreDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FullBackupStoreDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FullBackupStoreDetails: %+v", err)
	}

	decoded["objectType"] = "FullBackupStoreDetails"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FullBackupStoreDetails: %+v", err)
	}

	return encoded, nil
}
