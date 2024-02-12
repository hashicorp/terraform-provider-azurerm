package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BasePolicyRule = AzureBackupRule{}

type AzureBackupRule struct {
	BackupParameters BackupParameters  `json:"backupParameters"`
	DataStore        DataStoreInfoBase `json:"dataStore"`
	Trigger          TriggerContext    `json:"trigger"`

	// Fields inherited from BasePolicyRule
	Name string `json:"name"`
}

var _ json.Marshaler = AzureBackupRule{}

func (s AzureBackupRule) MarshalJSON() ([]byte, error) {
	type wrapper AzureBackupRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBackupRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupRule: %+v", err)
	}
	decoded["objectType"] = "AzureBackupRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBackupRule: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureBackupRule{}

func (s *AzureBackupRule) UnmarshalJSON(bytes []byte) error {
	type alias AzureBackupRule
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AzureBackupRule: %+v", err)
	}

	s.DataStore = decoded.DataStore
	s.Name = decoded.Name

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureBackupRule into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["backupParameters"]; ok {
		impl, err := unmarshalBackupParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'BackupParameters' for 'AzureBackupRule': %+v", err)
		}
		s.BackupParameters = impl
	}

	if v, ok := temp["trigger"]; ok {
		impl, err := unmarshalTriggerContextImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Trigger' for 'AzureBackupRule': %+v", err)
		}
		s.Trigger = impl
	}
	return nil
}
