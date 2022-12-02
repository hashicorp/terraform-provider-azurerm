package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateProtectionContainerInputProperties struct {
	ProviderSpecificInput *[]ReplicationProviderSpecificContainerCreationInput `json:"providerSpecificInput,omitempty"`
}

var _ json.Unmarshaler = &CreateProtectionContainerInputProperties{}

func (s *CreateProtectionContainerInputProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CreateProtectionContainerInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificInput"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling ProviderSpecificInput into list []json.RawMessage: %+v", err)
		}

		output := make([]ReplicationProviderSpecificContainerCreationInput, 0)
		for i, val := range listTemp {
			impl, err := unmarshalReplicationProviderSpecificContainerCreationInputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'ProviderSpecificInput' for 'CreateProtectionContainerInputProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.ProviderSpecificInput = &output
	}
	return nil
}
