package integrationruntimeenableinteractivequery

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SelfHostedIntegrationRuntimeTypeProperties struct {
	LinkedInfo                               LinkedIntegrationRuntimeType `json:"linkedInfo"`
	SelfContainedInteractiveAuthoringEnabled *bool                        `json:"selfContainedInteractiveAuthoringEnabled,omitempty"`
}

var _ json.Unmarshaler = &SelfHostedIntegrationRuntimeTypeProperties{}

func (s *SelfHostedIntegrationRuntimeTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		SelfContainedInteractiveAuthoringEnabled *bool `json:"selfContainedInteractiveAuthoringEnabled,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.SelfContainedInteractiveAuthoringEnabled = decoded.SelfContainedInteractiveAuthoringEnabled

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SelfHostedIntegrationRuntimeTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["linkedInfo"]; ok {
		impl, err := UnmarshalLinkedIntegrationRuntimeTypeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'LinkedInfo' for 'SelfHostedIntegrationRuntimeTypeProperties': %+v", err)
		}
		s.LinkedInfo = impl
	}

	return nil
}
