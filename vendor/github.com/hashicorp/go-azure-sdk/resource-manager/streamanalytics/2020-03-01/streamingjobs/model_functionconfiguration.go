package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionConfiguration struct {
	Binding FunctionBinding  `json:"binding"`
	Inputs  *[]FunctionInput `json:"inputs,omitempty"`
	Output  *FunctionOutput  `json:"output,omitempty"`
}

var _ json.Unmarshaler = &FunctionConfiguration{}

func (s *FunctionConfiguration) UnmarshalJSON(bytes []byte) error {
	type alias FunctionConfiguration
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into FunctionConfiguration: %+v", err)
	}

	s.Inputs = decoded.Inputs
	s.Output = decoded.Output

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FunctionConfiguration into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["binding"]; ok {
		impl, err := unmarshalFunctionBindingImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Binding' for 'FunctionConfiguration': %+v", err)
		}
		s.Binding = impl
	}
	return nil
}
