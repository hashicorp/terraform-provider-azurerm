package sqldedicatedgateway

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceResourceCreateUpdateParameters struct {
	Properties ServiceResourceCreateUpdateProperties `json:"properties"`
}

var _ json.Unmarshaler = &ServiceResourceCreateUpdateParameters{}

func (s *ServiceResourceCreateUpdateParameters) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ServiceResourceCreateUpdateParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalServiceResourceCreateUpdatePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ServiceResourceCreateUpdateParameters': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
