package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RestoreTargetInfoBase = ItemLevelRestoreTargetInfo{}

type ItemLevelRestoreTargetInfo struct {
	DatasourceAuthCredentials AuthCredentials            `json:"datasourceAuthCredentials"`
	DatasourceInfo            Datasource                 `json:"datasourceInfo"`
	DatasourceSetInfo         *DatasourceSet             `json:"datasourceSetInfo,omitempty"`
	RestoreCriteria           []ItemLevelRestoreCriteria `json:"restoreCriteria"`

	// Fields inherited from RestoreTargetInfoBase
	RecoveryOption  RecoveryOption `json:"recoveryOption"`
	RestoreLocation *string        `json:"restoreLocation,omitempty"`
}

var _ json.Marshaler = ItemLevelRestoreTargetInfo{}

func (s ItemLevelRestoreTargetInfo) MarshalJSON() ([]byte, error) {
	type wrapper ItemLevelRestoreTargetInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ItemLevelRestoreTargetInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ItemLevelRestoreTargetInfo: %+v", err)
	}
	decoded["objectType"] = "ItemLevelRestoreTargetInfo"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ItemLevelRestoreTargetInfo: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ItemLevelRestoreTargetInfo{}

func (s *ItemLevelRestoreTargetInfo) UnmarshalJSON(bytes []byte) error {
	type alias ItemLevelRestoreTargetInfo
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ItemLevelRestoreTargetInfo: %+v", err)
	}

	s.DatasourceInfo = decoded.DatasourceInfo
	s.DatasourceSetInfo = decoded.DatasourceSetInfo
	s.RecoveryOption = decoded.RecoveryOption
	s.RestoreLocation = decoded.RestoreLocation

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ItemLevelRestoreTargetInfo into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["datasourceAuthCredentials"]; ok {
		impl, err := unmarshalAuthCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DatasourceAuthCredentials' for 'ItemLevelRestoreTargetInfo': %+v", err)
		}
		s.DatasourceAuthCredentials = impl
	}

	if v, ok := temp["restoreCriteria"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling RestoreCriteria into list []json.RawMessage: %+v", err)
		}

		output := make([]ItemLevelRestoreCriteria, 0)
		for i, val := range listTemp {
			impl, err := unmarshalItemLevelRestoreCriteriaImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'RestoreCriteria' for 'ItemLevelRestoreTargetInfo': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.RestoreCriteria = output
	}
	return nil
}
