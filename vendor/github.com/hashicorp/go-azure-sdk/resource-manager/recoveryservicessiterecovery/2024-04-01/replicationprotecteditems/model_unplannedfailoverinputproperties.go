package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UnplannedFailoverInputProperties struct {
	FailoverDirection       *string                                `json:"failoverDirection,omitempty"`
	ProviderSpecificDetails UnplannedFailoverProviderSpecificInput `json:"providerSpecificDetails"`
	SourceSiteOperations    *string                                `json:"sourceSiteOperations,omitempty"`
}

var _ json.Unmarshaler = &UnplannedFailoverInputProperties{}

func (s *UnplannedFailoverInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FailoverDirection    *string `json:"failoverDirection,omitempty"`
		SourceSiteOperations *string `json:"sourceSiteOperations,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FailoverDirection = decoded.FailoverDirection
	s.SourceSiteOperations = decoded.SourceSiteOperations

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UnplannedFailoverInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalUnplannedFailoverProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'UnplannedFailoverInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
