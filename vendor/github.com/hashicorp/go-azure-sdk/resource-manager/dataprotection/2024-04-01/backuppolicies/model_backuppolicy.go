package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BaseBackupPolicy = BackupPolicy{}

type BackupPolicy struct {
	PolicyRules []BasePolicyRule `json:"policyRules"`

	// Fields inherited from BaseBackupPolicy

	DatasourceTypes []string `json:"datasourceTypes"`
	ObjectType      string   `json:"objectType"`
}

func (s BackupPolicy) BaseBackupPolicy() BaseBaseBackupPolicyImpl {
	return BaseBaseBackupPolicyImpl{
		DatasourceTypes: s.DatasourceTypes,
		ObjectType:      s.ObjectType,
	}
}

var _ json.Marshaler = BackupPolicy{}

func (s BackupPolicy) MarshalJSON() ([]byte, error) {
	type wrapper BackupPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BackupPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupPolicy: %+v", err)
	}

	decoded["objectType"] = "BackupPolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BackupPolicy: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &BackupPolicy{}

func (s *BackupPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DatasourceTypes []string `json:"datasourceTypes"`
		ObjectType      string   `json:"objectType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DatasourceTypes = decoded.DatasourceTypes
	s.ObjectType = decoded.ObjectType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BackupPolicy into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["policyRules"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling PolicyRules into list []json.RawMessage: %+v", err)
		}

		output := make([]BasePolicyRule, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalBasePolicyRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'PolicyRules' for 'BackupPolicy': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.PolicyRules = output
	}

	return nil
}
