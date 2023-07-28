package batchdeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AssetReferenceBase = OutputPathAssetReference{}

type OutputPathAssetReference struct {
	JobId *string `json:"jobId,omitempty"`
	Path  *string `json:"path,omitempty"`

	// Fields inherited from AssetReferenceBase
}

var _ json.Marshaler = OutputPathAssetReference{}

func (s OutputPathAssetReference) MarshalJSON() ([]byte, error) {
	type wrapper OutputPathAssetReference
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OutputPathAssetReference: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OutputPathAssetReference: %+v", err)
	}
	decoded["referenceType"] = "OutputPath"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OutputPathAssetReference: %+v", err)
	}

	return encoded, nil
}
