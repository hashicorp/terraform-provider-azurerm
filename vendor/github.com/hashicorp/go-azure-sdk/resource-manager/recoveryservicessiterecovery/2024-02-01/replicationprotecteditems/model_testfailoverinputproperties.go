package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestFailoverInputProperties struct {
	FailoverDirection       *string                           `json:"failoverDirection,omitempty"`
	NetworkId               *string                           `json:"networkId,omitempty"`
	NetworkType             *string                           `json:"networkType,omitempty"`
	ProviderSpecificDetails TestFailoverProviderSpecificInput `json:"providerSpecificDetails"`
}

var _ json.Unmarshaler = &TestFailoverInputProperties{}

func (s *TestFailoverInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias TestFailoverInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TestFailoverInputProperties: %+v", err)
	}

	s.FailoverDirection = decoded.FailoverDirection
	s.NetworkId = decoded.NetworkId
	s.NetworkType = decoded.NetworkType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TestFailoverInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := unmarshalTestFailoverProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'TestFailoverInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}
	return nil
}
