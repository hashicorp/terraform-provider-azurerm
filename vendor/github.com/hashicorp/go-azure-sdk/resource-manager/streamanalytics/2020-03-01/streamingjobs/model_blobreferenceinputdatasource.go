package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReferenceInputDataSource = BlobReferenceInputDataSource{}

type BlobReferenceInputDataSource struct {
	Properties *BlobDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from ReferenceInputDataSource

	Type string `json:"type"`
}

func (s BlobReferenceInputDataSource) ReferenceInputDataSource() BaseReferenceInputDataSourceImpl {
	return BaseReferenceInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = BlobReferenceInputDataSource{}

func (s BlobReferenceInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper BlobReferenceInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobReferenceInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobReferenceInputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.Storage/Blob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobReferenceInputDataSource: %+v", err)
	}

	return encoded, nil
}
