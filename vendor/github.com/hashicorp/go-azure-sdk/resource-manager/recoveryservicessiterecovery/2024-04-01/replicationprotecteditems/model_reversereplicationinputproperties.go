package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReverseReplicationInputProperties struct {
	FailoverDirection       *string                                 `json:"failoverDirection,omitempty"`
	ProviderSpecificDetails ReverseReplicationProviderSpecificInput `json:"providerSpecificDetails"`
}

var _ json.Unmarshaler = &ReverseReplicationInputProperties{}

func (s *ReverseReplicationInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FailoverDirection *string `json:"failoverDirection,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FailoverDirection = decoded.FailoverDirection

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ReverseReplicationInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalReverseReplicationProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'ReverseReplicationInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
