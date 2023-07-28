package batchdeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AssetReferenceBase = DataPathAssetReference{}

type DataPathAssetReference struct {
	DatastoreId *string `json:"datastoreId,omitempty"`
	Path        *string `json:"path,omitempty"`

	// Fields inherited from AssetReferenceBase
}

var _ json.Marshaler = DataPathAssetReference{}

func (s DataPathAssetReference) MarshalJSON() ([]byte, error) {
	type wrapper DataPathAssetReference
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DataPathAssetReference: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DataPathAssetReference: %+v", err)
	}
	decoded["referenceType"] = "DataPath"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DataPathAssetReference: %+v", err)
	}

	return encoded, nil
}
