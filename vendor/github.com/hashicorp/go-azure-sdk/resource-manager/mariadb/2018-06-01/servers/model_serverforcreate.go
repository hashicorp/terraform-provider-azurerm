package servers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerForCreate struct {
	Location   string                    `json:"location"`
	Properties ServerPropertiesForCreate `json:"properties"`
	Sku        *Sku                      `json:"sku,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
}

var _ json.Unmarshaler = &ServerForCreate{}

func (s *ServerForCreate) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Location string             `json:"location"`
		Sku      *Sku               `json:"sku,omitempty"`
		Tags     *map[string]string `json:"tags,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Location = decoded.Location
	s.Sku = decoded.Sku
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ServerForCreate into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalServerPropertiesForCreateImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ServerForCreate': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
