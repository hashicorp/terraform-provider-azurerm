package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ InputProperties = ReferenceInputProperties{}

type ReferenceInputProperties struct {
	Datasource ReferenceInputDataSource `json:"datasource"`

	// Fields inherited from InputProperties
	Compression   *Compression  `json:"compression,omitempty"`
	Diagnostics   *Diagnostics  `json:"diagnostics,omitempty"`
	Etag          *string       `json:"etag,omitempty"`
	PartitionKey  *string       `json:"partitionKey,omitempty"`
	Serialization Serialization `json:"serialization"`
}

var _ json.Marshaler = ReferenceInputProperties{}

func (s ReferenceInputProperties) MarshalJSON() ([]byte, error) {
	type wrapper ReferenceInputProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ReferenceInputProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ReferenceInputProperties: %+v", err)
	}
	decoded["type"] = "Reference"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ReferenceInputProperties: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ReferenceInputProperties{}

func (s *ReferenceInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias ReferenceInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ReferenceInputProperties: %+v", err)
	}

	s.Compression = decoded.Compression
	s.Diagnostics = decoded.Diagnostics
	s.Etag = decoded.Etag
	s.PartitionKey = decoded.PartitionKey

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ReferenceInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["datasource"]; ok {
		impl, err := unmarshalReferenceInputDataSourceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Datasource' for 'ReferenceInputProperties': %+v", err)
		}
		s.Datasource = impl
	}

	if v, ok := temp["serialization"]; ok {
		impl, err := unmarshalSerializationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Serialization' for 'ReferenceInputProperties': %+v", err)
		}
		s.Serialization = impl
	}
	return nil
}
