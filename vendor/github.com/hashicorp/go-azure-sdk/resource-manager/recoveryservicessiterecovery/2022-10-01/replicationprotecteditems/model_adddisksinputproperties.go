package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddDisksInputProperties struct {
	ProviderSpecificDetails AddDisksProviderSpecificInput `json:"providerSpecificDetails"`
}

var _ json.Unmarshaler = &AddDisksInputProperties{}

func (s *AddDisksInputProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AddDisksInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := unmarshalAddDisksProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'AddDisksInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}
	return nil
}
