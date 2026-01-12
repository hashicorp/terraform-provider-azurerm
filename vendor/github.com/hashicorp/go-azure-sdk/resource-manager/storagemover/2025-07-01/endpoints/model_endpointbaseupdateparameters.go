package endpoints

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointBaseUpdateParameters struct {
	Identity   *identity.LegacySystemAndUserAssignedMap `json:"identity,omitempty"`
	Properties EndpointBaseUpdateProperties             `json:"properties"`
}

var _ json.Unmarshaler = &EndpointBaseUpdateParameters{}

func (s *EndpointBaseUpdateParameters) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Identity *identity.LegacySystemAndUserAssignedMap `json:"identity,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Identity = decoded.Identity

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EndpointBaseUpdateParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalEndpointBaseUpdatePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'EndpointBaseUpdateParameters': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
