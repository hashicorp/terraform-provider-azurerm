package batchdeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AssetReferenceBase = IdAssetReference{}

type IdAssetReference struct {
	AssetId string `json:"assetId"`

	// Fields inherited from AssetReferenceBase
}

var _ json.Marshaler = IdAssetReference{}

func (s IdAssetReference) MarshalJSON() ([]byte, error) {
	type wrapper IdAssetReference
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IdAssetReference: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IdAssetReference: %+v", err)
	}
	decoded["referenceType"] = "Id"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IdAssetReference: %+v", err)
	}

	return encoded, nil
}
