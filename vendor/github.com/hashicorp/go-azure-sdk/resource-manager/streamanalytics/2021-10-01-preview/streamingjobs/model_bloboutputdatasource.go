package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = BlobOutputDataSource{}

type BlobOutputDataSource struct {
	Properties *BlobOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource
}

var _ json.Marshaler = BlobOutputDataSource{}

func (s BlobOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper BlobOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobOutputDataSource: %+v", err)
	}
	decoded["type"] = "Microsoft.Storage/Blob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobOutputDataSource: %+v", err)
	}

	return encoded, nil
}
