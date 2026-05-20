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
	var decoded struct {
		DataConnectionName *string `json:"dataConnectionName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DataConnectionName = decoded.DataConnectionName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DataConnectionValidation into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalDataConnectionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'DataConnectionValidation': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
