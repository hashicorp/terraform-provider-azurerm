package accountconnectionresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionUpdateContent struct {
	Properties ConnectionPropertiesV2 `json:"properties"`
}

var _ json.Unmarshaler = &ConnectionUpdateContent{}

func (s *ConnectionUpdateContent) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ConnectionUpdateContent into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalConnectionPropertiesV2Implementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ConnectionUpdateContent': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
