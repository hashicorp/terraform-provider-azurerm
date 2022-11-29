package inputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StreamInputDataSource = BlobStreamInputDataSource{}

type BlobStreamInputDataSource struct {
	Properties *BlobStreamInputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from StreamInputDataSource
}

var _ json.Marshaler = BlobStreamInputDataSource{}

func (s BlobStreamInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper BlobStreamInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobStreamInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobStreamInputDataSource: %+v", err)
	}
	decoded["type"] = "Microsoft.Storage/Blob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobStreamInputDataSource: %+v", err)
	}

	return encoded, nil
}
