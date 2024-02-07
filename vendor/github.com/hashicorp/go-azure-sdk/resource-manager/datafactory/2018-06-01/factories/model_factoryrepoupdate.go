package factories

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FactoryRepoUpdate struct {
	FactoryResourceId *string                  `json:"factoryResourceId,omitempty"`
	RepoConfiguration FactoryRepoConfiguration `json:"repoConfiguration"`
}

var _ json.Unmarshaler = &FactoryRepoUpdate{}

func (s *FactoryRepoUpdate) UnmarshalJSON(bytes []byte) error {
	type alias FactoryRepoUpdate
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into FactoryRepoUpdate: %+v", err)
	}

	s.FactoryResourceId = decoded.FactoryResourceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FactoryRepoUpdate into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["repoConfiguration"]; ok {
		impl, err := unmarshalFactoryRepoConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RepoConfiguration' for 'FactoryRepoUpdate': %+v", err)
		}
		s.RepoConfiguration = impl
	}
	return nil
}
