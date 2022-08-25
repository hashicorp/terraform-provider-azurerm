package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyParameters struct {
	DataStoreParametersList *[]DataStoreParameters `json:"dataStoreParametersList,omitempty"`
}

var _ json.Unmarshaler = &PolicyParameters{}

func (s *PolicyParameters) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PolicyParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["dataStoreParametersList"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling DataStoreParametersList into list []json.RawMessage: %+v", err)
		}

		output := make([]DataStoreParameters, 0)
		for i, val := range listTemp {
			impl, err := unmarshalDataStoreParametersImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'DataStoreParametersList' for 'PolicyParameters': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.DataStoreParametersList = &output
	}
	return nil
}
