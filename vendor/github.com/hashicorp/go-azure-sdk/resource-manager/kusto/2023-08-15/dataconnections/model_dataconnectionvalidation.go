package dataconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectionValidation struct {
	DataConnectionName *string        `json:"dataConnectionName,omitempty"`
	Properties         DataConnection `json:"properties"`
}

var _ json.Unmarshaler = &DataConnectionValidation{}

func (s *DataConnectionValidation) UnmarshalJSON(bytes []byte) error {
	type alias DataConnectionValidation
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into DataConnectionValidation: %+v", err)
	}

	s.DataConnectionName = decoded.DataConnectionName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DataConnectionValidation into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalDataConnectionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'DataConnectionValidation': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
