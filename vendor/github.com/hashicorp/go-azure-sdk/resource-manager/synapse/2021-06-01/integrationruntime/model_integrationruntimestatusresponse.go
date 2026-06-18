package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeStatusResponse struct {
	Name       *string                  `json:"name,omitempty"`
	Properties IntegrationRuntimeStatus `json:"properties"`
}

var _ json.Unmarshaler = &IntegrationRuntimeStatusResponse{}

func (s *IntegrationRuntimeStatusResponse) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Name *string `json:"name,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Name = decoded.Name

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling IntegrationRuntimeStatusResponse into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalIntegrationRuntimeStatusImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'IntegrationRuntimeStatusResponse': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
