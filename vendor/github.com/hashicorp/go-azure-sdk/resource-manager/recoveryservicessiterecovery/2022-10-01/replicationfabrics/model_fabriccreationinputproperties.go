package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricCreationInputProperties struct {
	CustomDetails FabricSpecificCreationInput `json:"customDetails"`
}

var _ json.Unmarshaler = &FabricCreationInputProperties{}

func (s *FabricCreationInputProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FabricCreationInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customDetails"]; ok {
		impl, err := unmarshalFabricSpecificCreationInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomDetails' for 'FabricCreationInputProperties': %+v", err)
		}
		s.CustomDetails = impl
	}
	return nil
}
