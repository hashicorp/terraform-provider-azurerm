package openapis

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OperationProgressResponseType = BackupAndExportResponseType{}

type BackupAndExportResponseType struct {
	BackupMetadata         *string `json:"backupMetadata,omitempty"`
	DataTransferredInBytes *int64  `json:"dataTransferredInBytes,omitempty"`
	DatasourceSizeInBytes  *int64  `json:"datasourceSizeInBytes,omitempty"`

	// Fields inherited from OperationProgressResponseType

	ObjectType ObjectType `json:"objectType"`
}

func (s BackupAndExportResponseType) OperationProgressResponseType() BaseOperationProgressResponseTypeImpl {
	return BaseOperationProgressResponseTypeImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = BackupAndExportResponseType{}

func (s BackupAndExportResponseType) MarshalJSON() ([]byte, error) {
	type wrapper BackupAndExportResponseType
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BackupAndExportResponseType: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupAndExportResponseType: %+v", err)
	}

	decoded["objectType"] = "BackupAndExportResponse"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BackupAndExportResponseType: %+v", err)
	}

	return encoded, nil
}
