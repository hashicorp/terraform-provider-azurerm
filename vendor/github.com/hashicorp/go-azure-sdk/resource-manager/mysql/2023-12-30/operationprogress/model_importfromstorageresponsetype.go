package operationprogress

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OperationProgressResponseType = ImportFromStorageResponseType{}

type ImportFromStorageResponseType struct {
	EstimatedCompletionTime *string `json:"estimatedCompletionTime,omitempty"`

	// Fields inherited from OperationProgressResponseType

	ObjectType ObjectType `json:"objectType"`
}

func (s ImportFromStorageResponseType) OperationProgressResponseType() BaseOperationProgressResponseTypeImpl {
	return BaseOperationProgressResponseTypeImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = ImportFromStorageResponseType{}

func (s ImportFromStorageResponseType) MarshalJSON() ([]byte, error) {
	type wrapper ImportFromStorageResponseType
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImportFromStorageResponseType: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImportFromStorageResponseType: %+v", err)
	}

	decoded["objectType"] = "ImportFromStorageResponse"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImportFromStorageResponseType: %+v", err)
	}

	return encoded, nil
}
