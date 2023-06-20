package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RestoreTargetInfoBase = RestoreTargetInfo{}

type RestoreTargetInfo struct {
	DatasourceAuthCredentials AuthCredentials `json:"datasourceAuthCredentials"`
	DatasourceInfo            Datasource      `json:"datasourceInfo"`
	DatasourceSetInfo         *DatasourceSet  `json:"datasourceSetInfo,omitempty"`

	// Fields inherited from RestoreTargetInfoBase
	RecoveryOption  RecoveryOption `json:"recoveryOption"`
	RestoreLocation *string        `json:"restoreLocation,omitempty"`
}

var _ json.Marshaler = RestoreTargetInfo{}

func (s RestoreTargetInfo) MarshalJSON() ([]byte, error) {
	type wrapper RestoreTargetInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RestoreTargetInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RestoreTargetInfo: %+v", err)
	}
	decoded["objectType"] = "RestoreTargetInfo"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RestoreTargetInfo: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &RestoreTargetInfo{}

func (s *RestoreTargetInfo) UnmarshalJSON(bytes []byte) error {
	type alias RestoreTargetInfo
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into RestoreTargetInfo: %+v", err)
	}

	s.DatasourceInfo = decoded.DatasourceInfo
	s.DatasourceSetInfo = decoded.DatasourceSetInfo
	s.RecoveryOption = decoded.RecoveryOption
	s.RestoreLocation = decoded.RestoreLocation

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RestoreTargetInfo into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["datasourceAuthCredentials"]; ok {
		impl, err := unmarshalAuthCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DatasourceAuthCredentials' for 'RestoreTargetInfo': %+v", err)
		}
		s.DatasourceAuthCredentials = impl
	}
	return nil
}
