package databases

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Database = ReadWriteDatabase{}

type ReadWriteDatabase struct {
	Properties *ReadWriteDatabaseProperties `json:"properties,omitempty"`

	// Fields inherited from Database
	Id       *string `json:"id,omitempty"`
	Location *string `json:"location,omitempty"`
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
}

var _ json.Marshaler = ReadWriteDatabase{}

func (s ReadWriteDatabase) MarshalJSON() ([]byte, error) {
	type wrapper ReadWriteDatabase
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ReadWriteDatabase: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ReadWriteDatabase: %+v", err)
	}
	decoded["kind"] = "ReadWrite"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ReadWriteDatabase: %+v", err)
	}

	return encoded, nil
}
