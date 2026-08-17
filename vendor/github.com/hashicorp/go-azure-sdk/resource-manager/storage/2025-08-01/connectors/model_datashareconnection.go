package connectors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageConnectorConnection = DataShareConnection{}

type DataShareConnection struct {
	DataShareUri string `json:"dataShareUri"`

	// Fields inherited from StorageConnectorConnection

	Type StorageConnectorConnectionType `json:"type"`
}

func (s DataShareConnection) StorageConnectorConnection() BaseStorageConnectorConnectionImpl {
	return BaseStorageConnectorConnectionImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = DataShareConnection{}

func (s DataShareConnection) MarshalJSON() ([]byte, error) {
	type wrapper DataShareConnection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DataShareConnection: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DataShareConnection: %+v", err)
	}

	decoded["type"] = "DataShare"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DataShareConnection: %+v", err)
	}

	return encoded, nil
}
