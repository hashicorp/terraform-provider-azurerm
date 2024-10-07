package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OperationExtendedInfo = OperationJobExtendedInfo{}

type OperationJobExtendedInfo struct {
	JobId *string `json:"jobId,omitempty"`

	// Fields inherited from OperationExtendedInfo

	ObjectType string `json:"objectType"`
}

func (s OperationJobExtendedInfo) OperationExtendedInfo() BaseOperationExtendedInfoImpl {
	return BaseOperationExtendedInfoImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = OperationJobExtendedInfo{}

func (s OperationJobExtendedInfo) MarshalJSON() ([]byte, error) {
	type wrapper OperationJobExtendedInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OperationJobExtendedInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OperationJobExtendedInfo: %+v", err)
	}

	decoded["objectType"] = "OperationJobExtendedInfo"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OperationJobExtendedInfo: %+v", err)
	}

	return encoded, nil
}
