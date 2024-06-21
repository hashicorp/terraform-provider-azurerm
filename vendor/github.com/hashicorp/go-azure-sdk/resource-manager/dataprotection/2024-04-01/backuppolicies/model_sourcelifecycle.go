package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceLifeCycle struct {
	DeleteAfter                 DeleteOption         `json:"deleteAfter"`
	SourceDataStore             DataStoreInfoBase    `json:"sourceDataStore"`
	TargetDataStoreCopySettings *[]TargetCopySetting `json:"targetDataStoreCopySettings,omitempty"`
}

var _ json.Unmarshaler = &SourceLifeCycle{}

func (s *SourceLifeCycle) UnmarshalJSON(bytes []byte) error {
	type alias SourceLifeCycle
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into SourceLifeCycle: %+v", err)
	}

	s.SourceDataStore = decoded.SourceDataStore
	s.TargetDataStoreCopySettings = decoded.TargetDataStoreCopySettings

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SourceLifeCycle into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["deleteAfter"]; ok {
		impl, err := unmarshalDeleteOptionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DeleteAfter' for 'SourceLifeCycle': %+v", err)
		}
		s.DeleteAfter = impl
	}
	return nil
}
