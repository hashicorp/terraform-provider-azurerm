package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdatePolicyInputProperties struct {
	ReplicationProviderSettings PolicyProviderSpecificInput `json:"replicationProviderSettings"`
}

var _ json.Unmarshaler = &UpdatePolicyInputProperties{}

func (s *UpdatePolicyInputProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UpdatePolicyInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["replicationProviderSettings"]; ok {
		impl, err := unmarshalPolicyProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ReplicationProviderSettings' for 'UpdatePolicyInputProperties': %+v", err)
		}
		s.ReplicationProviderSettings = impl
	}
	return nil
}
